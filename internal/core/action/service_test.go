package action

import (
	"context"
	"testing"

	"github.com/google/uuid"
	actionStorage "github.com/malyshevhen/rule-engine/internal/storage/action"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockActionRepository is a mock implementation of ActionRepository
type mockActionRepository struct {
	mock.Mock
}

func (m *mockActionRepository) Create(ctx context.Context, action *actionStorage.Action) error {
	args := m.Called(ctx, action)
	return args.Error(0)
}

func (m *mockActionRepository) GetByID(ctx context.Context, id uuid.UUID) (*actionStorage.Action, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*actionStorage.Action), args.Error(1)
}

func TestService_Create(t *testing.T) {
	mockRepo := &mockActionRepository{}
	service := NewService(mockRepo)

	action := &Action{
		LuaScript: "print('action executed')",
		Enabled:   true,
	}

	mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(a *actionStorage.Action) bool {
		return a.LuaScript == action.LuaScript && a.Enabled == action.Enabled
	})).Return(nil)

	err := service.Create(context.Background(), action)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestService_Create_Error(t *testing.T) {
	mockRepo := &mockActionRepository{}
	service := NewService(mockRepo)

	action := &Action{
		LuaScript: "print('action executed')",
		Enabled:   true,
	}

	mockRepo.On("Create", mock.Anything, mock.Anything).Return(assert.AnError)

	err := service.Create(context.Background(), action)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}
