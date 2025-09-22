package e2e

import (
	"context"
	"testing"

	"github.com/malyshevhen/rule-engine/client"
	"github.com/stretchr/testify/require"
)

// TestHealthCheck
// - Verify the tests configured correctly
func TestHealthCheck(t *testing.T) {
	ctx := context.Background()

	// Verify environment is set up correctly
	require.NotNil(t, testEnv)

	// Create client
	baseURL := testEnv.GetRuleEngineURL(ctx, t)
	c := client.NewClient(baseURL, client.AuthConfig{
		APIKey: "test-api-key",
	})

	// Check health
	health, err := c.Health(ctx)
	require.NoError(t, err)

	// Verify database and redis status
	require.Equal(t, "ok", health.Database)
	require.Contains(t, []string{"ok", "disabled"}, health.Redis)
}
