package api

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/malyshevhen/rule-engine/internal/core/action"
	"github.com/malyshevhen/rule-engine/internal/core/rule"
	"github.com/malyshevhen/rule-engine/internal/core/trigger"
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

func (s *Server) CreateTrigger(w http.ResponseWriter, r *http.Request) {}
func (s *Server) ListTriggers(w http.ResponseWriter, r *http.Request)  {}
func (s *Server) GetTrigger(w http.ResponseWriter, r *http.Request)    {}

func (s *Server) CreateAction(w http.ResponseWriter, r *http.Request) {}
func (s *Server) ListActions(w http.ResponseWriter, r *http.Request)  {}
func (s *Server) GetAction(w http.ResponseWriter, r *http.Request)    {}
