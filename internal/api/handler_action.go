package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/malyshevhen/rule-engine/internal/action"
	actionStorage "github.com/malyshevhen/rule-engine/internal/storage/action"
)

// createAction creates a new action
//
//	@Summary		Create a new action
//	@Description	Create a new action with the given name and Lua script.
//	@Tags			actions
//	@Accept			json
//	@Produce		json
//	@Param			action	body		CreateActionRequest	true	"Action to create"
//	@Success		201		{object}	ActionInfo
//	@Failure		400		{object}	APIErrorResponse
//	@Failure		500		{object}	APIErrorResponse
//	@Router			/api/v1/actions [post]
func createAction(actionSvc ActionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateActionRequest
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

		action := &action.Action{
			Name:      req.Name,
			LuaScript: req.LuaScript,
			Enabled:   enabled,
		}

		if err := actionSvc.Create(r.Context(), action); err != nil {
			slog.Error("Failed to create action", "error", err)
			ErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create action")
			return
		}

		CreatedResponse(w, ActionToActionInfo(action))
	}
}

// listActions lists all existing actions
//
//	@Summary		List all actions
//	@Description	Get a list of all actions with optional pagination.
//	@Tags			actions
//	@Produce		json
//	@Param			limit	query		int	false	"Limit number of actions returned"
//	@Param			offset	query		int	false	"Offset for pagination"
//	@Success		200		{object}	map[string]any
//	@Failure		400		{object}	APIErrorResponse
//	@Failure		500		{object}	APIErrorResponse
//	@Router			/api/v1/actions [get]
func listActions(actionSvc ActionService) http.HandlerFunc {
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

		actions, total, err := actionSvc.List(r.Context(), limit, offset)
		if err != nil {
			ErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to list actions")
			return
		}

		// Convert to DTOs
		actionInfos := make([]ActionInfo, len(actions))
		for i, a := range actions {
			actionInfos[i] = *ActionToActionInfo(a)
		}

		// Create response with pagination metadata
		response := map[string]any{
			"actions": actionInfos,
			"limit":   limit,
			"offset":  offset,
			"count":   len(actionInfos),
			"total":   total,
		}

		SuccessResponse(w, response)
	}
}

// getAction gets an action by its ID
//
//	@Summary		Get an action by ID
//	@Description	Get a single action by its unique ID.
//	@Tags			actions
//	@Produce		json
//	@Param			id	path		string	true	"Action ID"
//	@Success		200	{object}	ActionInfo
//	@Failure		400	{object}	APIErrorResponse
//	@Failure		404	{object}	APIErrorResponse
//	@Failure		500	{object}	APIErrorResponse
//	@Router			/api/v1/actions/{id} [get]
func getAction(actionSvc ActionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]
		id, err := uuid.Parse(idStr)
		if err != nil {
			slog.Error("Invalid action ID format", "id", idStr, "error", err)
			ErrorResponse(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid action ID format")
			return
		}

		action, err := actionSvc.GetByID(r.Context(), id)
		if err != nil {
			if errors.Is(err, actionStorage.ErrNotFound) {
				ErrorResponse(w, http.StatusNotFound, "NOT_FOUND", "Action not found")
				return
			}
			slog.Error("Failed to get action", "action_id", id, "error", err)
			ErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve action")
			return
		}

		SuccessResponse(w, ActionToActionInfo(action))
	}
}

// updateAction updates an action
//
//	@Summary		Update an action
//	@Description	Update an existing action using a JSON Patch.
//	@Tags			actions
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string			true	"Action ID"
//	@Param			patch	body		PatchRequest	true	"JSON Patch operations"
//	@Success		200		{object}	ActionInfo
//	@Failure		400		{object}	APIErrorResponse
//	@Failure		404		{object}	APIErrorResponse
//	@Failure		500		{object}	APIErrorResponse
//	@Router			/api/v1/actions/{id} [patch]
func updateAction(actionSvc ActionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]
		id, err := uuid.Parse(idStr)
		if err != nil {
			slog.Error("Invalid action ID format for update", "id", idStr, "error", err)
			ErrorResponse(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid action ID format")
			return
		}

		// Get the current action
		currentAction, err := actionSvc.GetByID(r.Context(), id)
		if err != nil {
			if errors.Is(err, actionStorage.ErrNotFound) {
				ErrorResponse(w, http.StatusNotFound, "NOT_FOUND", "Action not found")
				return
			}
			slog.Error("Failed to get action for update", "action_id", id, "error", err)
			ErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve action")
			return
		}

		// Apply JSON Patch
		actionJSON, err := json.Marshal(currentAction)
		if err != nil {
			ErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to serialize action")
			return
		}

		modifiedJSON, err := ApplyJSONPatch(r, actionJSON, "action", id.String())
		if err != nil {
			ErrorResponse(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
			return
		}

		var updatedAction action.Action
		if err := json.Unmarshal(modifiedJSON, &updatedAction); err != nil {
			ErrorResponse(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid patch result")
			return
		}

		// Validate the updated action
		if updatedAction.LuaScript == "" {
			ErrorResponse(w, http.StatusBadRequest, "VALIDATION_ERROR", "lua_script cannot be empty")
			return
		}

		// Ensure ID is preserved
		updatedAction.ID = id

		// Update the action
		if err := actionSvc.Update(r.Context(), &updatedAction); err != nil {
			slog.Error("Failed to update action", "action_id", id, "error", err)
			ErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update action")
			return
		}

		// Return the updated action
		SuccessResponse(w, ActionToActionInfo(&updatedAction))
	}
}

func deleteAction(actionSvc ActionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]
		id, err := uuid.Parse(idStr)
		if err != nil {
			slog.Error("Invalid action ID format for delete", "id", idStr, "error", err)
			ErrorResponse(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid action ID format")
			return
		}

		if err := actionSvc.Delete(r.Context(), id); err != nil {
			if errors.Is(err, actionStorage.ErrNotFound) {
				ErrorResponse(w, http.StatusNotFound, "NOT_FOUND", "Action not found")
				return
			}
			slog.Error("Failed to delete action", "action_id", id, "error", err)
			ErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to delete action")
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
