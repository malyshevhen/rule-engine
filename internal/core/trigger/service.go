package trigger

import (
	"context"

	"github.com/google/uuid"
	triggerStorage "github.com/malyshevhen/rule-engine/internal/storage/trigger"
)

// TriggerRepository interface for trigger storage operations
type TriggerRepository interface {
	Create(ctx context.Context, trigger *triggerStorage.Trigger) error
	GetByID(ctx context.Context, id uuid.UUID) (*triggerStorage.Trigger, error)
	List(ctx context.Context) ([]*triggerStorage.Trigger, error)
}

// Service handles business logic for triggers
type Service struct {
	repo TriggerRepository
}

// NewService creates a new trigger service
func NewService(repo TriggerRepository) *Service {
	return &Service{repo: repo}
}

// Create creates a new trigger
func (s *Service) Create(ctx context.Context, trigger *Trigger) error {
	storageTrigger := &triggerStorage.Trigger{
		RuleID:          trigger.RuleID,
		Type:            triggerStorage.TriggerType(trigger.Type),
		ConditionScript: trigger.ConditionScript,
		Enabled:         trigger.Enabled,
	}
	err := s.repo.Create(ctx, storageTrigger)
	if err != nil {
		return err
	}
	// Copy the generated ID back to the domain trigger
	trigger.ID = storageTrigger.ID
	trigger.CreatedAt = storageTrigger.CreatedAt
	trigger.UpdatedAt = storageTrigger.UpdatedAt
	return nil
}

// GetByID retrieves a trigger by its ID
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*Trigger, error) {
	storageTrigger, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	trigger := &Trigger{
		ID:              storageTrigger.ID,
		RuleID:          storageTrigger.RuleID,
		Type:            TriggerType(storageTrigger.Type),
		ConditionScript: storageTrigger.ConditionScript,
		Enabled:         storageTrigger.Enabled,
		CreatedAt:       storageTrigger.CreatedAt,
		UpdatedAt:       storageTrigger.UpdatedAt,
	}

	return trigger, nil
}

// List retrieves all triggers
func (s *Service) List(ctx context.Context) ([]*Trigger, error) {
	storageTriggers, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	triggers := make([]*Trigger, len(storageTriggers))
	for i, storageTrigger := range storageTriggers {
		triggers[i] = &Trigger{
			ID:              storageTrigger.ID,
			RuleID:          storageTrigger.RuleID,
			Type:            TriggerType(storageTrigger.Type),
			ConditionScript: storageTrigger.ConditionScript,
			Enabled:         storageTrigger.Enabled,
			CreatedAt:       storageTrigger.CreatedAt,
			UpdatedAt:       storageTrigger.UpdatedAt,
		}
	}

	return triggers, nil
}

// TODO: Add methods for trigger evaluation
