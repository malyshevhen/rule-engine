package api

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func addActionToRule(ruleSvc RuleService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ruleIDStr := vars["id"]
		ruleID, err := uuid.Parse(ruleIDStr)
		if err != nil {
			ErrorResponse(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid rule ID")
			return
		}

		var req AddActionToRuleRequest
		if err := ParseJSONBody(r, &req); err != nil {
			ErrorResponse(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body")
			return
		}

		if err := ruleSvc.AddAction(r.Context(), ruleID, req.ActionID); err != nil {
			ErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to add action to rule")
			return
		}

		SuccessResponse(w, map[string]string{"status": "success"})
	}
}
