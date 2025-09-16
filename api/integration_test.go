//go:build integration

package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/malyshevhen/rule-engine/internal/analytics"
	"github.com/malyshevhen/rule-engine/internal/core/action"
	"github.com/malyshevhen/rule-engine/internal/core/rule"
	"github.com/malyshevhen/rule-engine/internal/core/trigger"
	"github.com/malyshevhen/rule-engine/internal/engine/executor"
	execCtx "github.com/malyshevhen/rule-engine/internal/engine/executor/context"
	"github.com/malyshevhen/rule-engine/internal/engine/executor/platform"
	actionStorage "github.com/malyshevhen/rule-engine/internal/storage/action"
	"github.com/malyshevhen/rule-engine/internal/storage/db"
	ruleStorage "github.com/malyshevhen/rule-engine/internal/storage/rule"
	triggerStorage "github.com/malyshevhen/rule-engine/internal/storage/trigger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupIntegrationTest(t *testing.T) (*Server, func()) {
	// Skip if not running integration tests
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Set test API key for authentication
	testAPIKey := "test-api-key-integration"
	originalAPIKey := os.Getenv("API_KEY")
	os.Setenv("API_KEY", testAPIKey)

	// Setup test containers
	ctx := context.Background()
	tc, cleanupContainers := SetupTestContainers(ctx, t)

	// Wait for services to be ready
	tc.WaitForServices(ctx, t)

	// Setup database pool
	pool := tc.SetupDatabasePool(ctx, t)

	// Run migrations
	err := db.RunMigrations(pool)
	require.NoError(t, err)

	// Setup Redis client
	redisClient := tc.GetRedisClient(ctx, t)

	// Create repositories
	ruleRepo := ruleStorage.NewRepository(pool)
	triggerRepo := triggerStorage.NewRepository(pool)
	actionRepo := actionStorage.NewRepository(pool)

	// Create services
	ruleSvc := rule.NewService(ruleRepo, triggerRepo, actionRepo, redisClient)
	triggerSvc := trigger.NewService(triggerRepo, redisClient)
	actionSvc := action.NewService(actionRepo)
	analyticsSvc := analytics.NewService()

	// Create server
	config := &ServerConfig{Port: "8080"}
	server := NewServer(config, ruleSvc, triggerSvc, actionSvc, analyticsSvc)

	// Return cleanup function
	cleanup := func() {
		redisClient.Close()
		pool.Close()
		cleanupContainers()
		if originalAPIKey != "" {
			os.Setenv("API_KEY", originalAPIKey)
		} else {
			os.Unsetenv("API_KEY")
		}
	}

	return server, cleanup
}

func TestIntegration_CreateAndGetRule(t *testing.T) {
	server, cleanup := setupIntegrationTest(t)
	defer cleanup()

	// Create a rule
	createReq := CreateRuleRequest{
		Name:      "Integration Test Rule",
		LuaScript: "return event.temperature > 25",
		Priority:  &[]int{5}[0],
		Enabled:   &[]bool{true}[0],
	}

	body, err := json.Marshal(createReq)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/rules", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "ApiKey test-api-key-integration")
	w := httptest.NewRecorder()

	server.Router().ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var createdRule rule.Rule
	err = json.Unmarshal(w.Body.Bytes(), &createdRule)
	require.NoError(t, err)

	// Check that the rule has an ID (was actually created)
	assert.NotEqual(t, uuid.Nil, createdRule.ID)
	assert.Equal(t, "Integration Test Rule", createdRule.Name)
	assert.Equal(t, "return event.temperature > 25", createdRule.LuaScript)

	// Get the rule back
	getReq := httptest.NewRequest(http.MethodGet, "/api/v1/rules/"+createdRule.ID.String(), nil)
	getReq.Header.Set("Authorization", "ApiKey test-api-key-integration")
	getW := httptest.NewRecorder()

	server.Router().ServeHTTP(getW, getReq)
	assert.Equal(t, http.StatusOK, getW.Code)

	var retrievedRule rule.Rule
	err = json.Unmarshal(getW.Body.Bytes(), &retrievedRule)
	require.NoError(t, err)

	assert.Equal(t, createdRule.ID, retrievedRule.ID)
	assert.Equal(t, createdRule.Name, retrievedRule.Name)
	assert.Equal(t, createdRule.LuaScript, retrievedRule.LuaScript)
}

