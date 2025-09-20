package api

import (
	"net/http"
	"strings"

	execCtx "github.com/malyshevhen/rule-engine/internal/engine/executor/context"
)

func evaluateScript(executorSvc ExecutorService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req EvaluateScriptRequest
		if err := ParseJSONBody(r, &req); err != nil {
			ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		// Sanitize and validate inputs
		req.Script = strings.TrimSpace(req.Script)
		if req.Script == "" {
			ErrorResponse(w, http.StatusBadRequest, "Script cannot be empty")
			return
		}
		if len(req.Script) > 10000 {
			ErrorResponse(w, http.StatusBadRequest, "Script too long (max 10000 characters)")
			return
		}

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
