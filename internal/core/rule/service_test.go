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

func (m *mockRuleRepository) GetTriggersByRuleID(ctx context.Context, ruleID uuid.UUID) ([]*triggerStorage.Trigger, error) {
	args := m.Called(ctx, ruleID)
	return args.Get(0).([]*triggerStorage.Trigger), args.Error(1)
}

func (m *mockRuleRepository) GetActionsByRuleID(ctx context.Context, ruleID uuid.UUID) ([]*actionStorage.Action, error) {
	args := m.Called(ctx, ruleID)
	return args.Get(0).([]*actionStorage.Action), args.Error(1)
}

func TestService_Create(t *testing.T) {
	mockRepo := &mockRuleRepository{}
	svc := NewService(mockRepo, nil, nil)

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
