package e2e

import (
	"context"
	"testing"

	"github.com/malyshevhen/rule-engine/client"
	"github.com/stretchr/testify/require"
)

func TestTrigger(t *testing.T) {
	ctx := context.Background()

	// Verify environment is set up correctly
	require.NotNil(t, testEnv)

	// Create client
	baseURL := testEnv.GetRuleEngineURL(ctx, t)
	c := client.NewClient(baseURL, client.AuthConfig{
		APIKey: "test-api-key",
	})

	t.Run("CreateTrigger", func(t *testing.T) {
		// Create a rule for this test
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
		ruleUUID := rule.ID

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
	})

	t.Run("GetTrigger", func(t *testing.T) {
		// Create rule and trigger for this test
		priority := 0
		enabled := true
		ruleReq := client.CreateRuleRequest{
			Name:      "Test Rule for Get Trigger",
			LuaScript: "if event.temperature > 25 then return true end",
			Priority:  &priority,
			Enabled:   &enabled,
		}
		rule, err := c.CreateRule(ctx, ruleReq)
		require.NoError(t, err)
		ruleUUID := rule.ID

		req := client.CreateTriggerRequest{
			RuleID:          ruleUUID,
			Type:            "CONDITIONAL",
			ConditionScript: "if event.device_id == 'sensor_1' then return true end",
			Enabled:         &enabled,
		}
		trigger, err := c.CreateTrigger(ctx, req)
		require.NoError(t, err)
		triggerID := trigger.ID

		retrievedTrigger, err := c.GetTrigger(ctx, triggerID)
		require.NoError(t, err)
		require.Equal(t, triggerID, retrievedTrigger.ID)
		require.Equal(t, ruleUUID, retrievedTrigger.RuleID)
		require.Equal(t, "if event.device_id == 'sensor_1' then return true end", retrievedTrigger.ConditionScript)
		require.Equal(t, "CONDITIONAL", retrievedTrigger.Type)
		require.Equal(t, true, retrievedTrigger.Enabled)
	})

	t.Run("GetTriggers", func(t *testing.T) {
		// Create rule and trigger for this test
		priority := 0
		enabled := true
		ruleReq := client.CreateRuleRequest{
			Name:      "Test Rule for List Triggers",
			LuaScript: "if event.temperature > 25 then return true end",
			Priority:  &priority,
			Enabled:   &enabled,
		}
		rule, err := c.CreateRule(ctx, ruleReq)
		require.NoError(t, err)
		ruleUUID := rule.ID

		req := client.CreateTriggerRequest{
			RuleID:          ruleUUID,
			Type:            "CONDITIONAL",
			ConditionScript: "if event.device_id == 'sensor_1' then return true end",
			Enabled:         &enabled,
		}
		trigger, err := c.CreateTrigger(ctx, req)
		require.NoError(t, err)
		createdTriggerID := trigger.ID.String()

		triggers, err := c.ListTriggers(ctx, 100, 0) // limit=100, offset=0
		require.NoError(t, err)
		require.Greater(t, len(triggers.Triggers), 0)

		// Check that our created trigger is in the list
		found := false
		for _, trigger := range triggers.Triggers {
			if trigger.ID.String() == createdTriggerID {
				found = true
				require.Equal(t, ruleUUID.String(), trigger.RuleID.String())
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
