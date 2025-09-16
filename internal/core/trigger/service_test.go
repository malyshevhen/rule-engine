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