func TestIntegration_ListRules(t *testing.T) {
	server, cleanup := setupIntegrationTest(t)
	defer cleanup()

	// Create multiple rules
	rules := []CreateRuleRequest{
		{Name: "Rule 1", LuaScript: "return true", Priority: &[]int{0}[0], Enabled: &[]bool{true}[0]},
		{Name: "Rule 2", LuaScript: "return false", Priority: &[]int{5}[0], Enabled: &[]bool{false}[0]},
		{Name: "Rule 3", LuaScript: "return event.value > 10", Priority: &[]int{10}[0], Enabled: &[]bool{true}[0]},
	}

	var createdRules []rule.Rule
	for _, r := range rules {
		body, err := json.Marshal(r)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/rules", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "ApiKey test-api-key-integration")
		w := httptest.NewRecorder()

		server.Router().ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Logf("POST failed with status %d, body: %s", w.Code, w.Body.String())
		}
		assert.Equal(t, http.StatusOK, w.Code)

		var createdRule rule.Rule
		err = json.Unmarshal(w.Body.Bytes(), &createdRule)
		if err != nil {
			t.Logf("Failed to unmarshal response: %s, body: %s", err, w.Body.String())
		}
		require.NoError(t, err)
		createdRules = append(createdRules, createdRule)
	}

	// List all rules
	listReq := httptest.NewRequest(http.MethodGet, "/api/v1/rules", nil)
	listReq.Header.Set("Authorization", "ApiKey test-api-key-integration")
	listW := httptest.NewRecorder()

	server.Router().ServeHTTP(listW, listReq)
	assert.Equal(t, http.StatusOK, listW.Code)

	var listedRules []rule.Rule
	err := json.Unmarshal(listW.Body.Bytes(), &listedRules)
	require.NoError(t, err)

	assert.True(t, len(listedRules) >= len(createdRules), "Should have at least the created rules")
	// Verify all created rules are in the list
	for _, created := range createdRules {
		found := false
		for _, listed := range listedRules {
			if listed.ID == created.ID {
				assert.Equal(t, created.Name, listed.Name)
				assert.Equal(t, created.LuaScript, listed.LuaScript)
				assert.Equal(t, created.Enabled, listed.Enabled)
				found = true
				break
			}
		}
		assert.True(t, found, "Created rule not found in list")
	}
}

func TestIntegration_UpdateRule(t *testing.T) {
	server, cleanup := setupIntegrationTest(t)
	defer cleanup()

	// Create a rule
	createReq := CreateRuleRequest{
		Name:      "Original Rule",
		LuaScript: "return true",
		Priority:  &[]int{0}[0],
		Enabled:   &[]bool{true}[0],
	}

	body, err := json.Marshal(createReq)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/rules", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "ApiKey test-api-key-integration")
	w := httptest.NewRecorder()

	server.Router().ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var createdRule rule.Rule
	err = json.Unmarshal(w.Body.Bytes(), &createdRule)
	require.NoError(t, err)

	// Update the rule
	newName := "Updated Rule"
	newScript := "return event.temperature > 30"
	newPriority := 10
	newEnabled := false

	updateReq := UpdateRuleRequest{
		Name:      &newName,
		LuaScript: &newScript,
		Priority:  &newPriority,
		Enabled:   &newEnabled,
	}

	updateBody, err := json.Marshal(updateReq)
	require.NoError(t, err)

	updateReqHTTP := httptest.NewRequest(http.MethodPut, "/api/v1/rules/"+createdRule.ID.String(), bytes.NewReader(updateBody))
	updateReqHTTP.Header.Set("Content-Type", "application/json")
	updateReqHTTP.Header.Set("Authorization", "ApiKey test-api-key-integration")
	updateW := httptest.NewRecorder()

	server.Router().ServeHTTP(updateW, updateReqHTTP)
	assert.Equal(t, http.StatusOK, updateW.Code)

	var updatedRule rule.Rule
	err = json.Unmarshal(updateW.Body.Bytes(), &updatedRule)
	require.NoError(t, err)

	assert.Equal(t, createdRule.ID, updatedRule.ID)
	assert.Equal(t, newName, updatedRule.Name)
	assert.Equal(t, newScript, updatedRule.LuaScript)
	assert.Equal(t, newEnabled, updatedRule.Enabled)
}

