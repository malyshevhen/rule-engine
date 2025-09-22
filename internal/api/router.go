package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// setupRoutes registers all API routes
func setupRoutes(
	healthSvc *Health,
	executorSvc ExecutorService,
	ruleSvc RuleService,
	triggerSvc TriggerService,
	actionSvc ActionService,
) *mux.Router {
	router := mux.NewRouter()

	// Add swagger UI and JSON endpoints
	router.HandleFunc("/swagger/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./docs/index.html")
	})
	router.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./docs/swagger.json")
	})

	// Public routes (no authentication required)
	router.HandleFunc("/health", healthSvc.healthCheckHandler()).Methods("GET")
	router.Handle("/metrics", promhttp.Handler())

	// Protected API routes (require authentication)
	api := router.PathPrefix("/api/v1").Subrouter()
	api.Use(AuthMiddleware)

	// Rules routes
	api.HandleFunc("/rules", createRule(ruleSvc)).Methods("POST")
	api.HandleFunc("/rules", listRules(ruleSvc)).Methods("GET")
	api.HandleFunc("/rules/{id}", getRule(ruleSvc)).Methods("GET")
	api.HandleFunc("/rules/{id}", updateRule(ruleSvc)).Methods("PATCH")
	api.HandleFunc("/rules/{id}", deleteRule(ruleSvc)).Methods("DELETE")
	api.HandleFunc("/rules/{id}/actions", addActionToRule(ruleSvc)).Methods("POST")

	// Triggers routes
	api.HandleFunc("/triggers", createTrigger(triggerSvc)).Methods("POST")
	api.HandleFunc("/triggers", listTriggers(triggerSvc)).Methods("GET")
	api.HandleFunc("/triggers/{id}", getTrigger(triggerSvc)).Methods("GET")
	api.HandleFunc("/triggers/{id}", updateTrigger(triggerSvc)).Methods("PATCH")
	api.HandleFunc("/triggers/{id}", deleteTrigger(triggerSvc)).Methods("DELETE")

	// Actions routes
	api.HandleFunc("/actions", createAction(actionSvc)).Methods("POST")
	api.HandleFunc("/actions", listActions(actionSvc)).Methods("GET")
	api.HandleFunc("/actions/{id}", getAction(actionSvc)).Methods("GET")
	api.HandleFunc("/actions/{id}", updateAction(actionSvc)).Methods("PATCH")
	api.HandleFunc("/actions/{id}", deleteAction(actionSvc)).Methods("DELETE")

	// Script evaluation route
	api.HandleFunc("/evaluate", evaluateScript(executorSvc)).Methods("POST")

	return router
}
