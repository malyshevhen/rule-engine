package e2e

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTrigger(t *testing.T) {
	ctx := context.Background()

	// Setup test environment
	env, cleanup := SetupTestEnvironment(ctx, t)
	defer cleanup()

	// Verify environment is set up correctly
	require.NotNil(t, env)

	// Create client
	baseURL := env.GetRuleEngineURL(ctx, t)
	client := NewTestClient(baseURL)

	// First create a rule for the trigger
	priority := 0
	enabled := true
	rule := client.CreateRule(ctx, t, "Test Rule for Trigger", "if event.temperature > 25 then return true end", &priority, &enabled)
	ruleID := rule.ID

	var createdTriggerID string

	t.Run("CreateTrigger", func(t *testing.T) {
		enabled := true
		trigger := client.CreateTrigger(ctx, t, ruleID, "CONDITIONAL", "if event.device_id == 'sensor_1' then return true end", &enabled)
		require.NotEmpty(t, trigger.ID)
		require.Equal(t, ruleID, trigger.RuleID)
		require.Equal(t, "if event.device_id == 'sensor_1' then return true end", trigger.ConditionScript)
		require.Equal(t, "CONDITIONAL", trigger.Type)
		require.Equal(t, true, trigger.Enabled)
		require.NotEmpty(t, trigger.CreatedAt)
		require.NotEmpty(t, trigger.UpdatedAt)

		createdTriggerID = trigger.ID
	})

	t.Run("GetTrigger", func(t *testing.T) {
		require.NotEmpty(t, createdTriggerID)
		trigger, err := client.GetTrigger(ctx, t, createdTriggerID)
		require.NoError(t, err)
		require.Equal(t, createdTriggerID, trigger.ID)
		require.Equal(t, ruleID, trigger.RuleID)
		require.Equal(t, "if event.device_id == 'sensor_1' then return true end", trigger.ConditionScript)
		require.Equal(t, "CONDITIONAL", trigger.Type)
		require.Equal(t, true, trigger.Enabled)
	})

	t.Run("GetTriggers", func(t *testing.T) {
		triggers := client.ListTriggers(ctx, t)
		require.Greater(t, len(triggers.Triggers), 0)

		// Check that our created trigger is in the list
		found := false
		for _, trigger := range triggers.Triggers {
			if trigger.ID == createdTriggerID {
				found = true
				require.Equal(t, ruleID, trigger.RuleID)
				break
			}
		}
		require.True(t, found, "Created trigger not found in list")
	})

	t.Run("UpdateTrigger", func(t *testing.T) {
		t.Skip("Update not supported for triggers")
	})

	t.Run("DeleteTrigger", func(t *testing.T) {
		t.Skip("Delete not supported for triggers")
	})
}
