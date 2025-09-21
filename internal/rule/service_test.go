package rule

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/malyshevhen/rule-engine/internal/storage"
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

func (m *mockRuleRepository) List(ctx context.Context, limit int, offset int) ([]*ruleStorage.Rule, int, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, 0, args.Error(1)
	}
	return args.Get(0).([]*ruleStorage.Rule), args.Int(1), args.Error(2)
}

func (m *mockRuleRepository) ListAll(ctx context.Context) ([]*ruleStorage.Rule, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
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

// mockSQLStore is a mock implementation of Store interface for testing
type mockSQLStore struct {
	mock.Mock
	ruleRepo storage.RuleRepository
}

func newMockSQLStore() *mockSQLStore {
	return &mockSQLStore{
		ruleRepo: &mockRuleRepository{},
	}
}

func (m *mockSQLStore) ExecTx(ctx context.Context, fn func(*storage.Store) error) error {
	store := &storage.Store{
		RuleRepository: m.ruleRepo,
	}
	err := fn(store)
	m.Called(ctx, mock.Anything)
	return err
}

func (m *mockSQLStore) GetStore() *storage.Store {
	return &storage.Store{
		RuleRepository: m.ruleRepo,
	}
}

func TestService_Create(t *testing.T) {
	mockStore := newMockSQLStore()
	svc := NewService(mockStore, nil)

	rule := &Rule{
		Name:      "Test Rule",
		LuaScript: "return true",
		Enabled:   true,
	}

	mockStore.ruleRepo.(*mockRuleRepository).On("Create", mock.Anything, mock.AnythingOfType("*rule.Rule")).Return(nil)

	mockStore.On("ExecTx", mock.Anything, mock.Anything).Return(nil)

	err := svc.Create(context.Background(), rule)

	assert.NoError(t, err)
	mockStore.ruleRepo.(*mockRuleRepository).AssertExpectations(t)
	mockStore.AssertExpectations(t)
}

func TestService_GetByID(t *testing.T) {
	mockStore := newMockSQLStore()
	svc := NewService(mockStore, nil)

	ruleID := uuid.New()
	expectedRule := &ruleStorage.Rule{
		ID:        ruleID,
		Name:      "Test Rule",
		LuaScript: "return true",
		Enabled:   true,
	}

	expectedTriggers := []*triggerStorage.Trigger{}
	expectedActions := []*actionStorage.Action{}

	mockStore.ruleRepo.(*mockRuleRepository).On("GetByIDWithAssociations", mock.Anything, ruleID).Return(expectedRule, expectedTriggers, expectedActions, nil)

	rule, err := svc.GetByID(context.Background(), ruleID)

	assert.NoError(t, err)
	assert.NotNil(t, rule)
	assert.Equal(t, ruleID, rule.ID)
	assert.Equal(t, "Test Rule", rule.Name)
	assert.Equal(t, "return true", rule.LuaScript)
	assert.True(t, rule.Enabled)
	assert.Empty(t, rule.Triggers)
	assert.Empty(t, rule.Actions)

	mockStore.ruleRepo.(*mockRuleRepository).AssertExpectations(t)
}

func TestService_GetByID_Error(t *testing.T) {
	mockStore := newMockSQLStore()
	svc := NewService(mockStore, nil)

	ruleID := uuid.New()

	mockStore.ruleRepo.(*mockRuleRepository).On("GetByIDWithAssociations", mock.Anything, ruleID).Return((*ruleStorage.Rule)(nil), ([]*triggerStorage.Trigger)(nil), ([]*actionStorage.Action)(nil), assert.AnError)

	rule, err := svc.GetByID(context.Background(), ruleID)

	assert.Error(t, err)
	assert.Nil(t, rule)
	mockStore.ruleRepo.(*mockRuleRepository).AssertExpectations(t)
}

func TestService_List(t *testing.T) {
	mockStore := newMockSQLStore()
	svc := NewService(mockStore, nil)

	expectedRules := []*ruleStorage.Rule{
		{
			ID:        uuid.New(),
			Name:      "Rule 1",
			LuaScript: "return true",
			Enabled:   true,
		},
	}

	mockStore.ruleRepo.(*mockRuleRepository).On("List", mock.Anything, 1000, 0).Return(expectedRules, 1, nil)

	rules, err := svc.ListAll(context.Background())

	assert.NoError(t, err)
	assert.Len(t, rules, 1)
	assert.Equal(t, "Rule 1", rules[0].Name)

	mockStore.ruleRepo.(*mockRuleRepository).AssertExpectations(t)
}

func TestService_Update(t *testing.T) {
	mockStore := newMockSQLStore()
	svc := NewService(mockStore, nil)

	rule := &Rule{
		ID:        uuid.New(),
		Name:      "Updated Rule",
		LuaScript: "return true",
		Enabled:   true,
	}

	mockStore.ruleRepo.(*mockRuleRepository).On("Update", mock.Anything, mock.Anything).Return(nil)

	mockStore.On("ExecTx", mock.Anything, mock.Anything).Return(nil)

	err := svc.Update(context.Background(), rule)

	assert.NoError(t, err)
	mockStore.ruleRepo.(*mockRuleRepository).AssertExpectations(t)
	mockStore.AssertExpectations(t)
}

func TestService_Delete(t *testing.T) {
	mockStore := newMockSQLStore()
	svc := NewService(mockStore, nil)

	ruleID := uuid.New()

	mockStore.ruleRepo.(*mockRuleRepository).On("Delete", mock.Anything, ruleID).Return(nil)

	mockStore.On("ExecTx", mock.Anything, mock.Anything).Return(nil)

	err := svc.Delete(context.Background(), ruleID)

	assert.NoError(t, err)
	mockStore.ruleRepo.(*mockRuleRepository).AssertExpectations(t)
	mockStore.AssertExpectations(t)
}
