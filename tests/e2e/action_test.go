package e2e

import (
	"context"
	"testing"

	"github.com/malyshevhen/rule-engine/client"
	"github.com/stretchr/testify/require"
)

func TestAction(t *testing.T) {
	ctx := context.Background()

	// Verify environment is set up correctly
	require.NotNil(t, testEnv)

	// Create client
	baseURL := testEnv.GetRuleEngineURL(ctx, t)
	c := client.NewClient(baseURL, client.AuthConfig{
		APIKey: "test-api-key",
	})

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
	})

	t.Run("GetAction", func(t *testing.T) {
		// Create an action for this test
		enabled := true
		req := client.CreateActionRequest{
			LuaScript: "log_message('info', 'get action test')",
			Enabled:   &enabled,
		}
		action, err := c.CreateAction(ctx, req)
		require.NoError(t, err)
		createdActionID := action.ID

		retrievedAction, err := c.GetAction(ctx, createdActionID)
		require.NoError(t, err)
		require.Equal(t, createdActionID, retrievedAction.ID)
		require.Equal(t, "log_message('info', 'get action test')", retrievedAction.LuaScript)
		require.Equal(t, true, retrievedAction.Enabled)
	})

	t.Run("GetActions", func(t *testing.T) {
		// Create an action for this test
		enabled := true
		req := client.CreateActionRequest{
			LuaScript: "log_message('info', 'list actions test')",
			Enabled:   &enabled,
		}
		action, err := c.CreateAction(ctx, req)
		require.NoError(t, err)
		createdActionID := action.ID

		actions, err := c.ListActions(ctx, 100, 0) // limit=100, offset=0
		require.NoError(t, err)
		require.Greater(t, len(actions.Actions), 0)

		// Check that our created action is in the list
		found := false
		for _, action := range actions.Actions {
			if action.ID == createdActionID {
				found = true
				require.Equal(t, "log_message('info', 'list actions test')", action.LuaScript)
				break
			}
		}
		require.True(t, found, "Created action not found in list")
	})

	t.Run("UpdateAction", func(t *testing.T) {
		t.Skip("Update not supported for actions")
	})

	t.Run("DeleteAction", func(t *testing.T) {
		// Create an action for this test
		enabled := true
		req := client.CreateActionRequest{
			LuaScript: "log_message('info', 'delete action test')",
			Enabled:   &enabled,
		}
		action, err := c.CreateAction(ctx, req)
		require.NoError(t, err)
		actionID := action.ID

		// Delete the action
		err = c.DeleteAction(ctx, actionID)
		require.NoError(t, err)

		// Verify it's deleted by trying to get it - this should return an error
		_, err = c.GetAction(ctx, actionID)
		require.Error(t, err, "Expected GetAction to fail for deleted action")
	})
}