func TestIntegration_DeleteRule(t *testing.T) {
	server, cleanup := setupIntegrationTest(t)
	defer cleanup()

	// Create a rule
	createReq := CreateRuleRequest{
		Name:      "Rule to Delete",
		LuaScript: "return false",
		Priority:  &[]int{0}[0],
		Enabled:   &[]bool{true}[0],
	}

	body, err := json.Marshal(createReq)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/rules", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "ApiKey test-api-key-integration")
	w := httptest.NewRecorder()

	server.Router().ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var createdRule rule.Rule
	err = json.Unmarshal(w.Body.Bytes(), &createdRule)
	require.NoError(t, err)

	// Delete the rule
	deleteReq := httptest.NewRequest(http.MethodDelete, "/api/v1/rules/"+createdRule.ID.String(), nil)
	deleteReq.Header.Set("Authorization", "ApiKey test-api-key-integration")
	deleteW := httptest.NewRecorder()

	server.Router().ServeHTTP(deleteW, deleteReq)
	assert.Equal(t, http.StatusNoContent, deleteW.Code)

	// Verify rule is deleted - should return 404
	getReq := httptest.NewRequest(http.MethodGet, "/api/v1/rules/"+createdRule.ID.String(), nil)
	getReq.Header.Set("Authorization", "ApiKey test-api-key-integration")
	getW := httptest.NewRecorder()

	server.Router().ServeHTTP(getW, getReq)
	assert.Equal(t, http.StatusNotFound, getW.Code)
}

func TestIntegration_ListTriggers(t *testing.T) {
	server, cleanup := setupIntegrationTest(t)
	defer cleanup()

	// First create a rule for triggers
	ruleSvc := server.ruleSvc.(*rule.Service)
	testRule := &rule.Rule{
		Name:      "Test Rule for Triggers",
		LuaScript: "return true",
		Enabled:   true,
	}
	err := ruleSvc.Create(context.Background(), testRule)
	require.NoError(t, err)

	// Create multiple triggers
	triggers := []CreateTriggerRequest{
		{RuleID: testRule.ID, Type: "conditional", ConditionScript: "event.type == 'temp'", Enabled: &[]bool{true}[0]},
		{RuleID: testRule.ID, Type: "conditional", ConditionScript: "event.value > 10", Enabled: &[]bool{false}[0]},
		{RuleID: testRule.ID, Type: "scheduled", ConditionScript: "0 */5 * * * *", Enabled: &[]bool{true}[0]},
	}

	var createdTriggers []trigger.Trigger
	for _, tr := range triggers {
		body, err := json.Marshal(tr)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/triggers", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "ApiKey test-api-key-integration")
		w := httptest.NewRecorder()

		server.Router().ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		var createdTrigger trigger.Trigger
		err = json.Unmarshal(w.Body.Bytes(), &createdTrigger)
		require.NoError(t, err)
		createdTriggers = append(createdTriggers, createdTrigger)
	}

	// List all triggers
	listReq := httptest.NewRequest(http.MethodGet, "/api/v1/triggers", nil)
	listReq.Header.Set("Authorization", "ApiKey test-api-key-integration")
	listW := httptest.NewRecorder()

	server.Router().ServeHTTP(listW, listReq)
	assert.Equal(t, http.StatusOK, listW.Code)

	var listedTriggers []trigger.Trigger
	err = json.Unmarshal(listW.Body.Bytes(), &listedTriggers)
	require.NoError(t, err)

	assert.Len(t, listedTriggers, len(createdTriggers))
	// Verify all created triggers are in the list
	for _, created := range createdTriggers {
		found := false
		for _, listed := range listedTriggers {
			if listed.ID == created.ID {
				assert.Equal(t, created.RuleID, listed.RuleID)
				assert.Equal(t, created.Type, listed.Type)
				assert.Equal(t, created.ConditionScript, listed.ConditionScript)
				assert.Equal(t, created.Enabled, listed.Enabled)
				found = true
				break
			}
		}
		assert.True(t, found, "Created trigger not found in list")
	}
}

