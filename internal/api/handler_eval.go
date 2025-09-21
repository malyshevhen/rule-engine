package api

import (
	"net/http"
	"strings"

	execCtx "github.com/malyshevhen/rule-engine/internal/engine/executor/context"
)

func evaluateScript(executorSvc ExecutorService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req EvaluateScriptRequest
		if err := ValidateAndParseJSON(r, &req); err != nil {
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
