package api

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/malyshevhen/rule-engine/internal/analytics"
	"github.com/malyshevhen/rule-engine/internal/core/action"
	"github.com/malyshevhen/rule-engine/internal/core/rule"
	"github.com/malyshevhen/rule-engine/internal/core/trigger"
	ruleStorage "github.com/malyshevhen/rule-engine/internal/storage/rule"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// RuleService interface
type RuleService interface {
	Create(ctx context.Context, rule *rule.Rule) error
	GetByID(ctx context.Context, id uuid.UUID) (*rule.Rule, error)
	List(ctx context.Context, limit int, offset int) ([]*rule.Rule, error)
	ListAll(ctx context.Context) ([]*rule.Rule, error)
	Update(ctx context.Context, rule *rule.Rule) error
	Delete(ctx context.Context, id uuid.UUID) error
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

// ServerConfig holds server configuration
type ServerConfig struct {
	Port string
}

// AnalyticsService interface
type AnalyticsService interface {
	GetDashboardData(ctx context.Context, timeRange string) (*analytics.DashboardData, error)
}

// ServerOption is a functional option for configuring the server
type ServerOption func(*Server)

// WithRateLimiting enables or disables rate limiting
func WithRateLimiting(enabled bool) ServerOption {
	return func(s *Server) {
		s.rateLimitingEnabled = enabled
	}
}

// Server represents the HTTP server
type Server struct {
	config              *ServerConfig
	router              *mux.Router
	server              *http.Server
	ruleSvc             RuleService
	triggerSvc          TriggerService
	actionSvc           ActionService
	analyticsSvc        AnalyticsService
	rateLimitingEnabled bool
}

// NewServer creates a new HTTP server with optional configuration
func NewServer(config *ServerConfig, ruleSvc RuleService, triggerSvc TriggerService, actionSvc ActionService, analyticsSvc AnalyticsService, opts ...ServerOption) *Server {
	router := mux.NewRouter()

	s := &Server{
		config:              config,
		router:              router,
		ruleSvc:             ruleSvc,
		triggerSvc:          triggerSvc,
		actionSvc:           actionSvc,
		analyticsSvc:        analyticsSvc,
		rateLimitingEnabled: true, // Default to enabled
	}

	// Apply options
	for _, opt := range opts {
		opt(s)
	}

	// Add middleware based on configuration
	if s.rateLimitingEnabled {
		router.Use(RateLimitMiddleware)
	}
	router.Use(TracingMiddleware)
	router.Use(LoggingMiddleware)

	// TODO: Setup CORS

	s.setupRoutes()

	addr := ":" + config.Port
	s.server = &http.Server{
		Addr:    addr,
		Handler: router,
	}

	return s
}

// NewServerWithoutRateLimit creates a new HTTP server without rate limiting (for performance testing)
func NewServerWithoutRateLimit(config *ServerConfig, ruleSvc RuleService, triggerSvc TriggerService, actionSvc ActionService, analyticsSvc AnalyticsService) *Server {
	return NewServer(config, ruleSvc, triggerSvc, actionSvc, analyticsSvc, WithRateLimiting(false))
}

// setupRoutes registers all API routes
func (s *Server) setupRoutes() {
	// Public routes (no authentication required)
	s.router.HandleFunc("/health", s.HealthCheck).Methods("GET")
	s.router.HandleFunc("/dashboard", s.ServeDashboard).Methods("GET")
	s.router.Handle("/metrics", promhttp.Handler())

	// Protected API routes (require authentication)
	api := s.router.PathPrefix("/api/v1").Subrouter()
	api.Use(AuthMiddleware)

	// Rules routes
	api.HandleFunc("/rules", s.CreateRule).Methods("POST")
	api.HandleFunc("/rules", s.ListRules).Methods("GET")
	api.HandleFunc("/rules/{id}", s.GetRule).Methods("GET")
	api.HandleFunc("/rules/{id}", s.UpdateRule).Methods("PUT")
	api.HandleFunc("/rules/{id}", s.DeleteRule).Methods("DELETE")

	// Triggers routes
	api.HandleFunc("/triggers", s.CreateTrigger).Methods("POST")
	api.HandleFunc("/triggers", s.ListTriggers).Methods("GET")
	api.HandleFunc("/triggers/{id}", s.GetTrigger).Methods("GET")

	// Actions routes
	api.HandleFunc("/actions", s.CreateAction).Methods("POST")
	api.HandleFunc("/actions", s.ListActions).Methods("GET")
	api.HandleFunc("/actions/{id}", s.GetAction).Methods("GET")

	// Analytics routes
	api.HandleFunc("/analytics/dashboard", s.GetDashboardData).Methods("GET")
}

// Start starts the HTTP server
func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

// Router returns the router for testing purposes
func (s *Server) Router() *mux.Router {
	return s.router
}

// GetDashboardData godoc
//
//	@Summary		Get analytics dashboard data
//	@Description	Get aggregated analytics data for the dashboard
//	@Tags			analytics
//	@Accept			json
//	@Produce		json
//	@Param			timeRange	query		string	false	"Time range (1h, 24h, 7d, 30d)"	Enums(1h,24h,7d,30d)
//	@Success		200			{object}	analytics.DashboardData
//	@Failure		500			{object}	APIErrorResponse
//	@Router			/analytics/dashboard [get]
func (s *Server) GetDashboardData(w http.ResponseWriter, r *http.Request) {
	timeRange := r.URL.Query().Get("timeRange")
	if timeRange == "" {
		timeRange = "24h" // default
	}

	// Validate time range
	validRanges := map[string]bool{
		"1h": true, "24h": true, "7d": true, "30d": true,
	}
	if !validRanges[timeRange] {
		ErrorResponse(w, http.StatusBadRequest, "Invalid time range. Valid values: 1h, 24h, 7d, 30d")
		return
	}

	data, err := s.analyticsSvc.GetDashboardData(r.Context(), timeRange)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Failed to get dashboard data")
		return
	}

	SuccessResponse(w, data)
}