func TestIntegration_GetTrigger(t *testing.T) {
	server, cleanup := setupIntegrationTest(t)
	defer cleanup()

	// First create a rule
	ruleSvc := server.ruleSvc.(*rule.Service)
	testRule := &rule.Rule{
		Name:      "Test Rule for Get Trigger",
		LuaScript: "return true",
		Enabled:   true,
	}
	err := ruleSvc.Create(context.Background(), testRule)
	require.NoError(t, err)

	// Create a trigger
	createReq := CreateTriggerRequest{
		RuleID:          testRule.ID,
		Type:            "conditional",
		ConditionScript: "event.temperature > 25",
		Enabled:         &[]bool{true}[0],
	}

	body, err := json.Marshal(createReq)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/triggers", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "ApiKey test-api-key-integration")
	w := httptest.NewRecorder()

	server.Router().ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var createdTrigger trigger.Trigger
	err = json.Unmarshal(w.Body.Bytes(), &createdTrigger)
	require.NoError(t, err)

	// Get the trigger back
	getReq := httptest.NewRequest(http.MethodGet, "/api/v1/triggers/"+createdTrigger.ID.String(), nil)
	getReq.Header.Set("Authorization", "ApiKey test-api-key-integration")
	getW := httptest.NewRecorder()

	server.Router().ServeHTTP(getW, getReq)
	assert.Equal(t, http.StatusOK, getW.Code)

	var retrievedTrigger trigger.Trigger
	err = json.Unmarshal(getW.Body.Bytes(), &retrievedTrigger)
	require.NoError(t, err)

	assert.Equal(t, createdTrigger.ID, retrievedTrigger.ID)
	assert.Equal(t, createdTrigger.RuleID, retrievedTrigger.RuleID)
	assert.Equal(t, createdTrigger.Type, retrievedTrigger.Type)
	assert.Equal(t, createdTrigger.ConditionScript, retrievedTrigger.ConditionScript)
	assert.Equal(t, createdTrigger.Enabled, retrievedTrigger.Enabled)
}

func TestIntegration_ListActions(t *testing.T) {
	server, cleanup := setupIntegrationTest(t)
	defer cleanup()

	// Create multiple actions
	actions := []CreateActionRequest{
		{LuaScript: "send_notification('alert 1')", Enabled: &[]bool{true}[0]},
		{LuaScript: "log_message('action 2 executed')", Enabled: &[]bool{false}[0]},
		{LuaScript: "update_device_state('device123', 'active')", Enabled: &[]bool{true}[0]},
	}

	var createdActions []action.Action
	for _, a := range actions {
		body, err := json.Marshal(a)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/actions", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "ApiKey test-api-key-integration")
		w := httptest.NewRecorder()

		server.Router().ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		var createdAction action.Action
		err = json.Unmarshal(w.Body.Bytes(), &createdAction)
		require.NoError(t, err)
		createdActions = append(createdActions, createdAction)
	}

	// List all actions
	listReq := httptest.NewRequest(http.MethodGet, "/api/v1/actions", nil)
	listReq.Header.Set("Authorization", "ApiKey test-api-key-integration")
	listW := httptest.NewRecorder()

	server.Router().ServeHTTP(listW, listReq)
	assert.Equal(t, http.StatusOK, listW.Code)

	var listedActions []action.Action
	err := json.Unmarshal(listW.Body.Bytes(), &listedActions)
	require.NoError(t, err)

	assert.Len(t, listedActions, len(createdActions))
	// Verify all created actions are in the list
	for _, created := range createdActions {
		found := false
		for _, listed := range listedActions {
			if listed.ID == created.ID {
				assert.Equal(t, created.LuaScript, listed.LuaScript)
				assert.Equal(t, created.Enabled, listed.Enabled)
				found = true
				break
			}
		}
		assert.True(t, found, "Created action not found in list")
	}
}

func TestIntegration_GetAction(t *testing.T) {
	server, cleanup := setupIntegrationTest(t)
	defer cleanup()

	// Create an action
	createReq := CreateActionRequest{
		LuaScript: "send_notification('test alert')",
		Enabled:   &[]bool{true}[0],
	}

	body, err := json.Marshal(createReq)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/actions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "ApiKey test-api-key-integration")
	w := httptest.NewRecorder()

	server.Router().ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var createdAction action.Action
	err = json.Unmarshal(w.Body.Bytes(), &createdAction)
	require.NoError(t, err)

	// Get the action back
	getReq := httptest.NewRequest(http.MethodGet, "/api/v1/actions/"+createdAction.ID.String(), nil)
	getReq.Header.Set("Authorization", "ApiKey test-api-key-integration")
	getW := httptest.NewRecorder()

	server.Router().ServeHTTP(getW, getReq)
	assert.Equal(t, http.StatusOK, getW.Code)

	var retrievedAction action.Action
	err = json.Unmarshal(getW.Body.Bytes(), &retrievedAction)
	require.NoError(t, err)

	assert.Equal(t, createdAction.ID, retrievedAction.ID)
	assert.Equal(t, createdAction.LuaScript, retrievedAction.LuaScript)
	assert.Equal(t, createdAction.Enabled, retrievedAction.Enabled)
}

