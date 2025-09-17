package e2e

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/require"
)

// Hoverfly simulations
var (
	//go:embed simulations/webhook_success.json
	webhook_success string
)

// Lua scripts
var (
	//go:embed fixtures/actions/post.lua
	post_action string
)

func TestRuleWorkflow_HappyPath(t *testing.T) {
	t.Skip("Skipping until we can get the test to work")

	ctx := context.Background()

	// Setup test environment
	env, cleanup := SetupTestEnvironment(ctx, t)
	defer cleanup()

	// Setup Hoverfly simulation for external API calls
	env.SetupHoverflySimulation(ctx, t, webhook_success)

	// Verify environment is set up correctly
	require.NotNil(t, env)

	// Create HTTP client
	client := &http.Client{Timeout: 10 * time.Second}
	baseURL := env.GetRuleEngineURL(ctx, t)

	// Create action
	actionReq := map[string]any{
		"luaScript": post_action,
		"enabled":   true,
	}
	actionBody, err := json.Marshal(actionReq)
	require.NoError(t, err)

	t.Logf("Creating action: %s", string(actionBody))

	req, err := http.NewRequest("POST", baseURL+"/api/v1/actions", bytes.NewReader(actionBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "ApiKey test-api-key")

	resp, err := client.Do(req)
	require.NoError(t, err)
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Logf("Action creation response: %s", string(body))
		resp.Body.Close()
		t.FailNow()
	}

	var actionResp map[string]any
	json.NewDecoder(resp.Body).Decode(&actionResp)
	resp.Body.Close()

	// Create rule
	ruleReq := map[string]any{
		"name":      "Test Rule",
		"luaScript": "return true",
		"enabled":   true,
		"priority":  0,
	}
	ruleBody, err := json.Marshal(ruleReq)
	require.NoError(t, err)

	req, err = http.NewRequest("POST", baseURL+"/api/v1/rules", bytes.NewReader(ruleBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "ApiKey test-api-key")

	resp, err = client.Do(req)
	require.NoError(t, err)
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Logf("Rule creation response: %s", string(body))
		resp.Body.Close()
		t.FailNow()
	}

	var ruleResp map[string]any
	json.NewDecoder(resp.Body).Decode(&ruleResp)
	resp.Body.Close()

	ruleID := ruleResp["id"].(string)

	// Create trigger
	triggerReq := map[string]any{
		"ruleId":          ruleID,
		"type":            "CONDITIONAL",
		"conditionScript": "return true",
		"enabled":         true,
	}

	triggerBody, err := json.Marshal(triggerReq)
	require.NoError(t, err)

	req, err = http.NewRequest("POST", baseURL+"/api/v1/triggers", bytes.NewReader(triggerBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "ApiKey test-api-key")

	resp, err = client.Do(req)
	require.NoError(t, err)
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Logf("Trigger creation response: %s", string(body))
		resp.Body.Close()
		t.FailNow()
	}
	resp.Body.Close()

	// Connect to NATS and publish event
	nc, err := nats.Connect(env.GetNATSURL(ctx, t))
	require.NoError(t, err)
	defer nc.Close()

	eventData := map[string]any{
		"event": "test",
		"data":  "value",
	}
	eventBody, err := json.Marshal(eventData)
	require.NoError(t, err)

	err = nc.Publish("events.test", eventBody)
	require.NoError(t, err)

	// Wait for processing
	time.Sleep(2 * time.Second)

	// Check analytics for executions
	req, err = http.NewRequest("GET", baseURL+"/api/v1/analytics/dashboard?timeRange=1h", nil)
	require.NoError(t, err)

	req.Header.Set("Authorization", "ApiKey test-api-key")
	resp, err = client.Do(req)
	require.NoError(t, err)
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Logf("Analytics response: %s", string(body))
		resp.Body.Close()
		t.FailNow()
	}

	var analyticsResp map[string]any
	json.NewDecoder(resp.Body).Decode(&analyticsResp)
	resp.Body.Close()

	overallStats := analyticsResp["overall_stats"].(map[string]any)
	totalExecutions := int(overallStats["total_executions"].(float64))
	require.Greater(t, totalExecutions, 0, "Expected at least one rule execution")
}

func TestRuleWorkflow_ComplexConditions(t *testing.T) {
	t.Skip("Complex rule conditions test - implementation pending")

	ctx := context.Background()

	// Setup test environment
	env, cleanup := SetupTestEnvironment(ctx, t)
	defer cleanup()

	// Setup Hoverfly simulation
	env.SetupHoverflySimulation(ctx, t, "simulations/webhook_success.json")

	// Verify environment is set up correctly
	require.NotNil(t, env)

	// TODO: Test complex rule with multiple conditions
	// - Multiple triggers
	// - Complex boolean logic
	// - Conditional actions
}

func TestRuleWorkflow_ActionFailure(t *testing.T) {
	t.Skip("Action failure handling test - implementation pending")

	ctx := context.Background()

	// Setup test environment
	env, cleanup := SetupTestEnvironment(ctx, t)
	defer cleanup()

	// Verify environment is set up correctly
	require.NotNil(t, env)

	// TODO: Setup Hoverfly simulation for failure scenarios
	// Test error handling and retry logic
}

func TestRuleWorkflow_InvalidInput(t *testing.T) {
	t.Skip("Input validation test - implementation pending")

	ctx := context.Background()

	// Setup test environment
	env, cleanup := SetupTestEnvironment(ctx, t)
	defer cleanup()

	// Verify environment is set up correctly
	require.NotNil(t, env)

	// TODO: Test API validation
	// - Invalid rule definitions
	// - Malformed triggers/actions
	// - Missing required fields
}

func TestRuleWorkflow_CRONTrigger(t *testing.T) {
	t.Skip("CRON trigger test - implementation pending")

	ctx := context.Background()

	// Setup test environment
	env, cleanup := SetupTestEnvironment(ctx, t)
	defer cleanup()

	// Verify environment is set up correctly
	require.NotNil(t, env)

	// TODO: Test scheduled rule execution
	// - CRON-based triggers
	// - Time-based rule firing
}
