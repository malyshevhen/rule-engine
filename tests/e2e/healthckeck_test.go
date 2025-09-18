package e2e

import (
	"context"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestHealthCheck
// - Verify the tests configured correctly
func TestHealthCheck(t *testing.T) {
	ctx := context.Background()

	// Setup test environment
	env, cleanup := SetupTestEnvironment(ctx, t)
	defer cleanup()

	// Verify environment is set up correctly
	require.NotNil(t, env)

	// Create HTTP client
	client := &http.Client{Timeout: 10 * time.Second}
	baseURL := env.GetRuleEngineURL(ctx, t)

	// Check healthcheck
	req, err := http.NewRequest("GET", baseURL+"/health", nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)

	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, []byte("healthy"), body)
}
