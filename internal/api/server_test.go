package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/malyshevhen/rule-engine/internal/action"
	"github.com/malyshevhen/rule-engine/internal/engine/executor"
	execCtx "github.com/malyshevhen/rule-engine/internal/engine/executor/context"
	"github.com/malyshevhen/rule-engine/internal/rule"
	actionStorage "github.com/malyshevhen/rule-engine/internal/storage/action"
	ruleStorage "github.com/malyshevhen/rule-engine/internal/storage/rule"
	triggerStorage "github.com/malyshevhen/rule-engine/internal/storage/trigger"
	"github.com/malyshevhen/rule-engine/internal/trigger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockRuleService is a mock implementation of RuleService
type mockRuleService struct {
	mock.Mock
}

func (m *mockRuleService) Create(ctx context.Context, rule *rule.Rule) error {
	args := m.Called(ctx, rule)
	return args.Error(0)
}

func (m *mockRuleService) GetByID(ctx context.Context, id uuid.UUID) (*rule.Rule, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*rule.Rule), args.Error(1)
}

func (m *mockRuleService) List(ctx context.Context, limit int, offset int) ([]*rule.Rule, int, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]*rule.Rule), args.Int(1), args.Error(2)
}

func (m *mockRuleService) ListAll(ctx context.Context) ([]*rule.Rule, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*rule.Rule), args.Error(1)
}

func (m *mockRuleService) Update(ctx context.Context, rule *rule.Rule) error {
	args := m.Called(ctx, rule)
	return args.Error(0)
}

func (m *mockRuleService) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockRuleService) AddAction(ctx context.Context, ruleID, actionID uuid.UUID) error {
	args := m.Called(ctx, ruleID, actionID)
	return args.Error(0)
}

// mockTriggerService is a mock implementation of TriggerService
type mockTriggerService struct {
	mock.Mock
}

func (m *mockTriggerService) Create(ctx context.Context, trigger *trigger.Trigger) error {
	args := m.Called(ctx, trigger)
	return args.Error(0)
}

func (m *mockTriggerService) GetByID(ctx context.Context, id uuid.UUID) (*trigger.Trigger, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*trigger.Trigger), args.Error(1)
}

func (m *mockTriggerService) List(ctx context.Context, limit, offset int) ([]*trigger.Trigger, int, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]*trigger.Trigger), args.Int(1), args.Error(2)
}

func (m *mockTriggerService) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// mockActionService is a mock implementation of ActionService
type mockActionService struct {
	mock.Mock
}

func (m *mockActionService) Create(ctx context.Context, action *action.Action) error {
	args := m.Called(ctx, action)
	return args.Error(0)
}

func (m *mockActionService) GetByID(ctx context.Context, id uuid.UUID) (*action.Action, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*action.Action), args.Error(1)
}

func (m *mockActionService) List(ctx context.Context, limit, offset int) ([]*action.Action, int, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]*action.Action), args.Int(1), args.Error(2)
}

func (m *mockActionService) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// mockExecutorService is a mock implementation of ExecutorService
type mockExecutorService struct {
	mock.Mock
}

func (m *mockExecutorService) ExecuteScript(ctx context.Context, script string, execCtx *execCtx.ExecutionContext) *executor.ExecuteResult {
	args := m.Called(ctx, script, execCtx)
	return args.Get(0).(*executor.ExecuteResult)
}

func TestServer_CreateRule(t *testing.T) {
	mockRuleSvc := &mockRuleService{}

	tests := []struct {
		name           string
		requestBody    CreateRuleRequest
		expectedStatus int
		setupMocks     func()
	}{
		{
			name: "successful creation",
			requestBody: CreateRuleRequest{
				Name:      "Test Rule",
				LuaScript: "return true",
				Priority:  &[]int{5}[0],
				Enabled:   &[]bool{true}[0],
			},
			expectedStatus: http.StatusCreated,
			setupMocks: func() {
				mockRuleSvc.On("Create", mock.Anything, mock.MatchedBy(func(r *rule.Rule) bool {
					return r.Name == "Test Rule" && r.LuaScript == "return true" && r.Priority == 5 && r.Enabled == true
				})).Return(nil)
			},
		},
		{
			name: "empty name",
			requestBody: CreateRuleRequest{
				Name:      "",
				LuaScript: "return true",
			},
			expectedStatus: http.StatusBadRequest,
			setupMocks:     func() {},
		},
		{
			name: "empty lua script",
			requestBody: CreateRuleRequest{
				Name:      "Test Rule",
				LuaScript: "",
			},
			expectedStatus: http.StatusBadRequest,
			setupMocks:     func() {},
		},
		{
			name: "name too long",
			requestBody: CreateRuleRequest{
				Name:      string(make([]byte, 256)),
				LuaScript: "return true",
			},
			expectedStatus: http.StatusBadRequest,
			setupMocks:     func() {},
		},
		{
			name: "lua script too long",
			requestBody: CreateRuleRequest{
				Name:      "Test Rule",
				LuaScript: string(make([]byte, 10001)),
			},
			expectedStatus: http.StatusBadRequest,
			setupMocks:     func() {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/rules", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			createRule(mockRuleSvc)(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockRuleSvc.AssertExpectations(t)
		})
	}
}

