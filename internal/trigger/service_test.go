package trigger

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/malyshevhen/rule-engine/internal/storage"
	triggerStorage "github.com/malyshevhen/rule-engine/internal/storage/trigger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockTriggerRepository is a mock implementation of TriggerRepository interface
type mockTriggerRepository struct {
	mock.Mock
}

func (m *mockTriggerRepository) Create(ctx context.Context, trigger *triggerStorage.Trigger) error {
	args := m.Called(ctx, trigger)
	return args.Error(0)
}

func (m *mockTriggerRepository) GetByID(ctx context.Context, id uuid.UUID) (*triggerStorage.Trigger, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*triggerStorage.Trigger), args.Error(1)
}

func (m *mockTriggerRepository) List(ctx context.Context, limit, offset int) ([]*triggerStorage.Trigger, int, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, 0, args.Error(1)
	}
	return args.Get(0).([]*triggerStorage.Trigger), args.Int(1), args.Error(2)
}

func (m *mockTriggerRepository) Update(ctx context.Context, trigger *triggerStorage.Trigger) error {
	args := m.Called(ctx, trigger)
	return args.Error(0)
}

func (m *mockTriggerRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// mockSQLStore is a mock implementation of Store interface for testing
type mockSQLStore struct {
	mock.Mock
	triggerRepo storage.TriggerRepository
}

func newMockSQLStore() *mockSQLStore {
	return &mockSQLStore{
		triggerRepo: &mockTriggerRepository{},
	}
}

func (m *mockSQLStore) ExecTx(ctx context.Context, fn func(*storage.Store) error) error {
	store := &storage.Store{
		TriggerRepository: m.triggerRepo,
	}
	err := fn(store)
	m.Called(ctx, mock.Anything)
	return err
}

func (m *mockSQLStore) GetStore() *storage.Store {
	return &storage.Store{
		TriggerRepository: m.triggerRepo,
	}
}

func TestService_Create(t *testing.T) {
	mockStore := newMockSQLStore()
	service := NewService(mockStore, nil)

	ruleID := uuid.New()
	trigger := &Trigger{
		RuleID:          ruleID,
		Type:            TriggerType("conditional"),
		ConditionScript: "event.temperature > 25",
		Enabled:         true,
	}

	mockStore.triggerRepo.(*mockTriggerRepository).On("Create", mock.Anything, mock.MatchedBy(func(tr *triggerStorage.Trigger) bool {
		return tr.RuleID == ruleID && tr.Type == triggerStorage.TriggerType("conditional") &&
			tr.ConditionScript == "event.temperature > 25" && tr.Enabled == true
	})).Return(nil)

	mockStore.On("ExecTx", mock.Anything, mock.Anything).Return(nil)

	err := service.Create(context.Background(), trigger)

	assert.NoError(t, err)
	mockStore.triggerRepo.(*mockTriggerRepository).AssertExpectations(t)
	mockStore.AssertExpectations(t)
}

func TestService_Create_Error(t *testing.T) {
	mockStore := newMockSQLStore()
	service := NewService(mockStore, nil)

	trigger := &Trigger{
		RuleID:          uuid.New(),
		Type:            TriggerType("conditional"),
		ConditionScript: "event.temperature > 25",
		Enabled:         true,
	}

	mockStore.triggerRepo.(*mockTriggerRepository).On("Create", mock.Anything, mock.Anything).Return(assert.AnError)

	mockStore.On("ExecTx", mock.Anything, mock.Anything).Return(assert.AnError)

	err := service.Create(context.Background(), trigger)

	assert.Error(t, err)
	mockStore.triggerRepo.(*mockTriggerRepository).AssertExpectations(t)
	mockStore.AssertExpectations(t)
}

func TestService_GetByID(t *testing.T) {
	mockStore := newMockSQLStore()
	service := NewService(mockStore, nil)

	triggerID := uuid.New()
	expectedTrigger := &triggerStorage.Trigger{
		ID:              triggerID,
		RuleID:          uuid.New(),
		Type:            triggerStorage.TriggerType("conditional"),
		ConditionScript: "event.temperature > 25",
		Enabled:         true,
	}

	mockStore.triggerRepo.(*mockTriggerRepository).On("GetByID", mock.Anything, triggerID).Return(expectedTrigger, nil)

	trigger, err := service.GetByID(context.Background(), triggerID)

	assert.NoError(t, err)
	assert.NotNil(t, trigger)
	assert.Equal(t, triggerID, trigger.ID)
	assert.Equal(t, TriggerType("conditional"), trigger.Type)
	assert.Equal(t, "event.temperature > 25", trigger.ConditionScript)
	assert.True(t, trigger.Enabled)

	mockStore.triggerRepo.(*mockTriggerRepository).AssertExpectations(t)
}

func TestService_GetByID_Error(t *testing.T) {
	mockStore := newMockSQLStore()
	service := NewService(mockStore, nil)

	triggerID := uuid.New()

	mockStore.triggerRepo.(*mockTriggerRepository).On("GetByID", mock.Anything, triggerID).Return((*triggerStorage.Trigger)(nil), assert.AnError)

	trigger, err := service.GetByID(context.Background(), triggerID)

	assert.Error(t, err)
	assert.Nil(t, trigger)
	mockStore.triggerRepo.(*mockTriggerRepository).AssertExpectations(t)
}
