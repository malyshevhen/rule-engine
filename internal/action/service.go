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
	List(ctx context.Context) ([]*actionStorage.Action, error)
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
		Type:    "lua_script",
		Params:  action.LuaScript,
		Enabled: action.Enabled,
	}
	err := s.repo.Create(ctx, storageAction)
	if err != nil {
		return err
	}
	// Copy the generated ID back to the domain action
	action.ID = storageAction.ID
	action.CreatedAt = storageAction.CreatedAt
	action.UpdatedAt = storageAction.UpdatedAt
	return nil
}

// GetByID retrieves an action by its ID
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*Action, error) {
	storageAction, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	action := &Action{
		ID:        storageAction.ID,
		Type:      storageAction.Type,
		Params:    storageAction.Params,
		Enabled:   storageAction.Enabled,
		CreatedAt: storageAction.CreatedAt,
		UpdatedAt: storageAction.UpdatedAt,
	}
	// For backward compatibility
	if storageAction.Type == "lua_script" {
		action.LuaScript = storageAction.Params
	}

	return action, nil
}

// List retrieves all actions
func (s *Service) List(ctx context.Context) ([]*Action, error) {
	storageActions, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	actions := make([]*Action, len(storageActions))
	for i, storageAction := range storageActions {
		action := &Action{
			ID:        storageAction.ID,
			Type:      storageAction.Type,
			Params:    storageAction.Params,
			Enabled:   storageAction.Enabled,
			CreatedAt: storageAction.CreatedAt,
			UpdatedAt: storageAction.UpdatedAt,
		}
		// For backward compatibility
		if storageAction.Type == "lua_script" {
			action.LuaScript = storageAction.Params
		}
		actions[i] = action
	}

	return actions, nil
}

// TODO: Add methods for action execution
