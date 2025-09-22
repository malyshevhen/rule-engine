package e2e

import (
	"bytes"
	"context"
	_ "embed"
	"net/http"
	"testing"

	re_client "github.com/malyshevhen/rule-engine/client"
	"github.com/stretchr/testify/require"
)

//go:embed fixtures/evaluate/complex_script.lua
var complexScript string

func TestEvaluateEndpoint(t *testing.T) {
	ctx := context.Background()

	// Verify environment is set up correctly
	require.NotNil(t, testEnv)

	// Create client
	client := re_client.NewClient(testEnv.GetRuleEngineURL(ctx, t), re_client.AuthConfig{
		APIKey: "test-api-key",
	})

	t.Run("SimpleScript", func(t *testing.T) {
		result, err := client.EvaluateScript(ctx, re_client.EvaluateScriptRequest{
			Script:  "return 2 + 3",
			Context: map[string]any{},
		})

		require.NoError(t, err)
		require.True(t, result.Success)
		require.Equal(t, []any{5.0}, result.Output)
		require.Empty(t, result.Error)
		require.NotEmpty(t, result.Duration)
	})

	t.Run("ScriptWithPlatformAPI", func(t *testing.T) {
		result, err := client.EvaluateScript(ctx, re_client.EvaluateScriptRequest{
			Script:  complexScript,
			Context: map[string]any{},
		})

		require.NoError(t, err)
		require.True(t, result.Success)
		require.Equal(t, []any{true}, result.Output)
		require.Empty(t, result.Error)
		require.NotEmpty(t, result.Duration)
	})

	t.Run("ScriptError", func(t *testing.T) {
		result, err := client.EvaluateScript(ctx, re_client.EvaluateScriptRequest{
			Script:  "invalid lua syntax {{{",
			Context: map[string]any{},
		})

		require.NoError(t, err)
		require.False(t, result.Success)
		require.Empty(t, result.Output)
		require.Contains(t, result.Error, "parse error")
		require.NotEmpty(t, result.Duration)
	})

	t.Run("EmptyScript", func(t *testing.T) {
		errorResp, err := client.EvaluateScript(ctx, re_client.EvaluateScriptRequest{
			Script:  "",
			Context: map[string]any{},
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "validation failed: Script is required")
		require.Nil(t, errorResp)
	})

	t.Run("ScriptTooLong", func(t *testing.T) {
		longScript := make([]byte, 10001)
		for i := range longScript {
			longScript[i] = 'a'
		}

		errorResp, err := client.EvaluateScript(ctx, re_client.EvaluateScriptRequest{
			Script:  string(longScript),
			Context: map[string]any{},
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "validation failed: Lua script must be between 1 and 10000 characters")
		require.Nil(t, errorResp)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		// Test without authorization - need to make direct HTTP call for this
		baseURL := testEnv.GetRuleEngineURL(ctx, t)
		req, err := http.NewRequest("POST", baseURL+"/api/v1/evaluate", bytes.NewReader([]byte(`{"script":"return 42"}`)))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		// Missing authorization header
		resp, body := DoRequest(t, req)
		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		require.Contains(t, string(body), "missing or invalid authentication")
	})
}
