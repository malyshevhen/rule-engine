package trigger

import (
	"context"
	"testing"

	"github.com/google/uuid"
	triggerStorage "github.com/malyshevhen/rule-engine/internal/storage/trigger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockTriggerRepository is a mock implementation of TriggerRepository
type mockTriggerRepository struct {
	mock.Mock
}

func (m *mockTriggerRepository) Create(ctx context.Context, trigger *triggerStorage.Trigger) error {
	args := m.Called(ctx, trigger)
	return args.Error(0)
}

func (m *mockTriggerRepository) GetByID(ctx context.Context, id uuid.UUID) (*triggerStorage.Trigger, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*triggerStorage.Trigger), args.Error(1)
}

func (m *mockTriggerRepository) List(ctx context.Context) ([]*triggerStorage.Trigger, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*triggerStorage.Trigger), args.Error(1)
}

func TestService_Create(t *testing.T) {
	mockRepo := &mockTriggerRepository{}
	service := NewService(mockRepo)

	ruleID := uuid.New()
	trigger := &Trigger{
		RuleID:          ruleID,
		Type:            TriggerType("conditional"),
		ConditionScript: "event.temperature > 25",
		Enabled:         true,
	}

	mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(tr *triggerStorage.Trigger) bool {
		return tr.RuleID == ruleID && tr.Type == triggerStorage.TriggerType("conditional") &&
			tr.ConditionScript == "event.temperature > 25" && tr.Enabled == true
	})).Return(nil)

	err := service.Create(context.Background(), trigger)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestService_Create_Error(t *testing.T) {
	mockRepo := &mockTriggerRepository{}
	service := NewService(mockRepo)

	trigger := &Trigger{
		RuleID:          uuid.New(),
		Type:            TriggerType("conditional"),
		ConditionScript: "event.temperature > 25",
		Enabled:         true,
	}

	mockRepo.On("Create", mock.Anything, mock.Anything).Return(assert.AnError)

	err := service.Create(context.Background(), trigger)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestService_GetByID(t *testing.T) {
	mockRepo := &mockTriggerRepository{}
	service := NewService(mockRepo)

	triggerID := uuid.New()
	expectedTrigger := &triggerStorage.Trigger{
		ID:              triggerID,
		RuleID:          uuid.New(),
		Type:            triggerStorage.TriggerType("conditional"),
		ConditionScript: "event.temperature > 25",
		Enabled:         true,
	}

	mockRepo.On("GetByID", mock.Anything, triggerID).Return(expectedTrigger, nil)

	trigger, err := service.GetByID(context.Background(), triggerID)

	assert.NoError(t, err)
	assert.NotNil(t, trigger)
	assert.Equal(t, triggerID, trigger.ID)
	assert.Equal(t, TriggerType("conditional"), trigger.Type)
	assert.Equal(t, "event.temperature > 25", trigger.ConditionScript)
	assert.True(t, trigger.Enabled)

	mockRepo.AssertExpectations(t)
}

func TestService_GetByID_Error(t *testing.T) {
	mockRepo := &mockTriggerRepository{}
	service := NewService(mockRepo)

	triggerID := uuid.New()

	mockRepo.On("GetByID", mock.Anything, triggerID).Return((*triggerStorage.Trigger)(nil), assert.AnError)

	trigger, err := service.GetByID(context.Background(), triggerID)

	assert.Error(t, err)
	assert.Nil(t, trigger)
	mockRepo.AssertExpectations(t)
}