func TestIntegration_RuleExecution(t *testing.T) {
	server, cleanup := setupIntegrationTest(t)
	defer cleanup()

	// Create a rule with a simple Lua script that returns true
	createRuleReq := CreateRuleRequest{
		Name:      "Execution Test Rule",
		LuaScript: "return event.temperature > 20",
		Priority:  &[]int{0}[0],
		Enabled:   &[]bool{true}[0],
	}

	body, err := json.Marshal(createRuleReq)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/rules", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "ApiKey test-api-key-integration")
	w := httptest.NewRecorder()

	server.Router().ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var createdRule rule.Rule
	err = json.Unmarshal(w.Body.Bytes(), &createdRule)
	require.NoError(t, err)

	// Test execution through the executor service
	// Note: In a real scenario, this would be triggered by events, but for testing we call it directly
	contextSvc := execCtx.NewService()
	platformSvc := platform.NewService()
	executorSvc := executor.NewService(contextSvc, platformSvc)

	// Create execution context
	execContext := &execCtx.ExecutionContext{
		RuleID:    createdRule.ID.String(),
		TriggerID: "test-trigger-id",
		Data: map[string]interface{}{
			"temperature": 25.5,
			"humidity":    60,
		},
	}

	// Execute the rule
	result := executorSvc.ExecuteScript(context.Background(), createdRule.LuaScript, execContext)

	// Verify execution result
	assert.True(t, result.Success, "Rule execution should succeed")
	assert.NotEmpty(t, result.Output, "Should have output")
	assert.Equal(t, true, result.Output[0], "Script should return true for temperature > 20")
	assert.Greater(t, result.Duration, time.Duration(0), "Execution should take some time")

	// Test with temperature below threshold
	execContext.Data["temperature"] = 15.0
	result2 := executorSvc.ExecuteScript(context.Background(), createdRule.LuaScript, execContext)

	assert.True(t, result2.Success, "Rule execution should succeed")
	assert.Equal(t, false, result2.Output[0], "Script should return false for temperature <= 20")
}

func TestIntegration_PlatformAPIFunctions(t *testing.T) {
	server, cleanup := setupIntegrationTest(t)
	defer cleanup()

	// Create a rule that uses platform API functions
	createRuleReq := CreateRuleRequest{
		Name: "Platform API Test Rule",
		LuaScript: `
			log_message("Rule executed with temperature: " .. event.temperature)
			local device_state = get_device_state("device123")
			store_data("last_execution", os.time())
			return device_state ~= nil
		`,
		Priority: &[]int{0}[0],
		Enabled:  &[]bool{true}[0],
	}

	body, err := json.Marshal(createRuleReq)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/rules", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "ApiKey test-api-key-integration")
	w := httptest.NewRecorder()

	server.Router().ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var createdRule rule.Rule
	err = json.Unmarshal(w.Body.Bytes(), &createdRule)
	require.NoError(t, err)

	// Test execution with platform API
	contextSvc := execCtx.NewService()
	platformSvc := platform.NewService()
	executorSvc := executor.NewService(contextSvc, platformSvc)

	execContext := &execCtx.ExecutionContext{
		RuleID:    createdRule.ID.String(),
		TriggerID: "test-trigger-id",
		Data: map[string]interface{}{
			"temperature": 22.0,
		},
	}

	// Execute the rule - this should not fail even though platform functions may not work in test
	result := executorSvc.ExecuteScript(context.Background(), createdRule.LuaScript, execContext)

	// The script should execute without syntax errors
	assert.True(t, result.Success, "Rule execution should succeed")
	assert.NotEmpty(t, result.Output, "Should have output")
	assert.Greater(t, result.Duration, time.Duration(0), "Execution should take some time")
}

