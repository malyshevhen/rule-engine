package e2e

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRule(t *testing.T) {
	ctx := context.Background()

	// Setup test environment
	env, cleanup := SetupTestEnvironment(ctx, t)
	defer cleanup()

	// Verify environment is set up correctly
	require.NotNil(t, env)

	// Create client
	baseURL := env.GetRuleEngineURL(ctx, t)
	client := NewTestClient(baseURL)

	var createdRuleID string

	t.Run("CreateRule", func(t *testing.T) {
		priority := 0
		enabled := true
		rule := client.CreateRule(ctx, t, "Test Rule", "if event.temperature > 25 then return true end", &priority, &enabled)
		require.NotEmpty(t, rule.ID)
		require.Equal(t, "Test Rule", rule.Name)
		require.Equal(t, "if event.temperature > 25 then return true end", rule.LuaScript)
		require.Equal(t, true, rule.Enabled)
		require.Equal(t, 0, rule.Priority)
		require.NotEmpty(t, rule.CreatedAt)
		require.NotEmpty(t, rule.UpdatedAt)

		createdRuleID = rule.ID
	})

	t.Run("GetRule", func(t *testing.T) {
		require.NotEmpty(t, createdRuleID)
		rule := client.GetRule(ctx, t, createdRuleID)
		require.Equal(t, createdRuleID, rule.ID)
		require.Equal(t, "Test Rule", rule.Name)
		require.Equal(t, "if event.temperature > 25 then return true end", rule.LuaScript)
		require.Equal(t, true, rule.Enabled)
		require.Equal(t, 0, rule.Priority)
	})

	t.Run("GetRules", func(t *testing.T) {
		rules := client.ListRules(ctx, t)
		require.Greater(t, rules.Count, 0)

		// Check that our created rule is in the list
		found := false
		for _, rule := range rules.Rules {
			if rule.ID == createdRuleID {
				found = true
				require.Equal(t, "Test Rule", rule.Name)
				break
			}
		}
		require.True(t, found, "Created rule not found in list")
	})

	t.Run("UpdateRule", func(t *testing.T) {
		// If createdRuleID is empty (when running this test individually), create a rule first
		ruleID := createdRuleID
		if ruleID == "" {
			priority := 0
			enabled := true
			rule := client.CreateRule(ctx, t, "Test Rule for Update", "if event.temperature > 25 then return true end", &priority, &enabled)
			ruleID = rule.ID
		}
		require.NotEmpty(t, ruleID)

		// For now, skip the update test since we don't have JSON Patch support in the client wrapper
		t.Skip("Update test requires JSON Patch support in client wrapper")
	})

	t.Run("DeleteRule", func(t *testing.T) {
		// If createdRuleID is empty (when running this test individually), create a rule first
		ruleID := createdRuleID
		if ruleID == "" {
			priority := 0
			enabled := true
			rule := client.CreateRule(ctx, t, "Test Rule for Delete", "if event.temperature > 25 then return true end", &priority, &enabled)
			ruleID = rule.ID
		}
		require.NotEmpty(t, ruleID)

		client.DeleteRule(ctx, t, ruleID)

		// Verify it's deleted by trying to get it - this should fail
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected GetRule to fail for deleted rule")
			}
		}()
		client.GetRule(ctx, t, ruleID)
	})
}

func TestAddActionToRule(t *testing.T) {
	ctx := context.Background()

	// Setup test environment
	env, cleanup := SetupTestEnvironment(ctx, t)
	defer cleanup()

	// Verify environment is set up correctly
	require.NotNil(t, env)

	// Create client
	baseURL := env.GetRuleEngineURL(ctx, t)
	client := NewTestClient(baseURL)

	priority := 0
	enabled := true
	rule := client.CreateRule(ctx, t, "Test Rule for Action", "if event.temperature > 25 then return true end", &priority, &enabled)
	createdRuleID := rule.ID

	action := client.CreateAction(ctx, t, "log_message('info', 'action added to rule')", "", &enabled)
	createdActionID := action.ID

	t.Run("AddActionToRule", func(t *testing.T) {
		require.NotEmpty(t, createdRuleID)
		require.NotEmpty(t, createdActionID)

		client.AddActionToRule(ctx, t, createdRuleID, createdActionID)
	})

	t.Run("AddActionToNonExistentRule", func(t *testing.T) {
		t.Skip("Skipping error handling tests - need to implement error handling in client wrapper")
	})

	t.Run("AddNonExistentActionToRule", func(t *testing.T) {
		t.Skip("Skipping error handling tests - need to implement error handling in client wrapper")
	})

	t.Run("AddActionWithInvalidRuleID", func(t *testing.T) {
		t.Skip("Skipping error handling tests - need to implement error handling in client wrapper")
	})

	t.Run("AddActionWithInvalidRequestBody", func(t *testing.T) {
		t.Skip("Skipping error handling tests - need to implement error handling in client wrapper")
	})
}
