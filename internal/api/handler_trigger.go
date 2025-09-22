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
	triggerStorage "github.com/malyshevhen/rule-engine/internal/storage/trigger"
	"github.com/malyshevhen/rule-engine/internal/trigger"
)

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