func TestIntegration_ErrorHandling_InvalidAuth(t *testing.T) {
	server, cleanup := setupIntegrationTest(t)
	defer cleanup()

	// Try to access API without auth
	req := httptest.NewRequest(http.MethodGet, "/api/v1/rules", nil)
	w := httptest.NewRecorder()

	server.Router().ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// Try with invalid API key
	req2 := httptest.NewRequest(http.MethodGet, "/api/v1/rules", nil)
	req2.Header.Set("Authorization", "ApiKey invalid-key")
	w2 := httptest.NewRecorder()

	server.Router().ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusUnauthorized, w2.Code)
}

func TestIntegration_ErrorHandling_InvalidRuleCreation(t *testing.T) {
	server, cleanup := setupIntegrationTest(t)
	defer cleanup()

	// Test empty name
	createReq := CreateRuleRequest{
		Name:      "",
		LuaScript: "return true",
		Priority:  &[]int{0}[0],
		Enabled:   &[]bool{true}[0],
	}

	body, err := json.Marshal(createReq)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/rules", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "ApiKey test-api-key-integration")
	w := httptest.NewRecorder()

	server.Router().ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Test empty script
	createReq2 := CreateRuleRequest{
		Name:      "Valid Name",
		LuaScript: "",
		Priority:  &[]int{0}[0],
		Enabled:   &[]bool{true}[0],
	}

	body2, err := json.Marshal(createReq2)
	require.NoError(t, err)

	req2 := httptest.NewRequest(http.MethodPost, "/api/v1/rules", bytes.NewReader(body2))
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("Authorization", "ApiKey test-api-key-integration")
	w2 := httptest.NewRecorder()

	server.Router().ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusBadRequest, w2.Code)

	// Test malformed JSON
	req3 := httptest.NewRequest(http.MethodPost, "/api/v1/rules", bytes.NewReader([]byte("{invalid json")))
	req3.Header.Set("Content-Type", "application/json")
	req3.Header.Set("Authorization", "ApiKey test-api-key-integration")
	w3 := httptest.NewRecorder()

	server.Router().ServeHTTP(w3, req3)
	assert.Equal(t, http.StatusBadRequest, w3.Code)
}

func TestIntegration_ErrorHandling_RuleNotFound(t *testing.T) {
	server, cleanup := setupIntegrationTest(t)
	defer cleanup()

	// Try to get non-existent rule
	nonExistentID := uuid.New()
	getReq := httptest.NewRequest(http.MethodGet, "/api/v1/rules/"+nonExistentID.String(), nil)
	getReq.Header.Set("Authorization", "ApiKey test-api-key-integration")
	getW := httptest.NewRecorder()

	server.Router().ServeHTTP(getW, getReq)
	assert.Equal(t, http.StatusNotFound, getW.Code)

	// Try to update non-existent rule
	updateReq := UpdateRuleRequest{
		Name: &[]string{"New Name"}[0],
	}
	updateBody, err := json.Marshal(updateReq)
	require.NoError(t, err)

	updateReqHTTP := httptest.NewRequest(http.MethodPut, "/api/v1/rules/"+nonExistentID.String(), bytes.NewReader(updateBody))
	updateReqHTTP.Header.Set("Content-Type", "application/json")
	updateReqHTTP.Header.Set("Authorization", "ApiKey test-api-key-integration")
	updateW := httptest.NewRecorder()

	server.Router().ServeHTTP(updateW, updateReqHTTP)
	assert.Equal(t, http.StatusNotFound, updateW.Code)

	// Try to delete non-existent rule
	deleteReq := httptest.NewRequest(http.MethodDelete, "/api/v1/rules/"+nonExistentID.String(), nil)
	deleteReq.Header.Set("Authorization", "ApiKey test-api-key-integration")
	deleteW := httptest.NewRecorder()

	server.Router().ServeHTTP(deleteW, deleteReq)
	assert.Equal(t, http.StatusNotFound, deleteW.Code) // Delete should return 404 for non-existent resource
}

