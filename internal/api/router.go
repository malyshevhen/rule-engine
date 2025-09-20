package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// setupRoutes registers all API routes
func setupRoutes(
	executorSvc ExecutorService,
	ruleSvc RuleService,
	triggerSvc TriggerService,
	actionSvc ActionService,
) *mux.Router {
	router := mux.NewRouter()
	// Public routes (no authentication required)
	router.HandleFunc("/health", HealthCheck).Methods("GET")
	router.Handle("/metrics", promhttp.Handler())

	// Protected API routes (require authentication)
	api := router.PathPrefix("/api/v1").Subrouter()
	api.Use(AuthMiddleware)

	// Rules routes
	api.HandleFunc("/rules", createRule(ruleSvc)).Methods("POST")
	api.HandleFunc("/rules", listRules(ruleSvc)).Methods("GET")
	api.HandleFunc("/rules/{id}", getRule(ruleSvc)).Methods("GET")
	api.HandleFunc("/rules/{id}", updateRule(ruleSvc)).Methods("PUT")
	api.HandleFunc("/rules/{id}", deleteRule(ruleSvc)).Methods("DELETE")

	// Triggers routes
	api.HandleFunc("/triggers", createTrigger(triggerSvc)).Methods("POST")
	api.HandleFunc("/triggers", listTriggers(triggerSvc)).Methods("GET")
	api.HandleFunc("/triggers/{id}", getTrigger(triggerSvc)).Methods("GET")

	// Actions routes
	api.HandleFunc("/actions", createAction(actionSvc)).Methods("POST")
	api.HandleFunc("/actions", listActions(actionSvc)).Methods("GET")
	api.HandleFunc("/actions/{id}", getAction(actionSvc)).Methods("GET")

	// Script evaluation route
	api.HandleFunc("/evaluate", evaluateScript(executorSvc)).Methods("POST")

	return router
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("healthy"))
}
