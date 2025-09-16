package action

import (
	"context"

	"github.com/google/uuid"
	actionStorage "github.com/malyshevhen/rule-engine/internal/storage/action"
)

// ActionRepository interface for action storage operations
type ActionRepository interface {
	Create(ctx context.Context, action *actionStorage.Action) error
	GetByID(ctx context.Context, id uuid.UUID) (*actionStorage.Action, error)
}

// Service handles business logic for actions
type Service struct {
	repo ActionRepository
}

// NewService creates a new action service
func NewService(repo ActionRepository) *Service {
	return &Service{repo: repo}
}

// Create creates a new action
func (s *Service) Create(ctx context.Context, action *Action) error {
	storageAction := &actionStorage.Action{
		LuaScript: action.LuaScript,
		Enabled:   action.Enabled,
	}
	return s.repo.Create(ctx, storageAction)
}

// TODO: Add methods for action execution