// ServeDashboard serves the analytics dashboard HTML page
func (s *Server) ServeDashboard(w http.ResponseWriter, r *http.Request) {
	// For now, serve a simple HTML response
	// In a production system, you might want to embed the HTML file or serve it from a static directory
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Rule Engine Analytics Dashboard</title>
    <script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; background: #f5f5f5; }
        .container { max-width: 1200px; margin: 0 auto; }
        .header { text-align: center; margin-bottom: 30px; }
        .stats { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 20px; margin-bottom: 30px; }
        .stat-card { background: white; padding: 20px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); text-align: center; }
        .stat-value { font-size: 2em; font-weight: bold; color: #007bff; }
        .stat-label { color: #666; margin-top: 5px; }
        .chart-container { background: white; padding: 20px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); margin-bottom: 20px; }
        .loading { text-align: center; padding: 50px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🚀 Rule Engine Analytics Dashboard</h1>
            <p>Real-time monitoring and insights</p>
        </div>

        <div class="stats" id="stats">
            <div class="stat-card">
                <div class="stat-value" id="total-executions">Loading...</div>
                <div class="stat-label">Total Executions</div>
            </div>
            <div class="stat-card">
                <div class="stat-value" id="success-rate">Loading...</div>
                <div class="stat-label">Success Rate</div>
            </div>
            <div class="stat-card">
                <div class="stat-value" id="avg-latency">Loading...</div>
                <div class="stat-label">Avg Latency</div>
            </div>
        </div>

        <div class="chart-container">
            <h3>📈 Execution Trend (Last 24 Hours)</h3>
            <canvas id="executionChart" width="400" height="200"></canvas>
        </div>

        <div class="chart-container">
            <h3>✅ Success Rate Trend</h3>
            <canvas id="successRateChart" width="400" height="200"></canvas>
        </div>
    </div>

    <script>
        let executionChart, successRateChart;

        function initCharts() {
            executionChart = new Chart(document.getElementById('executionChart'), {
                type: 'line',
                data: { labels: [], datasets: [{ label: 'Executions', data: [], borderColor: '#007bff', tension: 0.4 }] },
                options: { responsive: true, maintainAspectRatio: false }
            });

            successRateChart = new Chart(document.getElementById('successRateChart'), {
                type: 'line',
                data: { labels: [], datasets: [{ label: 'Success Rate %', data: [], borderColor: '#28a745', tension: 0.4 }] },
                options: { responsive: true, maintainAspectRatio: false }
            });
        }

        async function loadDashboard() {
            try {
                const response = await fetch('/api/v1/analytics/dashboard?timeRange=24h');
                const data = await response.json();

                // Update stats
                document.getElementById('total-executions').textContent = data.overall_stats.total_executions;
                document.getElementById('success-rate').textContent = data.overall_stats.success_rate.toFixed(1) + '%';
                document.getElementById('avg-latency').textContent = data.overall_stats.average_latency.toFixed(1) + 'ms';

                // Update charts
                const execLabels = data.execution_trend.data.map(p => new Date(p.timestamp).toLocaleTimeString());
                const execData = data.execution_trend.data.map(p => p.value);
                executionChart.data.labels = execLabels;
                executionChart.data.datasets[0].data = execData;
                executionChart.update();

                const successLabels = data.success_rate_trend.data.map(p => new Date(p.timestamp).toLocaleTimeString());
                const successData = data.success_rate_trend.data.map(p => p.value);
                successRateChart.data.labels = successLabels;
                successRateChart.data.datasets[0].data = successData;
                successRateChart.update();

            } catch (error) {
                console.error('Failed to load dashboard data:', error);
            }
        }

        document.addEventListener('DOMContentLoaded', function() {
            initCharts();
            loadDashboard();
            // Refresh every 30 seconds
            setInterval(loadDashboard, 30000);
        });
    </script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

// HealthCheck godoc
//
//	@Summary		Health check
//	@Description	Get the health status of the service
//	@Tags			system
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	map[string]string
//	@Router			/health [get]
func (s *Server) HealthCheck(w http.ResponseWriter, r *http.Request) {
	SuccessResponse(w, map[string]string{"status": "healthy"})
}

// CreateRule godoc
//
//	@Summary		Create a new rule
//	@Description	Create a new automation rule with Lua script
//	@Tags			rules
//	@Accept			json
//	@Produce		json
//	@Param			rule	body		CreateRuleRequest	true	"Rule data"
//	@Success		200		{object}	RuleDTO
//	@Failure		400		{object}	APIErrorResponse
//	@Failure		500		{object}	APIErrorResponse
//	@Router			/rules [post]
func (s *Server) CreateRule(w http.ResponseWriter, r *http.Request) {
	var req CreateRuleRequest
	if err := ParseJSONBody(r, &req); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Sanitize and validate inputs
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		ErrorResponse(w, http.StatusBadRequest, "Rule name cannot be empty")
		return
	}
	if len(req.Name) > 255 {
		ErrorResponse(w, http.StatusBadRequest, "Rule name too long (max 255 characters)")
		return
	}

	req.LuaScript = strings.TrimSpace(req.LuaScript)
	if req.LuaScript == "" {
		ErrorResponse(w, http.StatusBadRequest, "Lua script cannot be empty")
		return
	}
	if len(req.LuaScript) > 10000 {
		ErrorResponse(w, http.StatusBadRequest, "Lua script too long (max 10000 characters)")
		return
	}

	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}

	priority := 0
	if req.Priority != nil {
		priority = *req.Priority
	}

	rule := &rule.Rule{
		Name:      req.Name,
		LuaScript: req.LuaScript,
		Priority:  priority,
		Enabled:   enabled,
	}

	if err := s.ruleSvc.Create(r.Context(), rule); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Failed to create rule")
		return
	}

	SuccessResponse(w, rule)
}

