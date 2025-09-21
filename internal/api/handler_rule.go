package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/malyshevhen/rule-engine/internal/rule"
	ruleStorage "github.com/malyshevhen/rule-engine/internal/storage/rule"
)

func createRule(ruleSvc RuleService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateRuleRequest
		if err := ValidateAndParseJSON(r, &req); err != nil {
			ErrorResponse(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
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
			ErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create rule")
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
				ErrorResponse(w, http.StatusBadRequest, "VALIDATION_ERROR", fmt.Sprintf("Invalid limit parameter (must be between 1 and %d)", apiConfig.MaxRulesLimit))
				return
			}
		}

		if offsetStr != "" {
			if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
				offset = parsedOffset
			} else {
				ErrorResponse(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid offset parameter (must be non-negative)")
				return
			}
		}

		rules, total, err := ruleSvc.List(r.Context(), limit, offset)
		if err != nil {
			ErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to list rules")
			return
		}

		// Create response with pagination metadata
		response := map[string]any{
			"rules":  RulesToRuleInfos(rules),
			"limit":  limit,
			"offset": offset,
			"count":  len(rules),
			"total":  total,
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
			ErrorResponse(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid rule ID")
			return
		}

		rule, err := ruleSvc.GetByID(r.Context(), id)
		if err != nil {
			ErrorResponse(w, http.StatusNotFound, "NOT_FOUND", "Rule not found")
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
			ErrorResponse(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid rule ID")
			return
		}

		// Parse JSON Patch operations
		var patchOps PatchRequest
		if err := ParseJSONBody(r, &patchOps); err != nil {
			ErrorResponse(w, http.StatusBadRequest, "VALIDATION_ERROR", fmt.Sprintf("Invalid JSON patch request body: %s", err.Error()))
			return
		}

		// Validate patch operations
		if len(patchOps) == 0 {
			ErrorResponse(w, http.StatusBadRequest, "VALIDATION_ERROR", "At least one patch operation is required")
			return
		}

		// Basic validation of patch operations
		for i, op := range patchOps {
			if strings.TrimSpace(op.Path) == "" {
				ErrorResponse(w, http.StatusBadRequest, "VALIDATION_ERROR", fmt.Sprintf("Patch operation %d: path cannot be empty", i+1))
				return
			}
			// Validate path format (should start with /)
			if !strings.HasPrefix(op.Path, "/") {
				ErrorResponse(w, http.StatusBadRequest, "VALIDATION_ERROR", fmt.Sprintf("Patch operation %d: path must start with '/'", i+1))
				return
			}
		}

		// Convert patch operations to JSON Patch format
		patchJSON, err := json.Marshal(patchOps)
		if err != nil {
			ErrorResponse(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid patch operations")
			return
		}

		// Apply the patch
		patch, err := jsonpatch.DecodePatch(patchJSON)
		if err != nil {
			ErrorResponse(w, http.StatusBadRequest, "VALIDATION_ERROR", fmt.Sprintf("Invalid JSON Patch format: %s", err.Error()))
			return
		}

		// Get existing rule
		existingRule, err := ruleSvc.GetByID(r.Context(), id)
		if err != nil {
			ErrorResponse(w, http.StatusNotFound, "NOT_FOUND", "Rule not found")
			return
		}

		// Convert rule to JSON for patching
		ruleJSON, err := json.Marshal(RuleToRuleInfo(existingRule))
		if err != nil {
			ErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to prepare rule for patching")
			return
		}

		modifiedJSON, err := patch.Apply(ruleJSON)
		if err != nil {
			// Provide more specific error messages based on the type of patch error
			var errorMsg string
			if strings.Contains(err.Error(), "path") {
				errorMsg = fmt.Sprintf("Invalid patch path: %s", err.Error())
			} else if strings.Contains(err.Error(), "operation") {
				errorMsg = fmt.Sprintf("Invalid patch operation: %s", err.Error())
			} else {
				errorMsg = fmt.Sprintf("Patch application failed: %s", err.Error())
			}
			ErrorResponse(w, http.StatusBadRequest, "VALIDATION_ERROR", errorMsg)
			return
		}

		// Convert back to RuleInfo
		var updatedRuleInfo RuleInfo
		if err := json.Unmarshal(modifiedJSON, &updatedRuleInfo); err != nil {
			ErrorResponse(w, http.StatusBadRequest, "VALIDATION_ERROR", fmt.Sprintf("Patch resulted in invalid JSON structure: %s", err.Error()))
			return
		}

		// Validate the updated rule
		var validationErrors []string

		trimmedName := strings.TrimSpace(updatedRuleInfo.Name)
		if trimmedName == "" {
			validationErrors = append(validationErrors, "rule name cannot be empty")
		} else if len(trimmedName) > apiConfig.MaxRuleNameLength {
			validationErrors = append(validationErrors, fmt.Sprintf("rule name too long (max %d characters, got %d)", apiConfig.MaxRuleNameLength, len(trimmedName)))
		}

		trimmedScript := strings.TrimSpace(updatedRuleInfo.LuaScript)
		if trimmedScript == "" {
			validationErrors = append(validationErrors, "Lua script cannot be empty")
		} else if len(trimmedScript) > apiConfig.MaxLuaScriptLength {
			validationErrors = append(validationErrors, fmt.Sprintf("Lua script too long (max %d characters, got %d)", apiConfig.MaxLuaScriptLength, len(trimmedScript)))
		}

		if len(validationErrors) > 0 {
			ErrorResponse(w, http.StatusBadRequest, "VALIDATION_ERROR", fmt.Sprintf("Validation failed: %s", strings.Join(validationErrors, "; ")))
			return
		}

		// Convert back to domain model
		updatedRule := &rule.Rule{
			ID:        existingRule.ID,
			Name:      strings.TrimSpace(updatedRuleInfo.Name),
			LuaScript: strings.TrimSpace(updatedRuleInfo.LuaScript),
			Priority:  updatedRuleInfo.Priority,
			Enabled:   updatedRuleInfo.Enabled,
			CreatedAt: existingRule.CreatedAt,
			UpdatedAt: existingRule.UpdatedAt,
		}

		if err := ruleSvc.Update(r.Context(), updatedRule); err != nil {
			ErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", fmt.Sprintf("Failed to update rule in database: %s", err.Error()))
			return
		}

		SuccessResponse(w, RuleToRuleInfo(updatedRule))
	}
}

func deleteRule(ruleSvc RuleService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]
		id, err := uuid.Parse(idStr)
		if err != nil {
			ErrorResponse(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid rule ID")
			return
		}

		if err := ruleSvc.Delete(r.Context(), id); err != nil {
			if errors.Is(err, ruleStorage.ErrNotFound) {
				ErrorResponse(w, http.StatusNotFound, "NOT_FOUND", "Rule not found")
				return
			}
			ErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to delete rule")
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
