package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/malyshevhen/rule-engine/internal/rule"
	ruleStorage "github.com/malyshevhen/rule-engine/internal/storage/rule"
)

func createRule(ruleSvc RuleService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateRuleRequest
		if err := ValidateAndParseJSON(r, &req); err != nil {
			ErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// Sanitize inputs
		req.Name = strings.TrimSpace(req.Name)
		req.LuaScript = strings.TrimSpace(req.LuaScript)

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

		CreatedResponse(w, RuleToRuleInfo(rule))
	}
}

func listRules(ruleSvc RuleService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse pagination parameters
		limitStr := GetQueryParam(r, "limit")
		offsetStr := GetQueryParam(r, "offset")

		limit := apiConfig.DefaultRulesLimit
		offset := apiConfig.DefaultRulesOffset

		if limitStr != "" {
			if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 && parsedLimit <= apiConfig.MaxRulesLimit {
				limit = parsedLimit
			} else {
				ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("Invalid limit parameter (must be between 1 and %d)", apiConfig.MaxRulesLimit))
				return
			}
		}

		if offsetStr != "" {
			if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
				offset = parsedOffset
			} else {
				ErrorResponse(w, http.StatusBadRequest, "Invalid offset parameter (must be non-negative)")
				return
			}
		}

		rules, err := ruleSvc.List(r.Context(), limit, offset)
		if err != nil {
			ErrorResponse(w, http.StatusInternalServerError, "Failed to list rules")
			return
		}

		// Create response with pagination metadata
		response := map[string]interface{}{
			"rules":  RulesToRuleInfos(rules),
			"limit":  limit,
			"offset": offset,
			"count":  len(rules),
		}

		SuccessResponse(w, response)
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

		SuccessResponse(w, RuleToRuleInfo(rule))
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
		if err := ValidateAndParseJSON(r, &req); err != nil {
			ErrorResponse(w, http.StatusBadRequest, err.Error())
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
			existingRule.Name = name
		}

		if req.LuaScript != nil {
			script := strings.TrimSpace(*req.LuaScript)
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

		SuccessResponse(w, RuleToRuleInfo(existingRule))
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
