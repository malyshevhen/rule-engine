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
	return s.repo.Create(ctx, storageTrigger)
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

// TODO: Add methods for trigger evaluation
