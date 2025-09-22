package e2e

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/malyshevhen/rule-engine/client"
	"github.com/stretchr/testify/require"
)

func TestRule(t *testing.T) {
	ctx := context.Background()

	// Verify environment is set up correctly
	require.NotNil(t, testEnv)

	// Create client
	baseURL := testEnv.GetRuleEngineURL(ctx, t)
	c := client.NewClient(baseURL, client.AuthConfig{
		APIKey: "test-api-key",
	})

	t.Run("CreateRule", func(t *testing.T) {
		priority := 0
		enabled := true
		req := client.CreateRuleRequest{
			Name:      "Test Rule",
			LuaScript: "if event.temperature > 25 then return true end",
			Priority:  &priority,
			Enabled:   &enabled,
		}

		rule, err := c.CreateRule(ctx, req)
		require.NoError(t, err)
		require.NotEmpty(t, rule.ID)
		require.Equal(t, "Test Rule", rule.Name)
		require.Equal(t, "if event.temperature > 25 then return true end", rule.LuaScript)
		require.Equal(t, true, rule.Enabled)
		require.Equal(t, 0, rule.Priority)
		require.NotEmpty(t, rule.CreatedAt)
		require.NotEmpty(t, rule.UpdatedAt)
	})

	t.Run("GetRule", func(t *testing.T) {
		// Create a rule for this test
		priority := 0
		enabled := true
		req := client.CreateRuleRequest{
			Name:      "Test Rule for Get",
			LuaScript: "if event.temperature > 25 then return true end",
			Priority:  &priority,
			Enabled:   &enabled,
		}
		rule, err := c.CreateRule(ctx, req)
		require.NoError(t, err)
		ruleID := rule.ID

		retrievedRule, err := c.GetRule(ctx, ruleID)
		require.NoError(t, err)
		require.Equal(t, ruleID, retrievedRule.ID)
		require.Equal(t, "Test Rule for Get", retrievedRule.Name)
		require.Equal(t, "if event.temperature > 25 then return true end", retrievedRule.LuaScript)
		require.Equal(t, true, retrievedRule.Enabled)
		require.Equal(t, 0, retrievedRule.Priority)
	})

	t.Run("GetRules", func(t *testing.T) {
		// Create a rule for this test
		priority := 0
		enabled := true
		req := client.CreateRuleRequest{
			Name:      "Test Rule for List",
			LuaScript: "if event.temperature > 25 then return true end",
			Priority:  &priority,
			Enabled:   &enabled,
		}
		rule, err := c.CreateRule(ctx, req)
		require.NoError(t, err)
		createdRuleID := rule.ID.String()

		rules, err := c.ListRules(ctx, 100, 0) // limit=100, offset=0
		require.NoError(t, err)
		require.Greater(t, rules.Count, 0)

		// Check that our created rule is in the list
		found := false
		for _, rule := range rules.Rules {
			if rule.ID.String() == createdRuleID {
				found = true
				require.Equal(t, "Test Rule for List", rule.Name)
				break
			}
		}
		require.True(t, found, "Created rule not found in list")
	})

	t.Run("UpdateRule", func(t *testing.T) {
		// Create a rule for this test
		priority := 0
		enabled := true
		req := client.CreateRuleRequest{
			Name:      "Test Rule for Update",
			LuaScript: "if event.temperature > 25 then return true end",
			Priority:  &priority,
			Enabled:   &enabled,
		}
		rule, err := c.CreateRule(ctx, req)
		require.NoError(t, err)
		ruleID := rule.ID.String()
		require.NotEmpty(t, ruleID)

		// For now, skip the update test since we don't have JSON Patch support in the client
		t.Skip("Update test requires JSON Patch support in client")
	})

	t.Run("DeleteRule", func(t *testing.T) {
		// Create a rule for this test
		priority := 0
		enabled := true
		req := client.CreateRuleRequest{
			Name:      "Test Rule for Delete",
			LuaScript: "if event.temperature > 25 then return true end",
			Priority:  &priority,
			Enabled:   &enabled,
		}
		rule, err := c.CreateRule(ctx, req)
		require.NoError(t, err)
		ruleID := rule.ID.String()
		require.NotEmpty(t, ruleID)

		ruleUUID, err := uuid.Parse(ruleID)
		require.NoError(t, err)
		err = c.DeleteRule(ctx, ruleUUID)
		require.NoError(t, err)

		// Verify it's deleted by trying to get it - this should return an error
		_, err = c.GetRule(ctx, ruleUUID)
		require.Error(t, err, "Expected GetRule to fail for deleted rule")
	})
}

func TestAddActionToRule(t *testing.T) {
	ctx := context.Background()

	// Verify environment is set up correctly
	require.NotNil(t, testEnv)

	// Create client
	baseURL := testEnv.GetRuleEngineURL(ctx, t)
	c := client.NewClient(baseURL, client.AuthConfig{
		APIKey: "test-api-key",
	})

	priority := 0
	enabled := true
	ruleReq := client.CreateRuleRequest{
		Name:      "Test Rule for Action",
		LuaScript: "if event.temperature > 25 then return true end",
		Priority:  &priority,
		Enabled:   &enabled,
	}
	rule, err := c.CreateRule(ctx, ruleReq)
	require.NoError(t, err)
	createdRuleID := rule.ID

	actionReq := client.CreateActionRequest{
		LuaScript: "log_message('info', 'action added to rule')",
		Enabled:   &enabled,
	}
	action, err := c.CreateAction(ctx, actionReq)
	require.NoError(t, err)
	createdActionID := action.ID

	t.Run("AddActionToRule", func(t *testing.T) {
		require.NotEmpty(t, createdRuleID)
		require.NotEmpty(t, createdActionID)

		req := client.AddActionToRuleRequest{
			ActionID: createdActionID,
		}
		err := c.AddActionToRule(ctx, createdRuleID, req)
		require.NoError(t, err)
	})

	t.Run("AddActionToNonExistentRule", func(t *testing.T) {
		t.Skip("Skipping error handling tests - need to implement error handling in client")
	})

	t.Run("AddNonExistentActionToRule", func(t *testing.T) {
		t.Skip("Skipping error handling tests - need to implement error handling in client")
	})

	t.Run("AddActionWithInvalidRuleID", func(t *testing.T) {
		t.Skip("Skipping error handling tests - need to implement error handling in client")
	})

	t.Run("AddActionWithInvalidRequestBody", func(t *testing.T) {
		t.Skip("Skipping error handling tests - need to implement error handling in client")
	})
}