// ListRules godoc
//
//	@Summary		List all rules
//	@Description	Get a list of all automation rules
//	@Tags			rules
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		RuleDTO
//	@Failure		500	{object}	APIErrorResponse
//	@Router			/rules [get]
func (s *Server) ListRules(w http.ResponseWriter, r *http.Request) {
	rules, err := s.ruleSvc.ListAll(r.Context())
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Failed to list rules")
		return
	}

	SuccessResponse(w, rules)
}

// GetRule godoc
//
//	@Summary		Get rule by ID
//	@Description	Get a specific automation rule by its ID
//	@Tags			rules
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Rule ID"
//	@Success		200	{object}	RuleDTO
//	@Failure		400	{object}	APIErrorResponse
//	@Failure		404	{object}	APIErrorResponse
//	@Router			/rules/{id} [get]
func (s *Server) GetRule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid rule ID")
		return
	}

	rule, err := s.ruleSvc.GetByID(r.Context(), id)
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "Rule not found")
		return
	}

	SuccessResponse(w, rule)
}

// UpdateRule godoc
//
//	@Summary		Update rule
//	@Description	Update an existing automation rule
//	@Tags			rules
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string				true	"Rule ID"
//	@Param			rule	body		UpdateRuleRequest	true	"Updated rule data"
//	@Success		200		{object}	RuleDTO
//	@Failure		400		{object}	APIErrorResponse
//	@Failure		404		{object}	APIErrorResponse
//	@Failure		500		{object}	APIErrorResponse
//	@Router			/rules/{id} [put]
func (s *Server) UpdateRule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid rule ID")
		return
	}

	var req UpdateRuleRequest
	if err := ParseJSONBody(r, &req); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Get existing rule
	existingRule, err := s.ruleSvc.GetByID(r.Context(), id)
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "Rule not found")
		return
	}

	// Apply updates
	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			ErrorResponse(w, http.StatusBadRequest, "Rule name cannot be empty")
			return
		}
		if len(name) > 255 {
			ErrorResponse(w, http.StatusBadRequest, "Rule name too long (max 255 characters)")
			return
		}
		existingRule.Name = name
	}

	if req.LuaScript != nil {
		script := strings.TrimSpace(*req.LuaScript)
		if script == "" {
			ErrorResponse(w, http.StatusBadRequest, "Lua script cannot be empty")
			return
		}
		if len(script) > 10000 {
			ErrorResponse(w, http.StatusBadRequest, "Lua script too long (max 10000 characters)")
			return
		}
		existingRule.LuaScript = script
	}

	if req.Enabled != nil {
		existingRule.Enabled = *req.Enabled
	}

	if req.Priority != nil {
		existingRule.Priority = *req.Priority
	}

	if err := s.ruleSvc.Update(r.Context(), existingRule); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Failed to update rule")
		return
	}

	SuccessResponse(w, existingRule)
}

