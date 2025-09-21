package e2e

import (
	"bytes"
	"context"
	_ "embed"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

//go:embed fixtures/evaluate/complex_script.lua
var complexScript string

func TestEvaluateEndpoint(t *testing.T) {
	ctx := context.Background()

	// Setup test environment
	env, cleanup := SetupTestEnvironment(ctx, t)
	defer cleanup()

	// Verify environment is set up correctly
	require.NotNil(t, env)

	// Create client
	baseURL := env.GetRuleEngineURL(ctx, t)
	client := NewTestClient(baseURL)

	t.Run("SimpleScript", func(t *testing.T) {
		result := client.EvaluateScript(ctx, t, "return 2 + 3", nil)

		require.True(t, result.Success)
		require.Equal(t, []any{5.0}, result.Output)
		require.Empty(t, result.Error)
		require.NotEmpty(t, result.Duration)
	})

	t.Run("ScriptWithPlatformAPI", func(t *testing.T) {
		result := client.EvaluateScript(ctx, t, complexScript, nil)

		require.True(t, result.Success)
		require.Equal(t, []any{true}, result.Output)
		require.Empty(t, result.Error)
		require.NotEmpty(t, result.Duration)
	})

	t.Run("ScriptError", func(t *testing.T) {
		result := client.EvaluateScript(ctx, t, "invalid lua syntax {{{", nil)

		require.False(t, result.Success)
		require.Empty(t, result.Output)
		require.Contains(t, result.Error, "parse error")
		require.NotEmpty(t, result.Duration)
	})

	t.Run("EmptyScript", func(t *testing.T) {
		errorResp := client.EvaluateScriptWithError(ctx, t, "", 400)
		require.Contains(t, errorResp.Error.Message, "validation failed: Script is required")
	})

	t.Run("ScriptTooLong", func(t *testing.T) {
		longScript := make([]byte, 10001)
		for i := range longScript {
			longScript[i] = 'a'
		}

		errorResp := client.EvaluateScriptWithError(ctx, t, string(longScript), 400)
		require.Contains(t, errorResp.Error.Message, "validation failed: Lua script must be between 1 and 10000 characters")
	})

	t.Run("Unauthorized", func(t *testing.T) {
		// Test without authorization - need to make direct HTTP call for this
		req, err := http.NewRequest("POST", baseURL+"/api/v1/evaluate", bytes.NewReader([]byte(`{"script":"return 42"}`)))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		// Missing Authorization header

		resp, body := DoRequest(t, req)
		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		require.Contains(t, string(body), "Missing authorization header")
	})
}
