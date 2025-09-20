package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/malyshevhen/rule-engine/internal/action"
	"github.com/malyshevhen/rule-engine/internal/analytics"
	"github.com/malyshevhen/rule-engine/internal/engine/executor"
	execCtx "github.com/malyshevhen/rule-engine/internal/engine/executor/context"
	"github.com/malyshevhen/rule-engine/internal/rule"
	"github.com/malyshevhen/rule-engine/internal/trigger"
)

// RuleService interface
type RuleService interface {
	Create(ctx context.Context, rule *rule.Rule) error
	GetByID(ctx context.Context, id uuid.UUID) (*rule.Rule, error)
	List(ctx context.Context, limit int, offset int) ([]*rule.Rule, error)
	ListAll(ctx context.Context) ([]*rule.Rule, error)
	Update(ctx context.Context, rule *rule.Rule) error
	Delete(ctx context.Context, id uuid.UUID) error
	AddAction(ctx context.Context, ruleID, actionID uuid.UUID) error
}

// TriggerService interface
type TriggerService interface {
	Create(ctx context.Context, trigger *trigger.Trigger) error
	GetByID(ctx context.Context, id uuid.UUID) (*trigger.Trigger, error)
	List(ctx context.Context) ([]*trigger.Trigger, error)
}

// ActionService interface
type ActionService interface {
	Create(ctx context.Context, action *action.Action) error
	GetByID(ctx context.Context, id uuid.UUID) (*action.Action, error)
	List(ctx context.Context) ([]*action.Action, error)
}

// AnalyticsService interface
type AnalyticsService interface {
	GetDashboardData(ctx context.Context, timeRange string) (*analytics.DashboardData, error)
}

// ExecutorService interface
type ExecutorService interface {
	ExecuteScript(ctx context.Context, script string, execCtx *execCtx.ExecutionContext) *executor.ExecuteResult
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port string
}

// NewServer creates a new HTTP server with optional configuration
func NewServer(
	config *ServerConfig,
	healthSvc *Health,
	ruleSvc RuleService,
	triggerSvc TriggerService,
	actionSvc ActionService,
	analyticsSvc AnalyticsService,
	executorSvc ExecutorService,
	rateLimitingEnabled bool,
) *http.Server {
	router := setupRoutes(
		healthSvc,
		executorSvc,
		ruleSvc,
		triggerSvc,
		actionSvc,
	)

	recoveryHandler := handlers.RecoveryHandler()
	corsHandler := handlers.CORS(
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "OPTIONS"}),
		handlers.AllowCredentials(),
		handlers.ExposedHeaders([]string{"Authorization", "Content-Type", "Content-Encoding", "Content-Length", "Location"}),
	)

	handler := loggingMiddleware(tracingMiddleware(corsHandler(recoveryHandler(router))))

	if rateLimitingEnabled {
		handler = rateLimitMiddleware(handler)
	}

	return &http.Server{
		Addr:         fmt.Sprintf(":%s", config.Port),
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      handler,
	}
}