// DeleteRule godoc
//
//	@Summary		Delete rule
//	@Description	Delete an automation rule by its ID
//	@Tags			rules
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Rule ID"
//	@Success		204	"No Content"
//	@Failure		400	{object}	APIErrorResponse
//	@Failure		404	{object}	APIErrorResponse
//	@Failure		500	{object}	APIErrorResponse
//	@Router			/rules/{id} [delete]
func (s *Server) DeleteRule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid rule ID")
		return
	}

	if err := s.ruleSvc.Delete(r.Context(), id); err != nil {
		if errors.Is(err, ruleStorage.ErrNotFound) {
			ErrorResponse(w, http.StatusNotFound, "Rule not found")
			return
		}
		ErrorResponse(w, http.StatusInternalServerError, "Failed to delete rule")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// CreateTrigger godoc
//
//	@Summary		Create a new trigger
//	@Description	Create a new trigger for rule execution
//	@Tags			triggers
//	@Accept			json
//	@Produce		json
//	@Param			trigger	body		CreateTriggerRequest	true	"Trigger data"
//	@Success		200		{object}	TriggerDTO
//	@Failure		400		{object}	APIErrorResponse
//	@Failure		500		{object}	APIErrorResponse
//	@Router			/triggers [post]
func (s *Server) CreateTrigger(w http.ResponseWriter, r *http.Request) {
	var req CreateTriggerRequest
	if err := ParseJSONBody(r, &req); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Sanitize and validate inputs
	req.Type = strings.TrimSpace(req.Type)
	if req.Type == "" {
		ErrorResponse(w, http.StatusBadRequest, "Trigger type cannot be empty")
		return
	}
	if req.Type != "CONDITIONAL" && req.Type != "CRON" {
		ErrorResponse(w, http.StatusBadRequest, "Invalid trigger type (must be 'CONDITIONAL' or 'CRON')")
		return
	}

	req.ConditionScript = strings.TrimSpace(req.ConditionScript)
	if req.ConditionScript == "" {
		ErrorResponse(w, http.StatusBadRequest, "Condition script cannot be empty")
		return
	}
	if len(req.ConditionScript) > 5000 {
		ErrorResponse(w, http.StatusBadRequest, "Condition script too long (max 5000 characters)")
		return
	}

	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}

	trigger := &trigger.Trigger{
		RuleID:          req.RuleID,
		Type:            trigger.TriggerType(req.Type),
		ConditionScript: req.ConditionScript,
		Enabled:         enabled,
	}

	if err := s.triggerSvc.Create(r.Context(), trigger); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Failed to create trigger")
		return
	}

	SuccessResponse(w, trigger)
}

