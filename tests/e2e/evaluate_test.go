package e2e

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

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

	t.Run("SimpleScript", func(t *testing.T) {
		// Create HTTP client
		client := &http.Client{Timeout: 10 * time.Second}
		baseURL := env.GetRuleEngineURL(ctx, t)

		// Test simple arithmetic script
		evaluateReq := map[string]string{
			"script": "return 2 + 3",
		}
		reqBody, err := json.Marshal(evaluateReq)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", baseURL+"/api/v1/evaluate", bytes.NewReader(reqBody))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "ApiKey test-api-key")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]any
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		require.True(t, response["success"].(bool))
		require.Equal(t, []any{5.0}, response["output"])
		require.Empty(t, response["error"])
		require.NotEmpty(t, response["duration"].(string))
	})

	t.Run("ScriptWithPlatformAPI", func(t *testing.T) {
		// Create HTTP client
		client := &http.Client{Timeout: 10 * time.Second}
		baseURL := env.GetRuleEngineURL(ctx, t)

		// Test script using platform API (time module)
		evaluateReq := map[string]string{
			"script": complexScript,
		}
		reqBody, err := json.Marshal(evaluateReq)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", baseURL+"/api/v1/evaluate", bytes.NewReader(reqBody))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "ApiKey test-api-key")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]any
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		require.True(t, response["success"].(bool))
		require.Equal(t, []any{true}, response["output"])
		require.Empty(t, response["error"])
		require.NotEmpty(t, response["duration"].(string))
	})

	t.Run("ScriptError", func(t *testing.T) {
		// Create HTTP client
		client := &http.Client{Timeout: 10 * time.Second}
		baseURL := env.GetRuleEngineURL(ctx, t)

		// Test script with syntax error
		evaluateReq := map[string]string{
			"script": "invalid lua syntax {{{",
		}
		reqBody, err := json.Marshal(evaluateReq)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", baseURL+"/api/v1/evaluate", bytes.NewReader(reqBody))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "ApiKey test-api-key")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]any
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		require.False(t, response["success"].(bool))
		require.Empty(t, response["output"])
		require.Contains(t, response["error"].(string), "parse error")
		require.NotEmpty(t, response["duration"].(string))
	})

	t.Run("EmptyScript", func(t *testing.T) {
		// Create HTTP client
		client := &http.Client{Timeout: 10 * time.Second}
		baseURL := env.GetRuleEngineURL(ctx, t)

		// Test empty script
		evaluateReq := map[string]string{
			"script": "",
		}
		reqBody, err := json.Marshal(evaluateReq)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", baseURL+"/api/v1/evaluate", bytes.NewReader(reqBody))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "ApiKey test-api-key")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var errorResp map[string]string
		err = json.Unmarshal(body, &errorResp)
		require.NoError(t, err)
		require.Contains(t, errorResp["error"], "Script cannot be empty")
	})

	t.Run("ScriptTooLong", func(t *testing.T) {
		// Create HTTP client
		client := &http.Client{Timeout: 10 * time.Second}
		baseURL := env.GetRuleEngineURL(ctx, t)

		// Test script that's too long (over 10,000 characters)
		longScript := make([]byte, 10001)
		for i := range longScript {
			longScript[i] = 'a'
		}

		evaluateReq := map[string]string{
			"script": string(longScript),
		}
		reqBody, err := json.Marshal(evaluateReq)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", baseURL+"/api/v1/evaluate", bytes.NewReader(reqBody))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "ApiKey test-api-key")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var errorResp map[string]string
		err = json.Unmarshal(body, &errorResp)
		require.NoError(t, err)
		require.Contains(t, errorResp["error"], "Script too long")
	})

	t.Run("Unauthorized", func(t *testing.T) {
		// Create HTTP client
		client := &http.Client{Timeout: 10 * time.Second}
		baseURL := env.GetRuleEngineURL(ctx, t)

		// Test without authorization header
		evaluateReq := map[string]string{
			"script": "return 42",
		}
		reqBody, err := json.Marshal(evaluateReq)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", baseURL+"/api/v1/evaluate", bytes.NewReader(reqBody))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		// Missing Authorization header

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}
