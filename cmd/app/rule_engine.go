package app

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/malyshevhen/rule-engine/api"
	"github.com/malyshevhen/rule-engine/internal/alerting"
	"github.com/malyshevhen/rule-engine/internal/analytics"
	"github.com/malyshevhen/rule-engine/internal/core/action"
	"github.com/malyshevhen/rule-engine/internal/core/manager"
	"github.com/malyshevhen/rule-engine/internal/core/queue"
	"github.com/malyshevhen/rule-engine/internal/core/rule"
	"github.com/malyshevhen/rule-engine/internal/core/trigger"
	"github.com/malyshevhen/rule-engine/internal/engine/executor"
	execCtx "github.com/malyshevhen/rule-engine/internal/engine/executor/context"
	"github.com/malyshevhen/rule-engine/internal/engine/executor/platform"
	actionStorage "github.com/malyshevhen/rule-engine/internal/storage/action"
	"github.com/malyshevhen/rule-engine/internal/storage/db"
	redisClient "github.com/malyshevhen/rule-engine/internal/storage/redis"
	ruleStorage "github.com/malyshevhen/rule-engine/internal/storage/rule"
	triggerStorage "github.com/malyshevhen/rule-engine/internal/storage/trigger"
	"github.com/malyshevhen/rule-engine/pkg/tracing"
	"github.com/nats-io/nats.go"
	"github.com/robfig/cron/v3"
)

// ruleServiceAdapter adapts the rule.Service to manager.RuleService interface
type ruleServiceAdapter struct {
	ruleSvc *rule.Service
}

func (a *ruleServiceAdapter) GetByID(ctx context.Context, id uuid.UUID) (*rule.Rule, error) {
	return a.ruleSvc.GetByID(ctx, id)
}

// triggerServiceAdapter adapts the trigger.Service to manager.TriggerService interface
type triggerServiceAdapter struct {
	triggerSvc *trigger.Service
}

func (a *triggerServiceAdapter) GetByID(ctx context.Context, id uuid.UUID) (*trigger.Trigger, error) {
	return a.triggerSvc.GetByID(ctx, id)
}

func (a *triggerServiceAdapter) GetEnabledConditionalTriggers(ctx context.Context) ([]*trigger.Trigger, error) {
	return a.triggerSvc.GetEnabledConditionalTriggers(ctx)
}

func (a *triggerServiceAdapter) GetEnabledScheduledTriggers(ctx context.Context) ([]*trigger.Trigger, error) {
	return a.triggerSvc.GetEnabledScheduledTriggers(ctx)
}

// Config holds application configuration
type Config struct {
	Port               string
	DBURL              string
	NATSURL            string
	RedisURL           string
	AlertingEnabled    bool
	AlertWebhookURL    string
	AlertRetryAttempts int
}

// App represents the rule engine application
type App struct {
	config      Config
	db          *pgxpool.Pool
	redis       *redisClient.Client
	server      *api.Server
	manager     *manager.Manager
	workerPool  *queue.WorkerPool
	alertingSvc *alerting.Service
	nc          *nats.Conn
	cron        *cron.Cron
}

// loadConfig loads configuration from environment variables
func loadConfig() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://user:password@localhost:5432/rule_engine?sslmode=disable"
	}
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = "nats://localhost:4222"
	}
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "localhost:6379"
	}

	// Alerting configuration
	alertingEnabled := os.Getenv("ALERTING_ENABLED") == "true"
	alertWebhookURL := os.Getenv("ALERT_WEBHOOK_URL")
	alertRetryAttempts := 3 // default
	if retryStr := os.Getenv("ALERT_RETRY_ATTEMPTS"); retryStr != "" {
		if retry, err := strconv.Atoi(retryStr); err == nil && retry > 0 {
			alertRetryAttempts = retry
		}
	}

	return Config{
		Port:               port,
		DBURL:              dbURL,
		NATSURL:            natsURL,
		RedisURL:           redisURL,
		AlertingEnabled:    alertingEnabled,
		AlertWebhookURL:    alertWebhookURL,
		AlertRetryAttempts: alertRetryAttempts,
	}
}