func TestIntegration_ErrorHandling_InvalidTriggerCreation(t *testing.T) {
	server, cleanup := setupIntegrationTest(t)
	defer cleanup()

	// Create a rule first
	ruleSvc := server.ruleSvc.(*rule.Service)
	testRule := &rule.Rule{
		Name:      "Test Rule for Invalid Trigger",
		LuaScript: "return true",
		Enabled:   true,
	}
	err := ruleSvc.Create(context.Background(), testRule)
	require.NoError(t, err)

	// Test invalid trigger type
	createReq := CreateTriggerRequest{
		RuleID:          testRule.ID,
		Type:            "invalid_type",
		ConditionScript: "event.value > 10",
		Enabled:         &[]bool{true}[0],
	}

	body, err := json.Marshal(createReq)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/triggers", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "ApiKey test-api-key-integration")
	w := httptest.NewRecorder()

	server.Router().ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Test empty condition script
	createReq2 := CreateTriggerRequest{
		RuleID:          testRule.ID,
		Type:            "conditional",
		ConditionScript: "",
		Enabled:         &[]bool{true}[0],
	}

	body2, err := json.Marshal(createReq2)
	require.NoError(t, err)

	req2 := httptest.NewRequest(http.MethodPost, "/api/v1/triggers", bytes.NewReader(body2))
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("Authorization", "ApiKey test-api-key-integration")
	w2 := httptest.NewRecorder()

	server.Router().ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusBadRequest, w2.Code)
}

func TestIntegration_ErrorHandling_TriggerNotFound(t *testing.T) {
	server, cleanup := setupIntegrationTest(t)
	defer cleanup()

	// Try to get non-existent trigger
	nonExistentID := uuid.New()
	getReq := httptest.NewRequest(http.MethodGet, "/api/v1/triggers/"+nonExistentID.String(), nil)
	getReq.Header.Set("Authorization", "ApiKey test-api-key-integration")
	getW := httptest.NewRecorder()

	server.Router().ServeHTTP(getW, getReq)
	assert.Equal(t, http.StatusNotFound, getW.Code)
}

func TestIntegration_ErrorHandling_InvalidActionCreation(t *testing.T) {
	server, cleanup := setupIntegrationTest(t)
	defer cleanup()

	// Test empty script
	createReq := CreateActionRequest{
		LuaScript: "",
		Enabled:   &[]bool{true}[0],
	}

	body, err := json.Marshal(createReq)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/actions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "ApiKey test-api-key-integration")
	w := httptest.NewRecorder()

	server.Router().ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestIntegration_ErrorHandling_ActionNotFound(t *testing.T) {
	server, cleanup := setupIntegrationTest(t)
	defer cleanup()

	// Try to get non-existent action
	nonExistentID := uuid.New()
	getReq := httptest.NewRequest(http.MethodGet, "/api/v1/actions/"+nonExistentID.String(), nil)
	getReq.Header.Set("Authorization", "ApiKey test-api-key-integration")
	getW := httptest.NewRecorder()

	server.Router().ServeHTTP(getW, getReq)
	assert.Equal(t, http.StatusNotFound, getW.Code)
}

