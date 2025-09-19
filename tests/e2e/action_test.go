package e2e

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAction(t *testing.T) {
	ctx := context.Background()

	// Setup test environment
	env, cleanup := SetupTestEnvironment(ctx, t)
	defer cleanup()

	// Verify environment is set up correctly
	require.NotNil(t, env)

	baseURL := env.GetRuleEngineURL(ctx, t) + "/api/v1"

	var createdActionID string

	t.Run("CreateAction", func(t *testing.T) {
		reqBody := `{"lua_script": "log_message('info', 'test action')", "enabled": true}`
		req, err := MakeAuthenticatedRequest("POST", baseURL+"/actions", reqBody)
		require.NoError(t, err)

		resp, body := DoRequest(t, req)
		if http.StatusCreated != resp.StatusCode {
			t.Logf("Response body: %s", string(body))
			t.Fail()
		}

		var action map[string]any
		err = json.Unmarshal(body, &action)
		require.NoError(t, err)
		require.NotEmpty(t, action["id"])
		require.Equal(t, "log_message('info', 'test action')", action["lua_script"])
		require.Equal(t, true, action["enabled"])
		require.NotEmpty(t, action["created_at"])
		require.NotEmpty(t, action["updated_at"])

		createdActionID = action["id"].(string)
	})

	t.Run("GetAction", func(t *testing.T) {
		require.NotEmpty(t, createdActionID)
		req, err := MakeAuthenticatedRequest("GET", baseURL+"/actions/"+createdActionID, "")
		require.NoError(t, err)

		resp, body := DoRequest(t, req)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var action map[string]interface{}
		err = json.Unmarshal(body, &action)
		require.NoError(t, err)
		require.Equal(t, createdActionID, action["id"])
		require.Equal(t, "log_message('info', 'test action')", action["lua_script"])
		require.Equal(t, true, action["enabled"])
	})

	t.Run("GetActions", func(t *testing.T) {
		req, err := MakeAuthenticatedRequest("GET", baseURL+"/actions", "")
		require.NoError(t, err)

		resp, body := DoRequest(t, req)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var actions []map[string]interface{}
		err = json.Unmarshal(body, &actions)
		require.NoError(t, err)
		require.Greater(t, len(actions), 0)

		// Check that our created action is in the list
		found := false
		for _, action := range actions {
			if action["id"] == createdActionID {
				found = true
				require.Equal(t, "log_message('info', 'test action')", action["lua_script"])
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