func TestServer_ListRules(t *testing.T) {
	mockRuleSvc := &mockRuleService{}

	expectedRules := []*rule.Rule{
		{
			ID:        uuid.New(),
			Name:      "Rule 1",
			LuaScript: "return true",
			Enabled:   true,
		},
		{
			ID:        uuid.New(),
			Name:      "Rule 2",
			LuaScript: "return false",
			Enabled:   false,
		},
	}

	mockRuleSvc.On("List", mock.Anything, 50, 0).Return(expectedRules, 2, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/rules", nil)
	w := httptest.NewRecorder()

	listRules(mockRuleSvc)(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockRuleSvc.AssertExpectations(t)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	rules := response["rules"].([]any)
	assert.Len(t, rules, 2)

	// Check the first rule
	rule1 := rules[0].(map[string]any)
	assert.Equal(t, "Rule 1", rule1["name"])

	// Check the second rule
	rule2 := rules[1].(map[string]any)
	assert.Equal(t, "Rule 2", rule2["name"])

	// Check pagination metadata
	assert.Equal(t, float64(50), response["limit"])
	assert.Equal(t, float64(0), response["offset"])
	assert.Equal(t, float64(2), response["count"])
}

func TestServer_UpdateRule(t *testing.T) {
	mockRuleSvc := &mockRuleService{}

	ruleID := uuid.New()
	existingRule := &rule.Rule{
		ID:        ruleID,
		Name:      "Old Name",
		LuaScript: "return false",
		Priority:  0,
		Enabled:   false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name           string
		ruleID         string
		requestBody    PatchRequest
		expectedStatus int
		setupMocks     func()
	}{
		{
			name:   "successful update with replace operations",
			ruleID: ruleID.String(),
			requestBody: PatchRequest{
				{Op: "replace", Path: "/name", Value: "New Name"},
				{Op: "replace", Path: "/lua_script", Value: "return true"},
				{Op: "replace", Path: "/priority", Value: 10},
				{Op: "replace", Path: "/enabled", Value: true},
			},
			expectedStatus: http.StatusOK,
			setupMocks: func() {
				mockRuleSvc.On("GetByID", mock.Anything, ruleID).Return(existingRule, nil)
				mockRuleSvc.On("Update", mock.Anything, mock.AnythingOfType("*rule.Rule")).Return(nil)
			},
		},
		{
			name:   "successful update with single operation",
			ruleID: ruleID.String(),
			requestBody: PatchRequest{
				{Op: "replace", Path: "/name", Value: "Updated Name"},
			},
			expectedStatus: http.StatusOK,
			setupMocks: func() {
				mockRuleSvc.On("GetByID", mock.Anything, ruleID).Return(existingRule, nil)
				mockRuleSvc.On("Update", mock.Anything, mock.AnythingOfType("*rule.Rule")).Return(nil)
			},
		},
		{
			name:   "invalid uuid",
			ruleID: "invalid-uuid",
			requestBody: PatchRequest{
				{Op: "replace", Path: "/name", Value: "New Name"},
			},
			expectedStatus: http.StatusBadRequest,
			setupMocks:     func() {},
		},
		{
			name:   "rule not found",
			ruleID: uuid.New().String(),
			requestBody: PatchRequest{
				{Op: "replace", Path: "/name", Value: "New Name"},
			},
			expectedStatus: http.StatusNotFound,
			setupMocks: func() {
				mockRuleSvc.On("GetByID", mock.Anything, mock.Anything).Return((*rule.Rule)(nil), ruleStorage.ErrNotFound)
			},
		},
		{
			name:   "empty name after patch",
			ruleID: ruleID.String(),
			requestBody: PatchRequest{
				{Op: "replace", Path: "/name", Value: ""},
			},
			expectedStatus: http.StatusBadRequest,
			setupMocks: func() {
				mockRuleSvc.On("GetByID", mock.Anything, ruleID).Return(existingRule, nil)
			},
		},
		{
			name:   "invalid patch operation",
			ruleID: ruleID.String(),
			requestBody: PatchRequest{
				{Op: "invalid", Path: "/name", Value: "New Name"},
			},
			expectedStatus: http.StatusBadRequest,
			setupMocks: func() {
				mockRuleSvc.On("GetByID", mock.Anything, ruleID).Return(existingRule, nil)
			},
		},
		{
			name:   "invalid patch path",
			ruleID: ruleID.String(),
			requestBody: PatchRequest{
				{Op: "replace", Path: "/invalid_field", Value: "value"},
			},
			expectedStatus: http.StatusBadRequest,
			setupMocks: func() {
				mockRuleSvc.On("GetByID", mock.Anything, ruleID).Return(existingRule, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPatch, "/api/v1/rules/"+tt.ruleID, bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json-patch+json")
			req = mux.SetURLVars(req, map[string]string{"id": tt.ruleID})
			w := httptest.NewRecorder()

			updateRule(mockRuleSvc)(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockRuleSvc.AssertExpectations(t)
		})
	}
}

func TestServer_DeleteRule(t *testing.T) {
	mockRuleSvc := &mockRuleService{}

	ruleID := uuid.New()

	tests := []struct {
		name           string
		ruleID         string
		expectedStatus int
		setupMocks     func()
	}{
		{
			name:           "successful delete",
			ruleID:         ruleID.String(),
			expectedStatus: http.StatusNoContent,
			setupMocks: func() {
				mockRuleSvc.On("Delete", mock.Anything, ruleID).Return(nil)
			},
		},
		{
			name:           "invalid uuid",
			ruleID:         "invalid-uuid",
			expectedStatus: http.StatusBadRequest,
			setupMocks:     func() {},
		},
		{
			name:           "rule not found",
			ruleID:         uuid.New().String(),
			expectedStatus: http.StatusNotFound,
			setupMocks: func() {
				mockRuleSvc.On("Delete", mock.Anything, mock.Anything).Return(ruleStorage.ErrNotFound)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			req := httptest.NewRequest(http.MethodDelete, "/api/v1/rules/"+tt.ruleID, nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.ruleID})
			w := httptest.NewRecorder()

			deleteRule(mockRuleSvc)(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockRuleSvc.AssertExpectations(t)
		})
	}
}

func TestServer_GetRule(t *testing.T) {
	mockRuleSvc := &mockRuleService{}

	ruleID := uuid.New()
	expectedRule := &rule.Rule{
		ID:        ruleID,
		Name:      "Test Rule",
		LuaScript: "return true",
		Enabled:   true,
	}

	tests := []struct {
		name           string
		ruleID         string
		expectedStatus int
		setupMocks     func()
	}{
		{
			name:           "successful get",
			ruleID:         ruleID.String(),
			expectedStatus: http.StatusOK,
			setupMocks: func() {
				mockRuleSvc.On("GetByID", mock.Anything, ruleID).Return(expectedRule, nil)
			},
		},
		{
			name:           "invalid uuid",
			ruleID:         "invalid-uuid",
			expectedStatus: http.StatusBadRequest,
			setupMocks:     func() {},
		},
		{
			name:           "rule not found",
			ruleID:         uuid.New().String(),
			expectedStatus: http.StatusNotFound,
			setupMocks: func() {
				mockRuleSvc.On("GetByID", mock.Anything, mock.Anything).Return((*rule.Rule)(nil), ruleStorage.ErrNotFound)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			req := httptest.NewRequest(http.MethodGet, "/api/v1/rules/"+tt.ruleID, nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.ruleID})
			w := httptest.NewRecorder()

			getRule(mockRuleSvc)(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockRuleSvc.AssertExpectations(t)
		})
	}
}

func TestServer_AddActionToRule(t *testing.T) {
	ruleID := uuid.New()
	actionID := uuid.New()

	tests := []struct {
		name           string
		ruleID         string
		requestBody    AddActionToRuleRequest
		expectedStatus int
		setupMocks     func(mock *mockRuleService)
	}{
		{
			name:   "successful add action",
			ruleID: ruleID.String(),
			requestBody: AddActionToRuleRequest{
				ActionID: actionID,
			},
			expectedStatus: http.StatusNoContent,
			setupMocks: func(m *mockRuleService) {
				m.On("AddAction", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
		},
		{
			name:   "invalid rule uuid",
			ruleID: "invalid-uuid",
			requestBody: AddActionToRuleRequest{
				ActionID: actionID,
			},
			expectedStatus: http.StatusBadRequest,
			setupMocks:     func(mock *mockRuleService) {},
		},
		{
			name:   "rule not found",
			ruleID: uuid.New().String(),
			requestBody: AddActionToRuleRequest{
				ActionID: actionID,
			},
			expectedStatus: http.StatusNotFound,
			setupMocks: func(m *mockRuleService) {
				m.On("AddAction", mock.Anything, mock.Anything, mock.Anything).Return(ruleStorage.ErrNotFound)
			},
		},
		{
			name:   "action not found",
			ruleID: ruleID.String(),
			requestBody: AddActionToRuleRequest{
				ActionID: uuid.New(),
			},
			expectedStatus: http.StatusBadRequest,
			setupMocks: func(m *mockRuleService) {
				m.On("AddAction", mock.Anything, mock.Anything, mock.Anything).Return(actionStorage.ErrNotFound)
			},
		},
		{
			name:   "service error",
			ruleID: ruleID.String(),
			requestBody: AddActionToRuleRequest{
				ActionID: actionID,
			},
			expectedStatus: http.StatusInternalServerError,
			setupMocks: func(m *mockRuleService) {
				m.On("AddAction", mock.Anything, mock.Anything, mock.Anything).Return(assert.AnError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRuleSvc := &mockRuleService{}
			tt.setupMocks(mockRuleSvc)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/rules/"+tt.ruleID+"/actions", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			req = mux.SetURLVars(req, map[string]string{"id": tt.ruleID})
			w := httptest.NewRecorder()

			addActionToRule(mockRuleSvc)(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockRuleSvc.AssertExpectations(t)
		})
	}
}

func TestServer_CreateTrigger(t *testing.T) {
	mockTriggerSvc := &mockTriggerService{}

	ruleID := uuid.New()

	tests := []struct {
		name           string
		requestBody    CreateTriggerRequest
		expectedStatus int
		setupMocks     func()
	}{
		{
			name: "successful creation",
			requestBody: CreateTriggerRequest{
				RuleID:          ruleID,
				Type:            "CONDITIONAL",
				ConditionScript: "event.type == 'device_update'",
				Enabled:         &[]bool{true}[0],
			},
			expectedStatus: http.StatusCreated,
			setupMocks: func() {
				mockTriggerSvc.On("Create", mock.Anything, mock.MatchedBy(func(tr *trigger.Trigger) bool {
					return tr.RuleID == ruleID && tr.Type == trigger.Conditional && tr.ConditionScript == "event.type == 'device_update'" && tr.Enabled == true
				})).Return(nil)
			},
		},
		{
			name: "empty type",
			requestBody: CreateTriggerRequest{
				RuleID:          ruleID,
				Type:            "",
				ConditionScript: "event.type == 'device_update'",
			},
			expectedStatus: http.StatusBadRequest,
			setupMocks:     func() {},
		},
		{
			name: "invalid type",
			requestBody: CreateTriggerRequest{
				RuleID:          ruleID,
				Type:            "invalid",
				ConditionScript: "event.type == 'device_update'",
			},
			expectedStatus: http.StatusBadRequest,
			setupMocks:     func() {},
		},
		{
			name: "empty condition script",
			requestBody: CreateTriggerRequest{
				RuleID:          ruleID,
				Type:            "conditional",
				ConditionScript: "",
			},
			expectedStatus: http.StatusBadRequest,
			setupMocks:     func() {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/triggers", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			createTrigger(mockTriggerSvc)(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockTriggerSvc.AssertExpectations(t)
		})
	}
}

func TestServer_ListTriggers(t *testing.T) {
	mockTriggerSvc := &mockTriggerService{}

	expectedTriggers := []*trigger.Trigger{
		{
			ID:              uuid.New(),
			Type:            trigger.Conditional,
			ConditionScript: "event.type == 'device_update'",
			Enabled:         true,
		},
	}

	mockTriggerSvc.On("List", mock.Anything, 50, 0).Return(expectedTriggers, 1, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/triggers", nil)
	w := httptest.NewRecorder()

	listTriggers(mockTriggerSvc)(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockTriggerSvc.AssertExpectations(t)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	triggers := response["triggers"].([]any)
	assert.Len(t, triggers, 1)

	triggerData := triggers[0].(map[string]any)
	assert.Equal(t, "CONDITIONAL", triggerData["type"])
}

func TestServer_GetTrigger(t *testing.T) {
	mockTriggerSvc := &mockTriggerService{}

	triggerID := uuid.New()
	expectedTrigger := &trigger.Trigger{
		ID:              triggerID,
		Type:            trigger.Conditional,
		ConditionScript: "event.type == 'device_update'",
		Enabled:         true,
	}

	tests := []struct {
		name           string
		triggerID      string
		expectedStatus int
		setupMocks     func()
	}{
		{
			name:           "successful get",
			triggerID:      triggerID.String(),
			expectedStatus: http.StatusOK,
			setupMocks: func() {
				mockTriggerSvc.On("GetByID", mock.Anything, triggerID).Return(expectedTrigger, nil)
			},
		},
		{
			name:           "invalid uuid",
			triggerID:      "invalid-uuid",
			expectedStatus: http.StatusBadRequest,
			setupMocks:     func() {},
		},
		{
			name:           "trigger not found",
			triggerID:      uuid.New().String(),
			expectedStatus: http.StatusNotFound,
			setupMocks: func() {
				mockTriggerSvc.On("GetByID", mock.Anything, mock.Anything).Return((*trigger.Trigger)(nil), triggerStorage.ErrNotFound)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			req := httptest.NewRequest(http.MethodGet, "/api/v1/triggers/"+tt.triggerID, nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.triggerID})
			w := httptest.NewRecorder()

			getTrigger(mockTriggerSvc)(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockTriggerSvc.AssertExpectations(t)
		})
	}
}

func TestServer_ListActions(t *testing.T) {
	mockActionSvc := &mockActionService{}

	expectedActions := []*action.Action{
		{
			ID:        uuid.New(),
			LuaScript: "print('action 1')",
			Enabled:   true,
		},
		{
			ID:        uuid.New(),
			LuaScript: "print('action 2')",
			Enabled:   false,
		},
	}

	mockActionSvc.On("List", mock.Anything, 50, 0).Return(expectedActions, 2, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/actions", nil)
	w := httptest.NewRecorder()

	listActions(mockActionSvc)(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockActionSvc.AssertExpectations(t)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	actions := response["actions"].([]any)
	assert.Len(t, actions, 2)

	actionData1 := actions[0].(map[string]any)
	actionData2 := actions[1].(map[string]any)
	assert.Equal(t, "print('action 1')", actionData1["lua_script"])
	assert.Equal(t, "print('action 2')", actionData2["lua_script"])
}

func TestServer_GetAction(t *testing.T) {
	mockActionSvc := &mockActionService{}

	actionID := uuid.New()
	expectedAction := &action.Action{
		ID:        actionID,
		LuaScript: "print('test action')",
		Enabled:   true,
	}

	tests := []struct {
		name           string
		actionID       string
		expectedStatus int
		setupMocks     func()
	}{
		{
			name:           "successful get",
			actionID:       actionID.String(),
			expectedStatus: http.StatusOK,
			setupMocks: func() {
				mockActionSvc.On("GetByID", mock.Anything, actionID).Return(expectedAction, nil)
			},
		},
		{
			name:           "invalid uuid",
			actionID:       "invalid-uuid",
			expectedStatus: http.StatusBadRequest,
			setupMocks:     func() {},
		},
		{
			name:           "action not found",
			actionID:       uuid.New().String(),
			expectedStatus: http.StatusNotFound,
			setupMocks: func() {
				mockActionSvc.On("GetByID", mock.Anything, mock.Anything).Return((*action.Action)(nil), actionStorage.ErrNotFound)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			req := httptest.NewRequest(http.MethodGet, "/api/v1/actions/"+tt.actionID, nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.actionID})
			w := httptest.NewRecorder()

			getAction(mockActionSvc)(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockActionSvc.AssertExpectations(t)
		})
	}
}

func TestServer_HealthCheck(t *testing.T) {
	// TODO: Implement health check test with proper mocking of database and redis clients
	// This test should verify the /health endpoint returns correct status based on service health
	// Mock database and redis health checks
	// Make GET request to /health
	// Assert response status and body
	t.Skip("Health check test requires complex mocking of database and redis clients. Integration tests cover this functionality.")
}

func TestServer_CreateAction(t *testing.T) {
	mockActionSvc := &mockActionService{}

	tests := []struct {
		name           string
		requestBody    CreateActionRequest
		expectedStatus int
		setupMocks     func()
	}{
		{
			name: "successful creation",
			requestBody: CreateActionRequest{
				LuaScript: "print('action executed')",
				Enabled:   &[]bool{true}[0],
			},
			expectedStatus: http.StatusCreated,
			setupMocks: func() {
				mockActionSvc.On("Create", mock.Anything, mock.MatchedBy(func(a *action.Action) bool {
					return a.LuaScript == "print('action executed')" && a.Enabled == true
				})).Return(nil)
			},
		},
		{
			name: "empty lua script",
			requestBody: CreateActionRequest{
				LuaScript: "",
			},
			expectedStatus: http.StatusBadRequest,
			setupMocks:     func() {},
		},
		{
			name: "lua script too long",
			requestBody: CreateActionRequest{
				LuaScript: string(make([]byte, 10001)),
			},
			expectedStatus: http.StatusBadRequest,
			setupMocks:     func() {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/actions", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			createAction(mockActionSvc)(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockActionSvc.AssertExpectations(t)
		})
	}
}

func TestServer_EvaluateScript(t *testing.T) {
	mockExecutorSvc := &mockExecutorService{}

	tests := []struct {
		name             string
		requestBody      EvaluateScriptRequest
		expectedStatus   int
		setupMocks       func()
		expectedResponse EvaluateScriptResponse
	}{
		{
			name: "successful evaluation",
			requestBody: EvaluateScriptRequest{
				Script: "return 2 + 2",
			},
			expectedStatus: http.StatusOK,
			setupMocks: func() {
				result := &executor.ExecuteResult{
					Success:  true,
					Output:   []any{4.0},
					Duration: 100 * time.Millisecond,
				}
				mockExecutorSvc.On("ExecuteScript", mock.Anything, "return 2 + 2", mock.AnythingOfType("*context.ExecutionContext")).Return(result)
			},
			expectedResponse: EvaluateScriptResponse{
				Success:  true,
				Output:   []any{4.0},
				Duration: "100ms",
			},
		},
		{
			name: "script execution error",
			requestBody: EvaluateScriptRequest{
				Script: "invalid lua",
			},
			expectedStatus: http.StatusOK,
			setupMocks: func() {
				result := &executor.ExecuteResult{
					Success:  false,
					Error:    "syntax error",
					Duration: 50 * time.Millisecond,
				}
				mockExecutorSvc.On("ExecuteScript", mock.Anything, "invalid lua", mock.AnythingOfType("*context.ExecutionContext")).Return(result)
			},
			expectedResponse: EvaluateScriptResponse{
				Success:  false,
				Error:    "syntax error",
				Duration: "50ms",
			},
		},
		{
			name: "empty script",
			requestBody: EvaluateScriptRequest{
				Script: "",
			},
			expectedStatus: http.StatusBadRequest,
			setupMocks:     func() {},
		},
		{
			name: "script too long",
			requestBody: EvaluateScriptRequest{
				Script: string(make([]byte, 10001)),
			},
			expectedStatus: http.StatusBadRequest,
			setupMocks:     func() {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/evaluate", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			evaluateScript(mockExecutorSvc)(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var response EvaluateScriptResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResponse.Success, response.Success)
				assert.Equal(t, tt.expectedResponse.Error, response.Error)
				assert.Equal(t, tt.expectedResponse.Duration, response.Duration)
				if tt.expectedResponse.Output != nil {
					assert.Equal(t, tt.expectedResponse.Output, response.Output)
				}
			}

			mockExecutorSvc.AssertExpectations(t)
		})
	}
}
