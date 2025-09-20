package api

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/malyshevhen/rule-engine/internal/trigger"
)

func createTrigger(triggerSvc TriggerService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateTriggerRequest
		if err := ValidateAndParseJSON(r, &req); err != nil {
			ErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// Sanitize inputs
		req.Type = strings.TrimSpace(req.Type)
		req.ConditionScript = strings.TrimSpace(req.ConditionScript)

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

		CreatedResponse(w, TriggerToTriggerInfo(trigger))
	}
}

func listTriggers(triggerSvc TriggerService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		triggers, err := triggerSvc.List(r.Context())
		if err != nil {
			ErrorResponse(w, http.StatusInternalServerError, "Failed to list triggers")
			return
		}

		// Convert to DTOs
		triggerInfos := make([]TriggerInfo, len(triggers))
		for i, t := range triggers {
			triggerInfos[i] = *TriggerToTriggerInfo(t)
		}

		SuccessResponse(w, triggerInfos)
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

		SuccessResponse(w, TriggerToTriggerInfo(trigger))
	}
}
