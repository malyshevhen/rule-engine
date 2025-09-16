package trigger

import (
	"context"

	triggerStorage "github.com/malyshevhen/rule-engine/internal/storage/trigger"
)

// Service handles business logic for triggers
type Service struct {
	repo *triggerStorage.Repository
}

// NewService creates a new trigger service
func NewService(repo *triggerStorage.Repository) *Service {
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

// TODO: Add methods for trigger evaluation
