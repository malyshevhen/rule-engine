package e2e

import (
	"context"
	"encoding/json"
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
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Parse JSON response
	var healthResponse map[string]string
	err = json.Unmarshal(body, &healthResponse)
	require.NoError(t, err)

	// Verify database and redis status
	require.Equal(t, "ok", healthResponse["database"])
	require.Contains(t, []string{"ok", "disabled"}, healthResponse["redis"])
}
