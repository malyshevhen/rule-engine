package api

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	actionStorage "github.com/malyshevhen/rule-engine/internal/storage/action"
	ruleStorage "github.com/malyshevhen/rule-engine/internal/storage/rule"
)

func addActionToRule(ruleSvc RuleService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ruleIDStr := vars["id"]
		ruleID, err := uuid.Parse(ruleIDStr)
		if err != nil {
			slog.Error("Invalid rule ID format for add action", "rule_id", ruleIDStr, "error", err)
			ErrorResponse(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid rule ID format")
			return
		}

		var req AddActionToRuleRequest
		if err := ParseJSONBody(r, &req); err != nil {
			slog.Error("Failed to parse add action to rule request body", "rule_id", ruleID, "error", err)
			ErrorResponse(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body")
			return
		}

		// Validate action ID
		if req.ActionID == uuid.Nil {
			ErrorResponse(w, http.StatusBadRequest, "VALIDATION_ERROR", "Action ID cannot be empty")
			return
		}

		if err := ruleSvc.AddAction(r.Context(), ruleID, req.ActionID); err != nil {
			// Check for specific errors
			if errors.Is(err, ruleStorage.ErrNotFound) {
				ErrorResponse(w, http.StatusNotFound, "NOT_FOUND", "Rule not found")
				return
			}
			if errors.Is(err, actionStorage.ErrNotFound) {
				ErrorResponse(w, http.StatusBadRequest, "VALIDATION_ERROR", "Action not found")
				return
			}
			slog.Error("Failed to add action to rule", "rule_id", ruleID, "action_id", req.ActionID, "error", err)
			ErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to add action to rule")
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
