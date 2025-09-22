package api

import (
	"log/slog"
	"net/http"
	"strings"

	execCtx "github.com/malyshevhen/rule-engine/internal/engine/executor/context"
)

// @Summary Evaluate a Lua script
// @Description Evaluate a Lua script in a sandboxed environment.
// @Tags evaluation
// @Accept  json
// @Produce  json
// @Param   script  body      EvaluateScriptRequest   true  "Script to evaluate"
// @Success 200     {object}  EvaluateScriptResponse
// @Failure 400     {object}  APIErrorResponse
// @Failure 500     {object}  APIErrorResponse
// @Router /api/v1/evaluate [post]
func evaluateScript(executorSvc ExecutorService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req EvaluateScriptRequest
		if err := ValidateAndParseJSON(r, &req); err != nil {
			slog.Error("Failed to validate evaluate script request", "error", err)
			ErrorResponse(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
			return
		}

		// Sanitize inputs
		req.Script = strings.TrimSpace(req.Script)

		// Create a minimal execution context for evaluation
		execContext := &execCtx.ExecutionContext{
			RuleID:    "evaluate", // Use a fixed ID for evaluation
			TriggerID: "evaluate",
			Data:      req.Context, // Empty data for evaluation
		}

		// Execute the script
		result := executorSvc.ExecuteScript(r.Context(), req.Script, execContext)

		// Convert duration to string
		durationStr := result.Duration.String()

		response := EvaluateScriptResponse{
			Success:  result.Success,
			Output:   result.Output,
			Error:    result.Error,
			Duration: durationStr,
		}

		SuccessResponse(w, response)
	}
}
