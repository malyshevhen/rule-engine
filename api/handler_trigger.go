package api

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/malyshevhen/rule-engine/internal/core/trigger"
)

func createTrigger(triggerSvc TriggerService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		if err := triggerSvc.Create(r.Context(), trigger); err != nil {
			ErrorResponse(w, http.StatusInternalServerError, "Failed to create trigger")
			return
		}

		SuccessResponse(w, trigger)
	}
}

func listTriggers(triggerSvc TriggerService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		triggers, err := triggerSvc.List(r.Context())
		if err != nil {
			ErrorResponse(w, http.StatusInternalServerError, "Failed to list triggers")
			return
		}

		SuccessResponse(w, triggers)
	}
}

func getTrigger(triggerSvc TriggerService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]
		id, err := uuid.Parse(idStr)
		if err != nil {
			ErrorResponse(w, http.StatusBadRequest, "Invalid trigger ID")
			return
		}

		trigger, err := triggerSvc.GetByID(r.Context(), id)
		if err != nil {
			ErrorResponse(w, http.StatusNotFound, "Trigger not found")
			return
		}

		SuccessResponse(w, trigger)
	}
}
