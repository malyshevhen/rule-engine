package e2e

import (
	"context"
	"encoding/json"
	"net/http"
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

	baseURL := env.GetRuleEngineURL(ctx, t) + "/api/v1"

	var createdRuleID string

	t.Run("CreateRule", func(t *testing.T) {
		reqBody := `{"name": "Test Rule", "lua_script": "if event.temperature > 25 then return true end", "enabled": true, "priority": 0}`
		req, err := MakeAuthenticatedRequest("POST", baseURL+"/rules", reqBody)
		require.NoError(t, err)

		resp, body := DoRequest(t, req)
		require.Equal(t, http.StatusCreated, resp.StatusCode)

		var rule map[string]any
		err = json.Unmarshal(body, &rule)
		require.NoError(t, err)
		require.NotEmpty(t, rule["id"])
		require.Equal(t, "Test Rule", rule["name"])
		require.Equal(t, "if event.temperature > 25 then return true end", rule["lua_script"])
		require.Equal(t, true, rule["enabled"])
		require.Equal(t, float64(0), rule["priority"])
		require.NotEmpty(t, rule["created_at"])
		require.NotEmpty(t, rule["updated_at"])
		// require.NotNil(t, rule["actions"])
		// require.NotNil(t, rule["triggers"])

		createdRuleID = rule["id"].(string)
	})

	t.Run("GetRule", func(t *testing.T) {
		require.NotEmpty(t, createdRuleID)
		req, err := MakeAuthenticatedRequest("GET", baseURL+"/rules/"+createdRuleID, "")
		require.NoError(t, err)

		resp, body := DoRequest(t, req)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var rule map[string]any
		err = json.Unmarshal(body, &rule)
		require.NoError(t, err)
		require.Equal(t, createdRuleID, rule["id"])
		require.Equal(t, "Test Rule", rule["name"])
		require.Equal(t, "if event.temperature > 25 then return true end", rule["lua_script"])
		require.Equal(t, true, rule["enabled"])
		require.Equal(t, float64(0), rule["priority"])
	})

	t.Run("GetRules", func(t *testing.T) {
		req, err := MakeAuthenticatedRequest("GET", baseURL+"/rules", "")
		require.NoError(t, err)

		resp, body := DoRequest(t, req)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var rules []map[string]any
		err = json.Unmarshal(body, &rules)
		require.NoError(t, err)
		require.Greater(t, len(rules), 0)

		// Check that our created rule is in the list
		found := false
		for _, rule := range rules {
			if rule["id"] == createdRuleID {
				found = true
				require.Equal(t, "Test Rule", rule["name"])
				break
			}
		}
		require.True(t, found, "Created rule not found in list")
	})

	t.Run("UpdateRule", func(t *testing.T) {
		require.NotEmpty(t, createdRuleID)
		reqBody := `{"name": "Updated Test Rule", "lua_script": "if event.temperature > 30 then return true end", "enabled": false, "priority": 5}`
		req, err := MakeAuthenticatedRequest("PUT", baseURL+"/rules/"+createdRuleID, reqBody)
		require.NoError(t, err)

		resp, body := DoRequest(t, req)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var rule map[string]any
		err = json.Unmarshal(body, &rule)
		require.NoError(t, err)
		require.Equal(t, createdRuleID, rule["id"])
		require.Equal(t, "Updated Test Rule", rule["name"])
		require.Equal(t, "if event.temperature > 30 then return true end", rule["lua_script"])
		require.Equal(t, false, rule["enabled"])
		require.Equal(t, float64(5), rule["priority"])
	})

	t.Run("DeleteRule", func(t *testing.T) {
		require.NotEmpty(t, createdRuleID)
		req, err := MakeAuthenticatedRequest("DELETE", baseURL+"/rules/"+createdRuleID, "")
		require.NoError(t, err)

		resp, body := DoRequest(t, req)
		require.Equal(t, http.StatusNoContent, resp.StatusCode)
		require.Empty(t, body)

		// Verify it's deleted by trying to get it
		req, err = MakeAuthenticatedRequest("GET", baseURL+"/rules/"+createdRuleID, "")
		require.NoError(t, err)
		resp, _ = DoRequest(t, req)
		require.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}
