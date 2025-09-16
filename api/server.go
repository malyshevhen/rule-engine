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

// ServerConfig holds server configuration
type ServerConfig struct {
	Port string
}

// Server represents the HTTP server
type Server struct {
	config     *ServerConfig
	router     *mux.Router
	server     *http.Server
	ruleSvc    *rule.Service
	triggerSvc *trigger.Service
	actionSvc  *action.Service
}

// NewServer creates a new HTTP server
func NewServer(config *ServerConfig, ruleSvc *rule.Service, triggerSvc *trigger.Service, actionSvc *action.Service) *Server {
	router := mux.NewRouter()

	// Add middleware
	router.Use(RateLimitMiddleware)
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

// TODO: Implement handler methods
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

	rule := &rule.Rule{
		Name:      req.Name,
		LuaScript: req.LuaScript,
		Enabled:   enabled,
	}

	if err := s.ruleSvc.Create(r.Context(), rule); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Failed to create rule")
		return
	}

	SuccessResponse(w, rule)
}

func (s *Server) ListRules(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement list
	SuccessResponse(w, []rule.Rule{})
}
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
func (s *Server) UpdateRule(w http.ResponseWriter, r *http.Request) {}
func (s *Server) DeleteRule(w http.ResponseWriter, r *http.Request) {}

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
	if req.Type != "conditional" && req.Type != "scheduled" {
		ErrorResponse(w, http.StatusBadRequest, "Invalid trigger type (must be 'conditional' or 'scheduled')")
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

func (s *Server) ListTriggers(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement list
	SuccessResponse(w, []trigger.Trigger{})
}
func (s *Server) GetTrigger(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement GetByID in trigger service
	ErrorResponse(w, http.StatusNotImplemented, "Get trigger not implemented")
}

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

func (s *Server) ListActions(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement list
	SuccessResponse(w, []action.Action{})
}
func (s *Server) GetAction(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement GetByID in action service
	ErrorResponse(w, http.StatusNotImplemented, "Get action not implemented")
}