func TestIntegration_CreateTrigger(t *testing.T) {
	server, cleanup := setupIntegrationTest(t)
	defer cleanup()

	// First create a rule
	ruleSvc := server.ruleSvc.(*rule.Service)
	testRule := &rule.Rule{
		Name:      "Test Rule for Trigger",
		LuaScript: "return true",
		Enabled:   true,
	}
	err := ruleSvc.Create(context.Background(), testRule)
	require.NoError(t, err)

	// Create a trigger
	createReq := CreateTriggerRequest{
		RuleID:          testRule.ID,
		Type:            "conditional",
		ConditionScript: "event.type == 'temperature_update'",
		Enabled:         &[]bool{true}[0],
	}

	body, err := json.Marshal(createReq)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/triggers", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "ApiKey test-api-key-integration")
	w := httptest.NewRecorder()

	server.Router().ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_CreateAction(t *testing.T) {
	server, cleanup := setupIntegrationTest(t)
	defer cleanup()

	// Create an action
	createReq := CreateActionRequest{
		LuaScript: "send_notification('alert')",
		Enabled:   &[]bool{true}[0],
	}

	body, err := json.Marshal(createReq)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/actions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "ApiKey test-api-key-integration")
	w := httptest.NewRecorder()

	server.Router().ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_TriggerEvaluation(t *testing.T) {
	server, cleanup := setupIntegrationTest(t)
	defer cleanup()

	// Create a rule that will be triggered
	ruleSvc := server.ruleSvc.(*rule.Service)
	testRule := &rule.Rule{
		Name:      "Trigger Evaluation Test Rule",
		LuaScript: "return event.temperature > 25",
		Priority:  5,
		Enabled:   true,
	}
	err := ruleSvc.Create(context.Background(), testRule)
	require.NoError(t, err)

	// Create a conditional trigger with a condition script
	triggerSvc := server.triggerSvc.(*trigger.Service)
	testTrigger := &trigger.Trigger{
		RuleID:          testRule.ID,
		Type:            trigger.Conditional,
		ConditionScript: "return event.type == 'temperature' and event.value > 20",
		Enabled:         true,
	}
	err = triggerSvc.Create(context.Background(), testTrigger)
	require.NoError(t, err)

	// Test trigger evaluation directly using the evaluator
	// Create evaluator with the same components used by the server
	contextSvc := execCtx.NewService()
	platformSvc := platform.NewService()
	executorSvc := executor.NewService(contextSvc, platformSvc)
	evaluator := trigger.NewEvaluator(executorSvc)

	// Test case 1: Event that should match the trigger condition
	eventData1 := map[string]any{
		"type":  "temperature",
		"value": 28.5,
	}
	results1 := evaluator.EvaluateTriggers(context.Background(), []*trigger.Trigger{testTrigger}, eventData1)
	assert.Len(t, results1, 1)
	assert.True(t, results1[0].Matched)
	assert.Equal(t, testTrigger.ID, results1[0].TriggerID)
	assert.Equal(t, testRule.ID, results1[0].RuleID)

	// Test case 2: Event that should NOT match the trigger condition
	eventData2 := map[string]any{
		"type":  "humidity",
		"value": 28.5,
	}
	results2 := evaluator.EvaluateTriggers(context.Background(), []*trigger.Trigger{testTrigger}, eventData2)
	assert.Len(t, results2, 1)
	assert.False(t, results2[0].Matched)

	// Test case 3: Event with wrong value
	eventData3 := map[string]any{
		"type":  "temperature",
		"value": 15.0,
	}
	results3 := evaluator.EvaluateTriggers(context.Background(), []*trigger.Trigger{testTrigger}, eventData3)
	assert.Len(t, results3, 1)
	assert.False(t, results3[0].Matched)

	// Test case 4: Complex condition with multiple criteria
	complexTrigger := &trigger.Trigger{
		RuleID:          testRule.ID,
		Type:            trigger.Conditional,
		ConditionScript: "return event.device_id == 'sensor_1' and event.temperature > 25 and event.status == 'active'",
		Enabled:         true,
	}

	eventData4 := map[string]any{
		"device_id":   "sensor_1",
		"temperature": 28.0,
		"status":      "active",
	}
	results4 := evaluator.EvaluateTriggers(context.Background(), []*trigger.Trigger{complexTrigger}, eventData4)
	assert.Len(t, results4, 1)
	assert.True(t, results4[0].Matched)
}

func TestIntegration_GetDashboardData(t *testing.T) {
	server, cleanup := setupIntegrationTest(t)
	defer cleanup()

	// Test with default time range
	req := httptest.NewRequest(http.MethodGet, "/api/v1/analytics/dashboard", nil)
	req.Header.Set("Authorization", "ApiKey test-api-key-integration")
	w := httptest.NewRecorder()

	server.Router().ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var data analytics.DashboardData
	err := json.Unmarshal(w.Body.Bytes(), &data)
	require.NoError(t, err)

	assert.NotNil(t, data.OverallStats)
	assert.NotNil(t, data.RuleStats)
	assert.NotNil(t, data.ExecutionTrend)
	assert.NotNil(t, data.SuccessRateTrend)
	assert.NotNil(t, data.LatencyTrend)
	assert.Equal(t, "24h", data.TimeRange)

	// Test with custom time range
	req2 := httptest.NewRequest(http.MethodGet, "/api/v1/analytics/dashboard?timeRange=1h", nil)
	req2.Header.Set("Authorization", "ApiKey test-api-key-integration")
	w2 := httptest.NewRecorder()

	server.Router().ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)

	var data2 analytics.DashboardData
	err = json.Unmarshal(w2.Body.Bytes(), &data2)
	require.NoError(t, err)
	assert.Equal(t, "1h", data2.TimeRange)

	// Test with invalid time range
	req3 := httptest.NewRequest(http.MethodGet, "/api/v1/analytics/dashboard?timeRange=invalid", nil)
	req3.Header.Set("Authorization", "ApiKey test-api-key-integration")
	w3 := httptest.NewRecorder()

	server.Router().ServeHTTP(w3, req3)
	assert.Equal(t, http.StatusBadRequest, w3.Code)
}
