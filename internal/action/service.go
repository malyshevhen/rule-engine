package action

import (
	"context"

	"github.com/google/uuid"
	"github.com/malyshevhen/rule-engine/internal/storage"
	actionStorage "github.com/malyshevhen/rule-engine/internal/storage/action"
)

// ActionRepository interface for action storage operations
type ActionRepository interface {
	Create(ctx context.Context, action *actionStorage.Action) error
	GetByID(ctx context.Context, id uuid.UUID) (*actionStorage.Action, error)
	List(ctx context.Context, limit, offset int) ([]*actionStorage.Action, int, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// Store interface for database operations
type Store interface {
	ExecTx(ctx context.Context, fn func(*storage.Store) error) error
	GetStore() *storage.Store
}

// Service handles business logic for actions
type Service struct {
	store Store
}

// NewService creates a new action service
func NewService(store Store) *Service {
	return &Service{store: store}
}

// Create creates a new action
func (s *Service) Create(ctx context.Context, action *Action) error {
	return s.store.ExecTx(ctx, func(q *storage.Store) error {
		storageAction := &actionStorage.Action{
			Name:    action.Name,
			Type:    "lua_script",
			Params:  action.LuaScript,
			Enabled: action.Enabled,
		}
		err := q.ActionRepository.Create(ctx, storageAction)
		if err != nil {
			return err
		}
		// Copy the generated ID back to the domain action
		action.ID = storageAction.ID
		action.Name = action.Name // Keep the name from the input
		action.CreatedAt = storageAction.CreatedAt
		action.UpdatedAt = storageAction.UpdatedAt
		return nil
	})
}

// GetByID retrieves an action by its ID
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*Action, error) {
	storageAction, err := s.store.GetStore().ActionRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	action := &Action{
		ID:        storageAction.ID,
		Name:      "", // Name is not stored in DB yet, will be added later
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

// List retrieves actions with pagination
func (s *Service) List(ctx context.Context, limit, offset int) ([]*Action, int, error) {
	storageActions, total, err := s.store.GetStore().ActionRepository.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	actions := make([]*Action, len(storageActions))
	for i, storageAction := range storageActions {
		action := &Action{
			ID:        storageAction.ID,
			Name:      "", // Name is not stored in DB yet, will be added later
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

	return actions, total, nil
}

// Delete removes an action
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.store.ExecTx(ctx, func(q *storage.Store) error {
		return q.ActionRepository.Delete(ctx, id)
	})
}

// TODO: Add methods for action execution
