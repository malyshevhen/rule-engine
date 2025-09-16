package app

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/malyshevhen/rule-engine/api"
	"github.com/malyshevhen/rule-engine/internal/storage/db"
)

// Config holds application configuration
type Config struct {
	Port  string
	DBURL string
}

// App represents the rule engine application
type App struct {
	config Config
	db     *pgxpool.Pool
	server *api.Server
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
	return Config{
		Port:  port,
		DBURL: dbURL,
	}
}

// New creates a new App instance
func New() *App {
	// Set up structured logging
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

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

	// Initialize HTTP server
	serverConfig := &api.ServerConfig{Port: config.Port}
	server := api.NewServer(serverConfig)

	return &App{
		config: config,
		db:     pool,
		server: server,
	}
}

// Run starts the application
func (a *App) Run() error {
	slog.Info("Starting rule engine app", "port", a.config.Port)

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

	a.db.Close()
	return nil
}
