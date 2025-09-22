package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/evanphx/json-patch/v5"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	triggerStorage "github.com/malyshevhen/rule-engine/internal/storage/trigger"
	"github.com/malyshevhen/rule-engine/internal/trigger"
)

// @Summary Create a new trigger
// @Description Create a new trigger for a rule.
// @Tags triggers
// @Accept  json
// @Produce  json
// @Param   trigger body      CreateTriggerRequest  true  "Trigger to create"
// @Success 201     {object}  TriggerInfo
// @Failure 400     {object}  APIErrorResponse
// @Failure 500     {object}  APIErrorResponse
// @Router /api/v1/triggers [post]
func createTrigger(triggerSvc TriggerService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateTriggerRequest
		if err := ValidateAndParseJSON(r, &req); err != nil {
			ErrorResponse(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
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
			ErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create trigger")
			return
		}

		CreatedResponse(w, TriggerToTriggerInfo(trigger))
	}
}

// @Summary List all triggers
// @Description Get a list of all triggers with optional pagination.
// @Tags triggers
// @Produce  json
// @Param   limit   query     int     false  "Limit number of triggers returned"
// @Param   offset  query     int     false  "Offset for pagination"
// @Success 200     {object}  map[string]any
// @Failure 400     {object}  APIErrorResponse
// @Failure 500     {object}  APIErrorResponse
// @Router /api/v1/triggers [get]
func listTriggers(triggerSvc TriggerService) http.HandlerFunc {
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

		triggers, total, err := triggerSvc.List(r.Context(), limit, offset)
		if err != nil {
			ErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to list triggers")
			return
		}

		// Convert to DTOs
		triggerInfos := make([]TriggerInfo, len(triggers))
		for i, t := range triggers {
			triggerInfos[i] = *TriggerToTriggerInfo(t)
		}

		// Create response with pagination metadata
		response := map[string]any{
			"triggers": triggerInfos,
			"limit":    limit,
			"offset":   offset,
			"count":    len(triggerInfos),
			"total":    total,
		}

		SuccessResponse(w, response)
	}
}

// @Summary Get a trigger by ID
// @Description Get a single trigger by its unique ID.
// @Tags triggers
// @Produce  json
// @Param   id      path      string  true  "Trigger ID"
// @Success 200     {object}  TriggerInfo
// @Failure 400     {object}  APIErrorResponse
// @Failure 404     {object}  APIErrorResponse
// @Failure 500     {object}  APIErrorResponse
// @Router /api/v1/triggers/{id} [get]
func getTrigger(triggerSvc TriggerService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]
		id, err := uuid.Parse(idStr)
		if err != nil {
			slog.Error("Invalid trigger ID format", "id", idStr, "error", err)
			ErrorResponse(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid trigger ID format")
			return
		}

		trigger, err := triggerSvc.GetByID(r.Context(), id)
		if err != nil {
			if errors.Is(err, triggerStorage.ErrNotFound) {
				ErrorResponse(w, http.StatusNotFound, "NOT_FOUND", "Trigger not found")
				return
			}
			slog.Error("Failed to get trigger", "trigger_id", id, "error", err)
			ErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve trigger")
			return
		}

		SuccessResponse(w, TriggerToTriggerInfo(trigger))
	}
}

// @Summary Update a trigger
// @Description Update an existing trigger using a JSON Patch.
// @Tags triggers
// @Accept  json
// @Produce  json
// @Param   id      path      string        true  "Trigger ID"
// @Param   patch   body      PatchRequest  true  "JSON Patch operations"
// @Success 200     {object}  TriggerInfo
// @Failure 400     {object}  APIErrorResponse
// @Failure 404     {object}  APIErrorResponse
// @Failure 500     {object}  APIErrorResponse
// @Router /api/v1/triggers/{id} [patch]
func updateTrigger(triggerSvc TriggerService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]
		id, err := uuid.Parse(idStr)
		if err != nil {
			slog.Error("Invalid trigger ID format for update", "id", idStr, "error", err)
			ErrorResponse(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid trigger ID format")
			return
		}

		// Parse JSON Patch operations
		var patchOps PatchRequest
		if err := ParseJSONBody(r, &patchOps); err != nil {
			slog.Error("Failed to parse JSON patch body", "trigger_id", id, "error", err)
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

		// Get the current trigger
		currentTrigger, err := triggerSvc.GetByID(r.Context(), id)
		if err != nil {
			if errors.Is(err, triggerStorage.ErrNotFound) {
				ErrorResponse(w, http.StatusNotFound, "NOT_FOUND", "Trigger not found")
				return
			}
			slog.Error("Failed to get trigger for update", "trigger_id", id, "error", err)
			ErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve trigger")
			return
		}

		// Apply JSON Patch
		triggerJSON, err := json.Marshal(currentTrigger)
		if err != nil {
			ErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to serialize trigger")
			return
		}

		patch, err := jsonpatch.DecodePatch(patchJSON)
		if err != nil {
			ErrorResponse(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid patch operations")
			return
		}

		modifiedJSON, err := patch.Apply(triggerJSON)
		if err != nil {
			ErrorResponse(w, http.StatusBadRequest, "VALIDATION_ERROR", fmt.Sprintf("Failed to apply patch: %s", err.Error()))
			return
		}

		var updatedTrigger trigger.Trigger
		if err := json.Unmarshal(modifiedJSON, &updatedTrigger); err != nil {
			ErrorResponse(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid patch result")
			return
		}

		// Validate the updated trigger
		if updatedTrigger.ConditionScript == "" {
			ErrorResponse(w, http.StatusBadRequest, "VALIDATION_ERROR", "condition_script cannot be empty")
			return
		}

		// Ensure ID and RuleID are preserved
		updatedTrigger.ID = id
		updatedTrigger.RuleID = currentTrigger.RuleID

		// Update the trigger
		if err := triggerSvc.Update(r.Context(), &updatedTrigger); err != nil {
			slog.Error("Failed to update trigger", "trigger_id", id, "error", err)
			ErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update trigger")
			return
		}

		// Return the updated trigger
		SuccessResponse(w, TriggerToTriggerInfo(&updatedTrigger))
	}
}

// @Summary Delete a trigger
// @Description Delete a trigger by its unique ID.
// @Tags triggers
// @Produce  json
// @Param   id      path      string  true  "Trigger ID"
// @Success 204     {string}  string  "No Content"
// @Failure 400     {object}  APIErrorResponse
// @Failure 404     {object}  APIErrorResponse
// @Failure 500     {object}  APIErrorResponse
// @Router /api/v1/triggers/{id} [delete]
func deleteTrigger(triggerSvc TriggerService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]
		id, err := uuid.Parse(idStr)
		if err != nil {
			slog.Error("Invalid trigger ID format for delete", "id", idStr, "error", err)
			ErrorResponse(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid trigger ID format")
			return
		}

		if err := triggerSvc.Delete(r.Context(), id); err != nil {
			if errors.Is(err, triggerStorage.ErrNotFound) {
				ErrorResponse(w, http.StatusNotFound, "NOT_FOUND", "Trigger not found")
				return
			}
			slog.Error("Failed to delete trigger", "trigger_id", id, "error", err)
			ErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to delete trigger")
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
