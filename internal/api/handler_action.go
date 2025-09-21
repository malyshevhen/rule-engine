package api

import (
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

func getAction(actionSvc ActionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]
		id, err := uuid.Parse(idStr)
		if err != nil {
			ErrorResponse(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid action ID")
			return
		}

		action, err := actionSvc.GetByID(r.Context(), id)
		if err != nil {
			ErrorResponse(w, http.StatusNotFound, "NOT_FOUND", "Action not found")
			return
		}

		SuccessResponse(w, ActionToActionInfo(action))
	}
}

func deleteAction(actionSvc ActionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]
		id, err := uuid.Parse(idStr)
		if err != nil {
			ErrorResponse(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid action ID")
			return
		}

		if err := actionSvc.Delete(r.Context(), id); err != nil {
			if errors.Is(err, actionStorage.ErrNotFound) {
				ErrorResponse(w, http.StatusNotFound, "NOT_FOUND", "Action not found")
				return
			}
			ErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to delete action")
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