// New creates a new App instance
func New() *App {
	// Set up structured logging
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	// Initialize OpenTelemetry tracing
	if err := tracing.InitTracing("rule-engine", "1.0.0"); err != nil {
		slog.Error("Failed to initialize tracing", "error", err)
		os.Exit(1)
	}
	slog.Info("OpenTelemetry tracing initialized")

	config := loadConfig()

	// Initialize database connection
	ctx := context.Background()
	pool, err := db.NewPostgresPool(ctx, config.DBURL)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	slog.Info("Connected to database")

	// Run database migrations
	if err := db.RunMigrations(pool); err != nil {
		slog.Error("Failed to run migrations", "error", err)
		os.Exit(1)
	}
	slog.Info("Database migrations completed")

	// Initialize repositories
	ruleRepo := ruleStorage.NewRepository(pool)
	triggerRepo := triggerStorage.NewRepository(pool)
	actionRepo := actionStorage.NewRepository(pool)

	// Initialize Redis client
	redisConfig := &redisClient.Config{
		Addr: config.RedisURL,
	}
	redisCli := redisClient.NewClient(redisConfig)

	// Test Redis connection
	if err := redisCli.Ping(ctx); err != nil {
		slog.Warn("Failed to connect to Redis, caching will be disabled", "error", err)
		redisCli = nil
	} else {
		slog.Info("Connected to Redis")
	}

	// Initialize services
	ruleSvc := rule.NewService(ruleRepo, triggerRepo, actionRepo, redisCli)
	triggerSvc := trigger.NewService(triggerRepo, redisCli)
	actionSvc := action.NewService(actionRepo)

	// Initialize executor components
	contextSvc := execCtx.NewService()
	platformSvc := platform.NewService()
	executorSvc := executor.NewService(contextSvc, platformSvc)

	// Initialize trigger evaluator
	triggerEval := trigger.NewEvaluator(executorSvc)

	// Initialize execution queue (use Redis if available, otherwise in-memory)
	var execQueue queue.Queue
	if redisCli != nil {
		execQueue = queue.NewRedisQueue(redisCli, "rule_engine:queue")
		slog.Info("Using Redis-backed execution queue")
	} else {
		execQueue = queue.NewInMemoryQueue()
		slog.Info("Using in-memory execution queue")
	}

	workerPool := queue.NewWorkerPool(execQueue, ruleSvc, executorSvc, 5)
	workerPool.Start(ctx)

	// Initialize alerting service
	alertingConfig := alerting.Config{
		Enabled:       config.AlertingEnabled,
		WebhookURL:    config.AlertWebhookURL,
		RetryAttempts: config.AlertRetryAttempts,
	}
	alertingSvc := alerting.NewService(alertingConfig)

	// Initialize analytics service
	analyticsSvc := analytics.NewService()

	// Initialize NATS connection
	nc, err := nats.Connect(config.NATSURL)
	if err != nil {
		slog.Error("Failed to connect to NATS", "error", err)
		os.Exit(1)
	}
	slog.Info("Connected to NATS")

	// Initialize cron scheduler
	c := cron.New()

	// Initialize trigger manager
	ruleSvcForManager := &ruleServiceAdapter{ruleSvc: ruleSvc}
	triggerSvcForManager := &triggerServiceAdapter{triggerSvc: triggerSvc}
	mgr := manager.NewManager(nc, c, ruleSvcForManager, triggerSvcForManager, triggerEval, executorSvc, alertingSvc, execQueue)

	// Initialize HTTP server
	serverConfig := &api.ServerConfig{Port: config.Port}
	server := api.NewServer(serverConfig, ruleSvc, triggerSvc, actionSvc, analyticsSvc)

	return &App{
		config:      config,
		db:          pool,
		redis:       redisCli,
		server:      server,
		manager:     mgr,
		workerPool:  workerPool,
		alertingSvc: alertingSvc,
		nc:          nc,
		cron:        c,
	}
}

// Run starts the application
func (a *App) Run() error {
	slog.Info("Starting rule engine app", "port", a.config.Port)

	// Start cron scheduler
	a.cron.Start()

	// Start trigger manager
	ctx := context.Background()
	if err := a.manager.Start(ctx); err != nil {
		slog.Error("Failed to start trigger manager", "error", err)
		os.Exit(1)
	}
	slog.Info("Trigger manager started")

	// Start rate limiter cleanup
	api.StartRateLimiterCleanup()
	slog.Info("Rate limiter cleanup started")

	// Start server in a goroutine
	go func() {
		if err := a.server.Start(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for shutdown signal
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	<-ctx.Done()
	slog.Info("Shutting down rule engine app")

	// Gracefully shutdown server
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := a.server.Shutdown(shutdownCtx); err != nil {
		slog.Error("Server shutdown failed", "error", err)
	}

	a.manager.Stop()
	a.workerPool.Stop()
	a.cron.Stop()
	a.nc.Close()
	a.db.Close()

	// Stop rate limiter cleanup
	api.StopRateLimiterCleanup()
	slog.Info("Rate limiter cleanup stopped")

	// Shutdown tracing
	if err := tracing.ShutdownTracing(shutdownCtx); err != nil {
		slog.Error("Tracing shutdown failed", "error", err)
	}

	return nil
}
