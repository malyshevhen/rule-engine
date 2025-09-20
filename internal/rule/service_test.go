package rule

import (
	"context"
	"testing"

	"github.com/google/uuid"
	actionStorage "github.com/malyshevhen/rule-engine/internal/storage/action"
	ruleStorage "github.com/malyshevhen/rule-engine/internal/storage/rule"
	triggerStorage "github.com/malyshevhen/rule-engine/internal/storage/trigger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockRuleRepository is a mock implementation of RuleRepository
type mockRuleRepository struct {
	mock.Mock
}

func (m *mockRuleRepository) Create(ctx context.Context, rule *ruleStorage.Rule) error {
	args := m.Called(ctx, rule)
	return args.Error(0)
}

func (m *mockRuleRepository) GetByID(ctx context.Context, id uuid.UUID) (*ruleStorage.Rule, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*ruleStorage.Rule), args.Error(1)
}

func (m *mockRuleRepository) GetByIDWithAssociations(ctx context.Context, id uuid.UUID) (*ruleStorage.Rule, []*triggerStorage.Trigger, []*actionStorage.Action, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*ruleStorage.Rule), args.Get(1).([]*triggerStorage.Trigger), args.Get(2).([]*actionStorage.Action), args.Error(3)
}

func (m *mockRuleRepository) GetTriggersByRuleID(ctx context.Context, ruleID uuid.UUID) ([]*triggerStorage.Trigger, error) {
	args := m.Called(ctx, ruleID)
	return args.Get(0).([]*triggerStorage.Trigger), args.Error(1)
}

func (m *mockRuleRepository) GetActionsByRuleID(ctx context.Context, ruleID uuid.UUID) ([]*actionStorage.Action, error) {
	args := m.Called(ctx, ruleID)
	return args.Get(0).([]*actionStorage.Action), args.Error(1)
}

func (m *mockRuleRepository) AddAction(ctx context.Context, ruleID, actionID uuid.UUID) error {
	args := m.Called(ctx, ruleID, actionID)
	return args.Error(0)
}

func (m *mockRuleRepository) List(ctx context.Context, limit int, offset int) ([]*ruleStorage.Rule, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]*ruleStorage.Rule), args.Error(1)
}

func (m *mockRuleRepository) ListAll(ctx context.Context) ([]*ruleStorage.Rule, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*ruleStorage.Rule), args.Error(1)
}

func (m *mockRuleRepository) Update(ctx context.Context, rule *ruleStorage.Rule) error {
	args := m.Called(ctx, rule)
	return args.Error(0)
}

func (m *mockRuleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestService_Create(t *testing.T) {
	mockRepo := &mockRuleRepository{}
	svc := NewService(mockRepo, nil, nil, nil)

	rule := &Rule{
		Name:      "Test Rule",
		LuaScript: "return true",
		Enabled:   true,
	}

	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*rule.Rule")).Return(nil)

	err := svc.Create(context.Background(), rule)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestService_GetByID(t *testing.T) {
	mockRepo := &mockRuleRepository{}
	svc := NewService(mockRepo, nil, nil, nil)

	ruleID := uuid.New()
	expectedRule := &ruleStorage.Rule{
		ID:        ruleID,
		Name:      "Test Rule",
		LuaScript: "return true",
		Enabled:   true,
	}

	expectedTriggers := []*triggerStorage.Trigger{}
	expectedActions := []*actionStorage.Action{}

	mockRepo.On("GetByIDWithAssociations", mock.Anything, ruleID).Return(expectedRule, expectedTriggers, expectedActions, nil)

	rule, err := svc.GetByID(context.Background(), ruleID)

	assert.NoError(t, err)
	assert.NotNil(t, rule)
	assert.Equal(t, ruleID, rule.ID)
	assert.Equal(t, "Test Rule", rule.Name)
	assert.Equal(t, "return true", rule.LuaScript)
	assert.True(t, rule.Enabled)
	assert.Empty(t, rule.Triggers)
	assert.Empty(t, rule.Actions)

	mockRepo.AssertExpectations(t)
}

func TestService_GetByID_Error(t *testing.T) {
	mockRepo := &mockRuleRepository{}
	svc := NewService(mockRepo, nil, nil, nil)

	ruleID := uuid.New()

	mockRepo.On("GetByIDWithAssociations", mock.Anything, ruleID).Return((*ruleStorage.Rule)(nil), ([]*triggerStorage.Trigger)(nil), ([]*actionStorage.Action)(nil), assert.AnError)

	rule, err := svc.GetByID(context.Background(), ruleID)

	assert.Error(t, err)
	assert.Nil(t, rule)
	mockRepo.AssertExpectations(t)
}

func TestService_List(t *testing.T) {
	mockRepo := &mockRuleRepository{}
	svc := NewService(mockRepo, nil, nil, nil)

	expectedRules := []*ruleStorage.Rule{
		{
			ID:        uuid.New(),
			Name:      "Rule 1",
			LuaScript: "return true",
			Enabled:   true,
		},
	}

	mockRepo.On("List", mock.Anything, 1000, 0).Return(expectedRules, nil)

	rules, err := svc.ListAll(context.Background())

	assert.NoError(t, err)
	assert.Len(t, rules, 1)
	assert.Equal(t, "Rule 1", rules[0].Name)

	mockRepo.AssertExpectations(t)
}

func TestService_Update(t *testing.T) {
	mockRepo := &mockRuleRepository{}
	svc := NewService(mockRepo, nil, nil, nil)

	rule := &Rule{
		ID:        uuid.New(),
		Name:      "Updated Rule",
		LuaScript: "return true",
		Enabled:   true,
	}

	mockRepo.On("Update", mock.Anything, mock.Anything).Return(nil)

	err := svc.Update(context.Background(), rule)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestService_Delete(t *testing.T) {
	mockRepo := &mockRuleRepository{}
	svc := NewService(mockRepo, nil, nil, nil)

	ruleID := uuid.New()

	mockRepo.On("Delete", mock.Anything, ruleID).Return(nil)

	err := svc.Delete(context.Background(), ruleID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
