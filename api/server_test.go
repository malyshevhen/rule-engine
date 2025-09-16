package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/malyshevhen/rule-engine/internal/core/action"
	"github.com/malyshevhen/rule-engine/internal/core/rule"
	"github.com/malyshevhen/rule-engine/internal/core/trigger"
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

func (m *mockRuleService) List(ctx context.Context, limit int, offset int) ([]*rule.Rule, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]*rule.Rule), args.Error(1)
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

func (m *mockTriggerService) List(ctx context.Context) ([]*trigger.Trigger, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*trigger.Trigger), args.Error(1)
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

func (m *mockActionService) List(ctx context.Context) ([]*action.Action, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*action.Action), args.Error(1)
}

func TestServer_CreateRule(t *testing.T) {
	mockRuleSvc := &mockRuleService{}
	mockTriggerSvc := &mockTriggerService{}
	mockActionSvc := &mockActionService{}

	server := &Server{
		ruleSvc:    mockRuleSvc,
		triggerSvc: mockTriggerSvc,
		actionSvc:  mockActionSvc,
	}

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
				Enabled:   &[]bool{true}[0],
			},
			expectedStatus: http.StatusOK,
			setupMocks: func() {
				mockRuleSvc.On("Create", mock.Anything, mock.MatchedBy(func(r *rule.Rule) bool {
					return r.Name == "Test Rule" && r.LuaScript == "return true" && r.Enabled == true
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

			server.CreateRule(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockRuleSvc.AssertExpectations(t)
		})
	}
}

func TestServer_ListRules(t *testing.T) {
	mockRuleSvc := &mockRuleService{}
	mockTriggerSvc := &mockTriggerService{}
	mockActionSvc := &mockActionService{}

	server := &Server{
		ruleSvc:    mockRuleSvc,
		triggerSvc: mockTriggerSvc,
		actionSvc:  mockActionSvc,
	}

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

	mockRuleSvc.On("ListAll", mock.Anything).Return(expectedRules, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/rules", nil)
	w := httptest.NewRecorder()

	server.ListRules(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockRuleSvc.AssertExpectations(t)

	var response []*rule.Rule
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)
	assert.Equal(t, "Rule 1", response[0].Name)
	assert.Equal(t, "Rule 2", response[1].Name)
}

func TestServer_UpdateRule(t *testing.T) {
	mockRuleSvc := &mockRuleService{}
	mockTriggerSvc := &mockTriggerService{}
	mockActionSvc := &mockActionService{}

	server := &Server{
		ruleSvc:    mockRuleSvc,
		triggerSvc: mockTriggerSvc,
		actionSvc:  mockActionSvc,
	}

	ruleID := uuid.New()
	existingRule := &rule.Rule{
		ID:        ruleID,
		Name:      "Old Name",
		LuaScript: "return false",
		Enabled:   false,
	}

	tests := []struct {
		name           string
		ruleID         string
		requestBody    UpdateRuleRequest
		expectedStatus int
		setupMocks     func()
	}{
		{
			name:   "successful update",
			ruleID: ruleID.String(),
			requestBody: UpdateRuleRequest{
				Name:      &[]string{"New Name"}[0],
				LuaScript: &[]string{"return true"}[0],
				Enabled:   &[]bool{true}[0],
			},
			expectedStatus: http.StatusOK,
			setupMocks: func() {
				mockRuleSvc.On("GetByID", mock.Anything, ruleID).Return(existingRule, nil)
				mockRuleSvc.On("Update", mock.Anything, mock.MatchedBy(func(r *rule.Rule) bool {
					return r.ID == ruleID && r.Name == "New Name" && r.LuaScript == "return true" && r.Enabled == true
				})).Return(nil)
			},
		},
		{
			name:   "invalid uuid",
			ruleID: "invalid-uuid",
			requestBody: UpdateRuleRequest{
				Name: &[]string{"New Name"}[0],
			},
			expectedStatus: http.StatusBadRequest,
			setupMocks:     func() {},
		},
		{
			name:   "rule not found",
			ruleID: uuid.New().String(),
			requestBody: UpdateRuleRequest{
				Name: &[]string{"New Name"}[0],
			},
			expectedStatus: http.StatusNotFound,
			setupMocks: func() {
				mockRuleSvc.On("GetByID", mock.Anything, mock.Anything).Return((*rule.Rule)(nil), assert.AnError)
			},
		},
		{
			name:   "empty name",
			ruleID: ruleID.String(),
			requestBody: UpdateRuleRequest{
				Name: &[]string{""}[0],
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
			req := httptest.NewRequest(http.MethodPut, "/api/v1/rules/"+tt.ruleID, bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			req = mux.SetURLVars(req, map[string]string{"id": tt.ruleID})
			w := httptest.NewRecorder()

			server.UpdateRule(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockRuleSvc.AssertExpectations(t)
		})
	}
}

func TestServer_DeleteRule(t *testing.T) {
	mockRuleSvc := &mockRuleService{}
	mockTriggerSvc := &mockTriggerService{}
	mockActionSvc := &mockActionService{}

	server := &Server{
		ruleSvc:    mockRuleSvc,
		triggerSvc: mockTriggerSvc,
		actionSvc:  mockActionSvc,
	}

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
			expectedStatus: http.StatusInternalServerError,
			setupMocks: func() {
				mockRuleSvc.On("Delete", mock.Anything, mock.Anything).Return(assert.AnError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			req := httptest.NewRequest(http.MethodDelete, "/api/v1/rules/"+tt.ruleID, nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.ruleID})
			w := httptest.NewRecorder()

			server.DeleteRule(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockRuleSvc.AssertExpectations(t)
		})
	}
}

func TestServer_GetRule(t *testing.T) {
	mockRuleSvc := &mockRuleService{}
	mockTriggerSvc := &mockTriggerService{}
	mockActionSvc := &mockActionService{}

	server := &Server{
		ruleSvc:    mockRuleSvc,
		triggerSvc: mockTriggerSvc,
		actionSvc:  mockActionSvc,
	}

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
				mockRuleSvc.On("GetByID", mock.Anything, mock.Anything).Return((*rule.Rule)(nil), assert.AnError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			req := httptest.NewRequest(http.MethodGet, "/api/v1/rules/"+tt.ruleID, nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.ruleID})
			w := httptest.NewRecorder()

			server.GetRule(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockRuleSvc.AssertExpectations(t)
		})
	}
}

func TestServer_CreateTrigger(t *testing.T) {
	mockRuleSvc := &mockRuleService{}
	mockTriggerSvc := &mockTriggerService{}
	mockActionSvc := &mockActionService{}

	server := &Server{
		ruleSvc:    mockRuleSvc,
		triggerSvc: mockTriggerSvc,
		actionSvc:  mockActionSvc,
	}

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
			expectedStatus: http.StatusOK,
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

			server.CreateTrigger(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockTriggerSvc.AssertExpectations(t)
		})
	}
}