// ListTriggers godoc
//
//	@Summary		List all triggers
//	@Description	Get a list of all triggers
//	@Tags			triggers
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		TriggerDTO
//	@Failure		500	{object}	APIErrorResponse
//	@Router			/triggers [get]
func (s *Server) ListTriggers(w http.ResponseWriter, r *http.Request) {
	triggers, err := s.triggerSvc.List(r.Context())
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Failed to list triggers")
		return
	}

	SuccessResponse(w, triggers)
}

// GetTrigger godoc
//
//	@Summary		Get trigger by ID
//	@Description	Get a specific trigger by its ID
//	@Tags			triggers
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Trigger ID"
//	@Success		200	{object}	TriggerDTO
//	@Failure		400	{object}	APIErrorResponse
//	@Failure		404	{object}	APIErrorResponse
//	@Router			/triggers/{id} [get]
func (s *Server) GetTrigger(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid trigger ID")
		return
	}

	trigger, err := s.triggerSvc.GetByID(r.Context(), id)
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "Trigger not found")
		return
	}

	SuccessResponse(w, trigger)
}

// CreateAction godoc
//
//	@Summary		Create a new action
//	@Description	Create a new action for rule execution
//	@Tags			actions
//	@Accept			json
//	@Produce		json
//	@Param			action	body		CreateActionRequest	true	"Action data"
//	@Success		200		{object}	ActionDTO
//	@Failure		400		{object}	APIErrorResponse
//	@Failure		500		{object}	APIErrorResponse
//	@Router			/actions [post]
func (s *Server) CreateAction(w http.ResponseWriter, r *http.Request) {
	var req CreateActionRequest
	if err := ParseJSONBody(r, &req); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Sanitize and validate inputs
	req.LuaScript = strings.TrimSpace(req.LuaScript)
	if req.LuaScript == "" {
		ErrorResponse(w, http.StatusBadRequest, "Lua script cannot be empty")
		return
	}
	if len(req.LuaScript) > 10000 {
		ErrorResponse(w, http.StatusBadRequest, "Lua script too long (max 10000 characters)")
		return
	}

	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}

	action := &action.Action{
		LuaScript: req.LuaScript,
		Enabled:   enabled,
	}

	if err := s.actionSvc.Create(r.Context(), action); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Failed to create action")
		return
	}

	SuccessResponse(w, action)
}

// ListActions godoc
//
//	@Summary		List all actions
//	@Description	Get a list of all actions
//	@Tags			actions
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		ActionDTO
//	@Failure		500	{object}	APIErrorResponse
//	@Router			/actions [get]
func (s *Server) ListActions(w http.ResponseWriter, r *http.Request) {
	actions, err := s.actionSvc.List(r.Context())
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Failed to list actions")
		return
	}

	SuccessResponse(w, actions)
}

// GetAction godoc
//
//	@Summary		Get action by ID
//	@Description	Get a specific action by its ID
//	@Tags			actions
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Action ID"
//	@Success		200	{object}	ActionDTO
//	@Failure		400	{object}	APIErrorResponse
//	@Failure		404	{object}	APIErrorResponse
//	@Router			/actions/{id} [get]
func (s *Server) GetAction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid action ID")
		return
	}

	action, err := s.actionSvc.GetByID(r.Context(), id)
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "Action not found")
		return
	}

	SuccessResponse(w, action)
}
