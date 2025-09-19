package e2e

import (
	"context"
	"encoding/json"
	"net/http"
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

	baseURL := env.GetRuleEngineURL(ctx, t) + "/api/v1"

	// First create a rule for the trigger
	ruleID := createTestRule(t, baseURL)

	var createdTriggerID string

	t.Run("CreateTrigger", func(t *testing.T) {
		reqBody := `{"rule_id": "` + ruleID + `", "condition_script": "if event.device_id == 'sensor_1' then return true end", "type": "CONDITIONAL", "enabled": true}`
		req, err := MakeAuthenticatedRequest("POST", baseURL+"/triggers", reqBody)
		require.NoError(t, err)

		resp, body := DoRequest(t, req)
		require.Equal(t, http.StatusCreated, resp.StatusCode)

		var trigger map[string]any
		err = json.Unmarshal(body, &trigger)
		require.NoError(t, err)
		require.NotEmpty(t, trigger["id"])
		require.Equal(t, ruleID, trigger["rule_id"])
		require.Equal(t, "if event.device_id == 'sensor_1' then return true end", trigger["condition_script"])
		require.Equal(t, "CONDITIONAL", trigger["type"])
		require.Equal(t, true, trigger["enabled"])
		require.NotEmpty(t, trigger["created_at"])
		require.NotEmpty(t, trigger["updated_at"])

		createdTriggerID = trigger["id"].(string)
	})

	t.Run("GetTrigger", func(t *testing.T) {
		require.NotEmpty(t, createdTriggerID)
		req, err := MakeAuthenticatedRequest("GET", baseURL+"/triggers/"+createdTriggerID, "")
		require.NoError(t, err)

		resp, body := DoRequest(t, req)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var trigger map[string]any
		err = json.Unmarshal(body, &trigger)
		require.NoError(t, err)
		require.Equal(t, createdTriggerID, trigger["id"])
		require.Equal(t, ruleID, trigger["rule_id"])
		require.Equal(t, "if event.device_id == 'sensor_1' then return true end", trigger["condition_script"])
		require.Equal(t, "CONDITIONAL", trigger["type"])
		require.Equal(t, true, trigger["enabled"])
	})

	t.Run("GetTriggers", func(t *testing.T) {
		req, err := MakeAuthenticatedRequest("GET", baseURL+"/triggers", "")
		require.NoError(t, err)

		resp, body := DoRequest(t, req)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var triggers []map[string]any
		err = json.Unmarshal(body, &triggers)
		require.NoError(t, err)
		require.Greater(t, len(triggers), 0)

		// Check that our created trigger is in the list
		found := false
		for _, trigger := range triggers {
			if trigger["id"] == createdTriggerID {
				found = true
				require.Equal(t, ruleID, trigger["rule_id"])
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

// Helper function to create a test rule
func createTestRule(t *testing.T, baseURL string) string {
	reqBody := `{"name": "Test Rule for Trigger", "lua_script": "return true", "enabled": true, "priority": 0}`
	req, err := MakeAuthenticatedRequest("POST", baseURL+"/rules", reqBody)
	require.NoError(t, err)

	resp, body := DoRequest(t, req)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var rule map[string]any
	err = json.Unmarshal(body, &rule)
	require.NoError(t, err)
	return rule["id"].(string)
}
