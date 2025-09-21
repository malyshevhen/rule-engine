package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// Simple client wrapper for e2e tests
type TestClient struct {
	baseURL string
}

func NewTestClient(baseURL string) *TestClient {
	return &TestClient{baseURL: baseURL}
}

func (c *TestClient) CreateAction(ctx context.Context, t *testing.T, luaScript, name string, enabled *bool) *ActionResponse {
	reqBody := map[string]any{
		"lua_script": luaScript,
	}
	if name != "" {
		reqBody["name"] = name
	}
	if enabled != nil {
		reqBody["enabled"] = *enabled
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, err := MakeAuthenticatedRequest("POST", c.baseURL+"/api/v1/actions", string(jsonBody))
	require.NoError(t, err)

	resp, body := DoRequest(t, req)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var action ActionResponse
	err = json.Unmarshal(body, &action)
	require.NoError(t, err)
	return &action
}

func (c *TestClient) GetAction(ctx context.Context, t *testing.T, id string) (*ActionResponse, error) {
	req, err := MakeAuthenticatedRequest("GET", c.baseURL+"/api/v1/actions/"+id, "")
	if err != nil {
		return nil, err
	}

	resp, body := DoRequest(t, req)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var action ActionResponse
	err = json.Unmarshal(body, &action)
	if err != nil {
		return nil, err
	}
	return &action, nil
}

func (c *TestClient) ListActions(ctx context.Context, t *testing.T) *ActionsListResponse {
	req, err := MakeAuthenticatedRequest("GET", c.baseURL+"/api/v1/actions", "")
	require.NoError(t, err)

	resp, body := DoRequest(t, req)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var actions ActionsListResponse
	err = json.Unmarshal(body, &actions)
	require.NoError(t, err)
	return &actions
}

func (c *TestClient) EvaluateScript(ctx context.Context, t *testing.T, script string, context map[string]any) *EvaluateResponse {
	reqBody := map[string]any{
		"script": script,
	}
	if context != nil {
		reqBody["context"] = context
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, err := MakeAuthenticatedRequest("POST", c.baseURL+"/api/v1/evaluate", string(jsonBody))
	require.NoError(t, err)

	resp, body := DoRequest(t, req)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var result EvaluateResponse
	err = json.Unmarshal(body, &result)
	require.NoError(t, err)
	return &result
}

func (c *TestClient) EvaluateScriptWithError(ctx context.Context, t *testing.T, script string, expectedStatus int) *ErrorResponse {
	reqBody := map[string]any{
		"script": script,
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, err := MakeAuthenticatedRequest("POST", c.baseURL+"/api/v1/evaluate", string(jsonBody))
	require.NoError(t, err)

	resp, body := DoRequest(t, req)
	require.Equal(t, expectedStatus, resp.StatusCode)

	var result ErrorResponse
	err = json.Unmarshal(body, &result)
	require.NoError(t, err)
	return &result
}

// Rule methods
func (c *TestClient) CreateRule(ctx context.Context, t *testing.T, name, luaScript string, priority *int, enabled *bool) *RuleResponse {
	reqBody := map[string]any{
		"name":       name,
		"lua_script": luaScript,
	}
	if priority != nil {
		reqBody["priority"] = *priority
	}
	if enabled != nil {
		reqBody["enabled"] = *enabled
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, err := MakeAuthenticatedRequest("POST", c.baseURL+"/api/v1/rules", string(jsonBody))
	require.NoError(t, err)

	resp, body := DoRequest(t, req)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var rule RuleResponse
	err = json.Unmarshal(body, &rule)
	require.NoError(t, err)
	return &rule
}

func (c *TestClient) GetRule(ctx context.Context, t *testing.T, id string) (*RuleResponse, error) {
	req, err := MakeAuthenticatedRequest("GET", c.baseURL+"/api/v1/rules/"+id, "")
	if err != nil {
		return nil, err
	}

	resp, body := DoRequest(t, req)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var rule RuleResponse
	err = json.Unmarshal(body, &rule)
	if err != nil {
		return nil, err
	}
	return &rule, nil
}

func (c *TestClient) ListRules(ctx context.Context, t *testing.T) *RulesListResponse {
	req, err := MakeAuthenticatedRequest("GET", c.baseURL+"/api/v1/rules", "")
	require.NoError(t, err)

	resp, body := DoRequest(t, req)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var rules RulesListResponse
	err = json.Unmarshal(body, &rules)
	require.NoError(t, err)
	return &rules
}

func (c *TestClient) DeleteRule(ctx context.Context, t *testing.T, id string) {
	req, err := MakeAuthenticatedRequest("DELETE", c.baseURL+"/api/v1/rules/"+id, "")
	require.NoError(t, err)

	resp, body := DoRequest(t, req)
	require.Equal(t, http.StatusNoContent, resp.StatusCode)
	require.Empty(t, body)
}

// Trigger methods
func (c *TestClient) CreateTrigger(ctx context.Context, t *testing.T, ruleID, triggerType, conditionScript string, enabled *bool) *TriggerResponse {
	reqBody := map[string]any{
		"rule_id":          ruleID,
		"type":             triggerType,
		"condition_script": conditionScript,
	}
	if enabled != nil {
		reqBody["enabled"] = *enabled
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, err := MakeAuthenticatedRequest("POST", c.baseURL+"/api/v1/triggers", string(jsonBody))
	require.NoError(t, err)

	resp, body := DoRequest(t, req)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var trigger TriggerResponse
	err = json.Unmarshal(body, &trigger)
	require.NoError(t, err)
	return &trigger
}

func (c *TestClient) GetTrigger(ctx context.Context, t *testing.T, id string) (*TriggerResponse, error) {
	req, err := MakeAuthenticatedRequest("GET", c.baseURL+"/api/v1/triggers/"+id, "")
	if err != nil {
		return nil, err
	}

	resp, body := DoRequest(t, req)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var trigger TriggerResponse
	err = json.Unmarshal(body, &trigger)
	if err != nil {
		return nil, err
	}
	return &trigger, nil
}

func (c *TestClient) ListTriggers(ctx context.Context, t *testing.T) *TriggersListResponse {
	req, err := MakeAuthenticatedRequest("GET", c.baseURL+"/api/v1/triggers", "")
	require.NoError(t, err)

	resp, body := DoRequest(t, req)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var triggers TriggersListResponse
	err = json.Unmarshal(body, &triggers)
	require.NoError(t, err)
	return &triggers
}

func (c *TestClient) DeleteTrigger(ctx context.Context, t *testing.T, id string) {
	req, err := MakeAuthenticatedRequest("DELETE", c.baseURL+"/api/v1/triggers/"+id, "")
	require.NoError(t, err)

	resp, body := DoRequest(t, req)
	require.Equal(t, http.StatusNoContent, resp.StatusCode)
	require.Empty(t, body)
}

func (c *TestClient) AddActionToRule(ctx context.Context, t *testing.T, ruleID, actionID string) {
	reqBody := map[string]any{
		"action_id": actionID,
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, err := MakeAuthenticatedRequest("POST", c.baseURL+"/api/v1/rules/"+ruleID+"/actions", string(jsonBody))
	require.NoError(t, err)

	resp, body := DoRequest(t, req)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Contains(t, string(body), "success")
}

func (c *TestClient) Health(ctx context.Context, t *testing.T) *HealthResponse {
	req, err := http.NewRequest("GET", c.baseURL+"/health", nil)
	require.NoError(t, err)

	resp, body := DoRequest(t, req)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var health HealthResponse
	err = json.Unmarshal(body, &health)
	require.NoError(t, err)
	return &health
}

// Response types for evaluate endpoint
type EvaluateResponse struct {
	Success  bool   `json:"success"`
	Result   any    `json:"result,omitempty"`
	Output   []any  `json:"output,omitempty"`
	Error    string `json:"error,omitempty"`
	Duration string `json:"duration"`
}

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type HealthResponse struct {
	Database string `json:"database"`
	Redis    string `json:"redis"`
}

// Rule response types
type RuleResponse struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	LuaScript string            `json:"lua_script"`
	Priority  int               `json:"priority"`
	Enabled   bool              `json:"enabled"`
	Triggers  []TriggerResponse `json:"triggers,omitempty"`
	Actions   []ActionResponse  `json:"actions,omitempty"`
	CreatedAt string            `json:"created_at"`
	UpdatedAt string            `json:"updated_at"`
}

type RulesListResponse struct {
	Rules  []RuleResponse `json:"rules"`
	Limit  int            `json:"limit"`
	Offset int            `json:"offset"`
	Count  int            `json:"count"`
	Total  int            `json:"total"`
}

type TriggerResponse struct {
	ID              string `json:"id"`
	RuleID          string `json:"rule_id"`
	Type            string `json:"type"`
	ConditionScript string `json:"condition_script"`
	Enabled         bool   `json:"enabled"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

type TriggersListResponse struct {
	Triggers []TriggerResponse `json:"triggers"`
	Limit    int               `json:"limit"`
	Offset   int               `json:"offset"`
	Count    int               `json:"count"`
	Total    int               `json:"total"`
}

// Response types for e2e tests
type ActionResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	LuaScript string `json:"lua_script"`
	Enabled   bool   `json:"enabled"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type ActionsListResponse struct {
	Actions []ActionResponse `json:"actions"`
	Limit   int              `json:"limit"`
	Offset  int              `json:"offset"`
	Count   int              `json:"count"`
	Total   int              `json:"total"`
}

func TestAction(t *testing.T) {
	ctx := context.Background()

	// Setup test environment
	env, cleanup := SetupTestEnvironment(ctx, t)
	defer cleanup()

	// Verify environment is set up correctly
	require.NotNil(t, env)

	// Create client
	baseURL := env.GetRuleEngineURL(ctx, t)
	client := NewTestClient(baseURL)

	var createdActionID string

	t.Run("CreateAction", func(t *testing.T) {
		enabled := true
		action := client.CreateAction(ctx, t, "log_message('info', 'test action')", "", &enabled)
		require.NotEmpty(t, action.ID)
		require.Equal(t, "log_message('info', 'test action')", action.LuaScript)
		require.Equal(t, true, action.Enabled)
		require.NotEmpty(t, action.CreatedAt)
		require.NotEmpty(t, action.UpdatedAt)

		createdActionID = action.ID
	})

	t.Run("GetAction", func(t *testing.T) {
		require.NotEmpty(t, createdActionID)
		action, err := client.GetAction(ctx, t, createdActionID)
		require.NoError(t, err)
		require.Equal(t, createdActionID, action.ID)
		require.Equal(t, "log_message('info', 'test action')", action.LuaScript)
		require.Equal(t, true, action.Enabled)
	})

	t.Run("GetActions", func(t *testing.T) {
		actions := client.ListActions(ctx, t)
		require.Greater(t, len(actions.Actions), 0)

		// Check that our created action is in the list
		found := false
		for _, action := range actions.Actions {
			if action.ID == createdActionID {
				found = true
				require.Equal(t, "log_message('info', 'test action')", action.LuaScript)
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
