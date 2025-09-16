package action

import (
	"context"

	actionStorage "github.com/malyshevhen/rule-engine/internal/storage/action"
)

// Service handles business logic for actions
type Service struct {
	repo *actionStorage.Repository
}

// NewService creates a new action service
func NewService(repo *actionStorage.Repository) *Service {
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
