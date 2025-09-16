package rule

import (
	"context"

	"github.com/google/uuid"
	"github.com/malyshevhen/rule-engine/internal/core/action"
	"github.com/malyshevhen/rule-engine/internal/core/trigger"
	actionStorage "github.com/malyshevhen/rule-engine/internal/storage/action"
	ruleStorage "github.com/malyshevhen/rule-engine/internal/storage/rule"
	triggerStorage "github.com/malyshevhen/rule-engine/internal/storage/trigger"
)

// Service handles business logic for rules
type Service struct {
	ruleRepo    *ruleStorage.Repository
	triggerRepo *triggerStorage.Repository
	actionRepo  *actionStorage.Repository
}

// NewService creates a new rule service
func NewService(ruleRepo *ruleStorage.Repository, triggerRepo *triggerStorage.Repository, actionRepo *actionStorage.Repository) *Service {
	return &Service{
		ruleRepo:    ruleRepo,
		triggerRepo: triggerRepo,
		actionRepo:  actionRepo,
	}
}

// Create creates a new rule
func (s *Service) Create(ctx context.Context, rule *Rule) error {
	storageRule := &ruleStorage.Rule{
		Name:      rule.Name,
		LuaScript: rule.LuaScript,
		Enabled:   rule.Enabled,
	}
	return s.ruleRepo.Create(ctx, storageRule)
}

// GetByID retrieves a rule with its triggers and actions
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*Rule, error) {
	// Get the rule
	ruleStorage, err := s.ruleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Load triggers and actions associated with the rule
	triggersStorage, err := s.ruleRepo.GetTriggersByRuleID(ctx, id)
	if err != nil {
		return nil, err
	}

	actionsStorage, err := s.ruleRepo.GetActionsByRuleID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Convert storage models to domain models
	triggers := make([]trigger.Trigger, len(triggersStorage))
	for i, t := range triggersStorage {
		triggers[i] = trigger.Trigger{
			ID:              t.ID,
			Type:            trigger.TriggerType(t.Type),
			ConditionScript: t.ConditionScript,
			Enabled:         t.Enabled,
			CreatedAt:       t.CreatedAt,
			UpdatedAt:       t.UpdatedAt,
		}
	}

	actions := make([]action.Action, len(actionsStorage))
	for i, a := range actionsStorage {
		actions[i] = action.Action{
			ID:        a.ID,
			LuaScript: a.LuaScript,
			Enabled:   a.Enabled,
			CreatedAt: a.CreatedAt,
			UpdatedAt: a.UpdatedAt,
		}
	}

	rule := &Rule{
		ID:        ruleStorage.ID,
		Name:      ruleStorage.Name,
		LuaScript: ruleStorage.LuaScript,
		Enabled:   ruleStorage.Enabled,
		CreatedAt: ruleStorage.CreatedAt,
		UpdatedAt: ruleStorage.UpdatedAt,
		Triggers:  triggers,
		Actions:   actions,
	}

	return rule, nil
}
