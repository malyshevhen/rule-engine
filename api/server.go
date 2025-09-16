package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/malyshevhen/rule-engine/internal/core/action"
	"github.com/malyshevhen/rule-engine/internal/core/rule"
	"github.com/malyshevhen/rule-engine/internal/core/trigger"
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

// Server represents the HTTP server
type Server struct {
	config     *ServerConfig
	router     *mux.Router
	server     *http.Server
	ruleSvc    RuleService
	triggerSvc TriggerService
	actionSvc  ActionService
}

// NewServer creates a new HTTP server
func NewServer(config *ServerConfig, ruleSvc *rule.Service, triggerSvc *trigger.Service, actionSvc *action.Service) *Server {
	router := mux.NewRouter()

	// Add middleware
	router.Use(RateLimitMiddleware)
	router.Use(TracingMiddleware)
	router.Use(LoggingMiddleware)
	router.Use(AuthMiddleware)

	// TODO: Setup CORS

	s := &Server{
		config:     config,
		router:     router,
		ruleSvc:    ruleSvc,
		triggerSvc: triggerSvc,
		actionSvc:  actionSvc,
	}

	s.setupRoutes()

	addr := ":" + config.Port
	s.server = &http.Server{
		Addr:    addr,
		Handler: router,
	}

	return s
}

// NewServerWithoutRateLimit creates a new HTTP server without rate limiting (for performance testing)
func NewServerWithoutRateLimit(config *ServerConfig, ruleSvc *rule.Service, triggerSvc *trigger.Service, actionSvc *action.Service) *Server {
	router := mux.NewRouter()

	// Add middleware (excluding rate limiting)
	router.Use(LoggingMiddleware)
	router.Use(AuthMiddleware)

	// TODO: Setup CORS

	s := &Server{
		config:     config,
		router:     router,
		ruleSvc:    ruleSvc,
		triggerSvc: triggerSvc,
		actionSvc:  actionSvc,
	}

	s.setupRoutes()

	addr := ":" + config.Port
	s.server = &http.Server{
		Addr:    addr,
		Handler: router,
	}

	return s
}

// setupRoutes registers all API routes
func (s *Server) setupRoutes() {
	// Health check endpoint
	s.router.HandleFunc("/health", s.HealthCheck).Methods("GET")

	// Metrics endpoint
	s.router.Handle("/metrics", promhttp.Handler())

	api := s.router.PathPrefix("/api/v1").Subrouter()

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
