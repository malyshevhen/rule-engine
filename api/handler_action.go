package api

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/malyshevhen/rule-engine/internal/core/action"
)

func createAction(actionSvc ActionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateActionRequest
		if err := ParseJSONBody(r, &req); err != nil {
			ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		// Sanitize and validate inputs
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

		action := &action.Action{
			LuaScript: req.LuaScript,
			Enabled:   enabled,
		}

		if err := actionSvc.Create(r.Context(), action); err != nil {
			slog.Error("Failed to create action", "error", err)
			ErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		CreatedResponse(w, action)
	}
}

func listActions(actionSvc ActionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		actions, err := actionSvc.List(r.Context())
		if err != nil {
			ErrorResponse(w, http.StatusInternalServerError, "Failed to list actions")
			return
		}

		SuccessResponse(w, actions)
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

		SuccessResponse(w, action)
	}
}
