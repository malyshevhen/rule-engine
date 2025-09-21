package e2e

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/malyshevhen/rule-engine/client"
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
	c := client.NewClient(baseURL, client.AuthConfig{
		APIKey: "test-api-key",
	})

	// First create a rule for the trigger
	priority := 0
	enabled := true
	ruleReq := client.CreateRuleRequest{
		Name:      "Test Rule for Trigger",
		LuaScript: "if event.temperature > 25 then return true end",
		Priority:  &priority,
		Enabled:   &enabled,
	}
	rule, err := c.CreateRule(ctx, ruleReq)
	require.NoError(t, err)
	ruleID := rule.ID.String()

	var createdTriggerID string

	t.Run("CreateTrigger", func(t *testing.T) {
		enabled := true
		ruleUUID, err := uuid.Parse(ruleID)
		require.NoError(t, err)

		req := client.CreateTriggerRequest{
			RuleID:          ruleUUID,
			Type:            "CONDITIONAL",
			ConditionScript: "if event.device_id == 'sensor_1' then return true end",
			Enabled:         &enabled,
		}
		trigger, err := c.CreateTrigger(ctx, req)
		require.NoError(t, err)
		require.NotEmpty(t, trigger.ID)
		require.Equal(t, ruleUUID, trigger.RuleID)
		require.Equal(t, "if event.device_id == 'sensor_1' then return true end", trigger.ConditionScript)
		require.Equal(t, "CONDITIONAL", trigger.Type)
		require.Equal(t, true, trigger.Enabled)
		require.NotEmpty(t, trigger.CreatedAt)
		require.NotEmpty(t, trigger.UpdatedAt)

		createdTriggerID = trigger.ID.String()
	})

	t.Run("GetTrigger", func(t *testing.T) {
		require.NotEmpty(t, createdTriggerID)
		triggerID, err := uuid.Parse(createdTriggerID)
		require.NoError(t, err)

		trigger, err := c.GetTrigger(ctx, triggerID)
		require.NoError(t, err)
		require.Equal(t, triggerID, trigger.ID)
		require.Equal(t, ruleID, trigger.RuleID.String())
		require.Equal(t, "if event.device_id == 'sensor_1' then return true end", trigger.ConditionScript)
		require.Equal(t, "CONDITIONAL", trigger.Type)
		require.Equal(t, true, trigger.Enabled)
	})

	t.Run("GetTriggers", func(t *testing.T) {
		triggers, err := c.ListTriggers(ctx, 100, 0) // limit=100, offset=0
		require.NoError(t, err)
		require.Greater(t, len(triggers.Triggers), 0)

		// Check that our created trigger is in the list
		found := false
		for _, trigger := range triggers.Triggers {
			if trigger.ID.String() == createdTriggerID {
				found = true
				require.Equal(t, ruleID, trigger.RuleID.String())
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