func TestServer_ListTriggers(t *testing.T) {
	mockRuleSvc := &mockRuleService{}
	mockTriggerSvc := &mockTriggerService{}
	mockActionSvc := &mockActionService{}

	server := &Server{
		ruleSvc:    mockRuleSvc,
		triggerSvc: mockTriggerSvc,
		actionSvc:  mockActionSvc,
	}

	expectedTriggers := []*trigger.Trigger{
		{
			ID:              uuid.New(),
			Type:            trigger.Conditional,
			ConditionScript: "event.type == 'device_update'",
			Enabled:         true,
		},
	}

	mockTriggerSvc.On("List", mock.Anything).Return(expectedTriggers, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/triggers", nil)
	w := httptest.NewRecorder()

	server.ListTriggers(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockTriggerSvc.AssertExpectations(t)

	var response []*trigger.Trigger
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 1)
	assert.Equal(t, trigger.Conditional, response[0].Type)
}

func TestServer_GetTrigger(t *testing.T) {
	mockRuleSvc := &mockRuleService{}
	mockTriggerSvc := &mockTriggerService{}
	mockActionSvc := &mockActionService{}

	server := &Server{
		ruleSvc:    mockRuleSvc,
		triggerSvc: mockTriggerSvc,
		actionSvc:  mockActionSvc,
	}

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
				mockTriggerSvc.On("GetByID", mock.Anything, mock.Anything).Return((*trigger.Trigger)(nil), assert.AnError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			req := httptest.NewRequest(http.MethodGet, "/api/v1/triggers/"+tt.triggerID, nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.triggerID})
			w := httptest.NewRecorder()

			server.GetTrigger(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockTriggerSvc.AssertExpectations(t)
		})
	}
}

func TestServer_ListActions(t *testing.T) {
	mockRuleSvc := &mockRuleService{}
	mockTriggerSvc := &mockTriggerService{}
	mockActionSvc := &mockActionService{}

	server := &Server{
		ruleSvc:    mockRuleSvc,
		triggerSvc: mockTriggerSvc,
		actionSvc:  mockActionSvc,
	}

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

	mockActionSvc.On("List", mock.Anything).Return(expectedActions, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/actions", nil)
	w := httptest.NewRecorder()

	server.ListActions(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockActionSvc.AssertExpectations(t)

	var response []*action.Action
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)
	assert.Equal(t, "print('action 1')", response[0].LuaScript)
	assert.Equal(t, "print('action 2')", response[1].LuaScript)
}

func TestServer_GetAction(t *testing.T) {
	mockRuleSvc := &mockRuleService{}
	mockTriggerSvc := &mockTriggerService{}
	mockActionSvc := &mockActionService{}

	server := &Server{
		ruleSvc:    mockRuleSvc,
		triggerSvc: mockTriggerSvc,
		actionSvc:  mockActionSvc,
	}

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
				mockActionSvc.On("GetByID", mock.Anything, mock.Anything).Return((*action.Action)(nil), assert.AnError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			req := httptest.NewRequest(http.MethodGet, "/api/v1/actions/"+tt.actionID, nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.actionID})
			w := httptest.NewRecorder()

			server.GetAction(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockActionSvc.AssertExpectations(t)
		})
	}
}

func TestServer_HealthCheck(t *testing.T) {
	mockRuleSvc := &mockRuleService{}
	mockTriggerSvc := &mockTriggerService{}
	mockActionSvc := &mockActionService{}

	server := &Server{
		ruleSvc:    mockRuleSvc,
		triggerSvc: mockTriggerSvc,
		actionSvc:  mockActionSvc,
	}

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	server.HealthCheck(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "healthy", response["status"])
}

func TestServer_CreateAction(t *testing.T) {
	mockRuleSvc := &mockRuleService{}
	mockTriggerSvc := &mockTriggerService{}
	mockActionSvc := &mockActionService{}

	server := &Server{
		ruleSvc:    mockRuleSvc,
		triggerSvc: mockTriggerSvc,
		actionSvc:  mockActionSvc,
	}

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
			expectedStatus: http.StatusOK,
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

			server.CreateAction(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockActionSvc.AssertExpectations(t)
		})
	}
}
