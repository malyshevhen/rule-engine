package api

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/malyshevhen/rule-engine/internal/action"
)

func createAction(actionSvc ActionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateActionRequest
		if err := ValidateAndParseJSON(r, &req); err != nil {
			ErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// Sanitize inputs
		req.LuaScript = strings.TrimSpace(req.LuaScript)

		enabled := true
		if req.Enabled != nil {
			enabled = *req.Enabled
		}

		action := &action.Action{
			LuaScript: req.LuaScript,
			Enabled:   enabled,
		}

		if err := actionSvc.Create(r.Context(), action); err != nil {
			slog.Error("Failed to create action", "error", err)
			ErrorResponse(w, http.StatusInternalServerError, "Failed to create action")
			return
		}

		CreatedResponse(w, ActionToActionInfo(action))
	}
}

func listActions(actionSvc ActionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		actions, err := actionSvc.List(r.Context())
		if err != nil {
			ErrorResponse(w, http.StatusInternalServerError, "Failed to list actions")
			return
		}

		// Convert to DTOs
		actionInfos := make([]ActionInfo, len(actions))
		for i, a := range actions {
			actionInfos[i] = *ActionToActionInfo(a)
		}

		SuccessResponse(w, actionInfos)
	}
}

func getAction(actionSvc ActionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]
		id, err := uuid.Parse(idStr)
		if err != nil {
			ErrorResponse(w, http.StatusBadRequest, "Invalid action ID")
			return
		}

		action, err := actionSvc.GetByID(r.Context(), id)
		if err != nil {
			ErrorResponse(w, http.StatusNotFound, "Action not found")
			return
		}

		SuccessResponse(w, ActionToActionInfo(action))
	}
}
