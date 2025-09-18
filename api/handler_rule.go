package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/malyshevhen/rule-engine/internal/core/rule"
	ruleStorage "github.com/malyshevhen/rule-engine/internal/storage/rule"
)

func createRule(ruleSvc RuleService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		if err := ruleSvc.Create(r.Context(), rule); err != nil {
			ErrorResponse(w, http.StatusInternalServerError, "Failed to create rule")
			return
		}

		SuccessResponse(w, rule)
	}
}

func listRules(ruleSvc RuleService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rules, err := ruleSvc.ListAll(r.Context())
		if err != nil {
			ErrorResponse(w, http.StatusInternalServerError, "Failed to list rules")
			return
		}

		SuccessResponse(w, rules)
	}
}

func getRule(ruleSvc RuleService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]
		id, err := uuid.Parse(idStr)
		if err != nil {
			ErrorResponse(w, http.StatusBadRequest, "Invalid rule ID")
			return
		}

		rule, err := ruleSvc.GetByID(r.Context(), id)
		if err != nil {
			ErrorResponse(w, http.StatusNotFound, "Rule not found")
			return
		}

		SuccessResponse(w, rule)
	}
}

func updateRule(ruleSvc RuleService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		existingRule, err := ruleSvc.GetByID(r.Context(), id)
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

		if err := ruleSvc.Update(r.Context(), existingRule); err != nil {
			ErrorResponse(w, http.StatusInternalServerError, "Failed to update rule")
			return
		}

		SuccessResponse(w, existingRule)
	}
}

func deleteRule(ruleSvc RuleService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]
		id, err := uuid.Parse(idStr)
		if err != nil {
			ErrorResponse(w, http.StatusBadRequest, "Invalid rule ID")
			return
		}

		if err := ruleSvc.Delete(r.Context(), id); err != nil {
			if errors.Is(err, ruleStorage.ErrNotFound) {
				ErrorResponse(w, http.StatusNotFound, "Rule not found")
				return
			}
			ErrorResponse(w, http.StatusInternalServerError, "Failed to delete rule")
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
