package e2e

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/malyshevhen/rule-engine/client"
	"github.com/stretchr/testify/require"
)

func TestAction(t *testing.T) {
	ctx := context.Background()

	// Setup test environment
	env, cleanup := SetupTestEnvironment(ctx, t)
	defer cleanup()

	// Verify environment is set up correctly
	require.NotNil(t, env)

	// Create client
	baseURL := env.GetRuleEngineURL(ctx, t)
	c := client.NewClient(baseURL, client.AuthConfig{
		APIKey: "test-api-key",
	})

	var createdActionID uuid.UUID

	t.Run("CreateAction", func(t *testing.T) {
		enabled := true
		req := client.CreateActionRequest{
			LuaScript: "log_message('info', 'test action')",
			Enabled:   &enabled,
		}

		action, err := c.CreateAction(ctx, req)
		require.NoError(t, err)
		require.NotEmpty(t, action.ID)
		require.Equal(t, "log_message('info', 'test action')", action.LuaScript)
		require.Equal(t, true, action.Enabled)
		require.NotEmpty(t, action.CreatedAt)
		require.NotEmpty(t, action.UpdatedAt)

		createdActionID = action.ID
	})

	t.Run("GetAction", func(t *testing.T) {
		require.NotEmpty(t, createdActionID)
		action, err := c.GetAction(ctx, createdActionID)
		require.NoError(t, err)
		require.Equal(t, createdActionID, action.ID)
		require.Equal(t, "log_message('info', 'test action')", action.LuaScript)
		require.Equal(t, true, action.Enabled)
	})

	t.Run("GetActions", func(t *testing.T) {
		actions, err := c.ListActions(ctx, 100, 0) // limit=100, offset=0
		require.NoError(t, err)
		require.Greater(t, len(actions.Actions), 0)

		// Check that our created action is in the list
		found := false
		for _, action := range actions.Actions {
			if action.ID == createdActionID {
				found = true
				require.Equal(t, "log_message('info', 'test action')", action.LuaScript)
				break
			}
		}
		require.True(t, found, "Created action not found in list")
	})

	t.Run("UpdateAction", func(t *testing.T) {
		t.Skip("Update not supported for actions")
	})

	t.Run("DeleteAction", func(t *testing.T) {
		t.Skip("Delete not supported for actions")
	})
}
