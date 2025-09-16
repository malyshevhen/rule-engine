package rule

import (
	"context"

	"github.com/google/uuid"
	"github.com/malyshevhen/rule-engine/internal/core/action"
	"github.com/malyshevhen/rule-engine/internal/core/trigger"
	actionStorage "github.com/malyshevhen/rule-engine/internal/storage/action"
	"github.com/malyshevhen/rule-engine/internal/storage/rule"
	triggerStorage "github.com/malyshevhen/rule-engine/internal/storage/trigger"
)

// Service handles business logic for rules
type Service struct {
	ruleRepo    *rule.Repository
	triggerRepo *triggerStorage.Repository
	actionRepo  *actionStorage.Repository
}

// NewService creates a new rule service
func NewService(ruleRepo *rule.Repository, triggerRepo *triggerStorage.Repository, actionRepo *actionStorage.Repository) *Service {
	return &Service{
		ruleRepo:    ruleRepo,
		triggerRepo: triggerRepo,
		actionRepo:  actionRepo,
	}
}

// GetByID retrieves a rule with its triggers and actions
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*Rule, error) {
	// Get the rule
	ruleStorage, err := s.ruleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// TODO: Load triggers and actions associated with the rule
	// For now, return rule without them
	rule := &Rule{
		ID:        ruleStorage.ID,
		Name:      ruleStorage.Name,
		LuaScript: ruleStorage.LuaScript,
		Enabled:   ruleStorage.Enabled,
		CreatedAt: ruleStorage.CreatedAt,
		UpdatedAt: ruleStorage.UpdatedAt,
		Triggers:  []trigger.Trigger{}, // TODO
		Actions:   []action.Action{},   // TODO
	}

	return rule, nil
}
