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

// RuleRepository interface for rule storage operations
type RuleRepository interface {
	Create(ctx context.Context, rule *ruleStorage.Rule) error
	GetByID(ctx context.Context, id uuid.UUID) (*ruleStorage.Rule, error)
	List(ctx context.Context) ([]*ruleStorage.Rule, error)
	Update(ctx context.Context, rule *ruleStorage.Rule) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetTriggersByRuleID(ctx context.Context, ruleID uuid.UUID) ([]*triggerStorage.Trigger, error)
	GetActionsByRuleID(ctx context.Context, ruleID uuid.UUID) ([]*actionStorage.Action, error)
}

// Service handles business logic for rules
type Service struct {
	ruleRepo    RuleRepository
	triggerRepo *triggerStorage.Repository
	actionRepo  *actionStorage.Repository
}

// NewService creates a new rule service
func NewService(ruleRepo RuleRepository, triggerRepo *triggerStorage.Repository, actionRepo *actionStorage.Repository) *Service {
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
	err := s.ruleRepo.Create(ctx, storageRule)
	if err != nil {
		return err
	}
	// Copy the generated ID back to the domain rule
	rule.ID = storageRule.ID
	rule.CreatedAt = storageRule.CreatedAt
	rule.UpdatedAt = storageRule.UpdatedAt
	return nil
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

// List retrieves all rules (TODO: add pagination and filtering)
func (s *Service) List(ctx context.Context) ([]*Rule, error) {
	// TODO: Implement list functionality in storage layer
	// For now, call repository but expect it to fail
	rules, err := s.ruleRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	// Convert storage models to domain models
	domainRules := make([]*Rule, len(rules))
	for i, rule := range rules {
		domainRules[i] = &Rule{
			ID:        rule.ID,
			Name:      rule.Name,
			LuaScript: rule.LuaScript,
			Enabled:   rule.Enabled,
			CreatedAt: rule.CreatedAt,
			UpdatedAt: rule.UpdatedAt,
		}
	}

	return domainRules, nil
}

// Update updates an existing rule
func (s *Service) Update(ctx context.Context, rule *Rule) error {
	storageRule := &ruleStorage.Rule{
		ID:        rule.ID,
		Name:      rule.Name,
		LuaScript: rule.LuaScript,
		Enabled:   rule.Enabled,
	}
	return s.ruleRepo.Update(ctx, storageRule)
}

// Delete deletes a rule by ID
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.ruleRepo.Delete(ctx, id)
}
