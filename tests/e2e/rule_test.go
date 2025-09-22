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
		ruleID := rule.ID

		// Update the rule using JSON Patch
		patches := client.PatchRequest{
			{Op: "replace", Path: "/name", Value: "Updated Test Rule"},
			{Op: "replace", Path: "/lua_script", Value: "if event.temperature > 30 then return true end"},
		}
		updatedRule, err := c.UpdateRule(ctx, ruleID, patches)
		require.NoError(t, err)
		require.Equal(t, "Updated Test Rule", updatedRule.Name)
		require.Equal(t, "if event.temperature > 30 then return true end", updatedRule.LuaScript)
		require.Equal(t, true, updatedRule.Enabled)
		require.Equal(t, 0, updatedRule.Priority)
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
		// Create a valid action
		actionReq := client.CreateActionRequest{
			LuaScript: "log_message('info', 'test action')",
			Enabled:   &enabled,
		}
		action, err := c.CreateAction(ctx, actionReq)
		require.NoError(t, err)
		actionID := action.ID

		// Try to add action to a non-existent rule
		nonExistentRuleID := uuid.New()
		req := client.AddActionToRuleRequest{
			ActionID: actionID,
		}
		err = c.AddActionToRule(ctx, nonExistentRuleID, req)
		require.Error(t, err)
		// Check that it's a 404 error
		if apiErr, ok := err.(*client.APIError); ok {
			require.Equal(t, 404, apiErr.StatusCode)
		}
	})

	t.Run("AddNonExistentActionToRule", func(t *testing.T) {
		// Create a valid rule
		ruleReq := client.CreateRuleRequest{
			Name:      "Test Rule for Error",
			LuaScript: "if event.temperature > 25 then return true end",
			Priority:  &priority,
			Enabled:   &enabled,
		}
		rule, err := c.CreateRule(ctx, ruleReq)
		require.NoError(t, err)
		ruleID := rule.ID

		// Try to add a non-existent action to the rule
		nonExistentActionID := uuid.New()
		req := client.AddActionToRuleRequest{
			ActionID: nonExistentActionID,
		}
		err = c.AddActionToRule(ctx, ruleID, req)
		require.Error(t, err)
		// This might return 400 or 404 depending on implementation
		if apiErr, ok := err.(*client.APIError); ok {
			require.True(t, apiErr.StatusCode >= 400)
		}
	})

	t.Run("AddActionWithInvalidRuleID", func(t *testing.T) {
		// TODO: Implement test for invalid UUID format in rule ID
		// This should test that providing an invalid UUID returns 400 Bad Request
		// Requires either modifying client to bypass UUID validation or using direct HTTP calls
		// Example: make HTTP POST to /api/v1/rules/invalid-uuid/actions with valid body
		// Expect 400 status
		t.Skip("Invalid UUID format testing requires direct HTTP calls")
	})

	t.Run("AddActionWithInvalidRequestBody", func(t *testing.T) {
		// Create a valid rule
		ruleReq := client.CreateRuleRequest{
			Name:      "Test Rule for Error",
			LuaScript: "if event.temperature > 25 then return true end",
			Priority:  &priority,
			Enabled:   &enabled,
		}
		rule, err := c.CreateRule(ctx, ruleReq)
		require.NoError(t, err)
		ruleID := rule.ID

		// Try to add action with invalid request body (empty action ID)
		req := client.AddActionToRuleRequest{
			// ActionID is zero UUID, which should be invalid
		}
		err = c.AddActionToRule(ctx, ruleID, req)
		require.Error(t, err)
		if apiErr, ok := err.(*client.APIError); ok {
			require.Equal(t, 400, apiErr.StatusCode)
		}
	})
}
