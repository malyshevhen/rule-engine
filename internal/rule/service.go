package rule

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/malyshevhen/rule-engine/internal/action"
	actionStorage "github.com/malyshevhen/rule-engine/internal/storage/action"
	redisClient "github.com/malyshevhen/rule-engine/internal/storage/redis"
	ruleStorage "github.com/malyshevhen/rule-engine/internal/storage/rule"
	triggerStorage "github.com/malyshevhen/rule-engine/internal/storage/trigger"
	"github.com/malyshevhen/rule-engine/internal/trigger"
)

// RuleRepository interface for rule storage operations
type RuleRepository interface {
	Create(ctx context.Context, rule *ruleStorage.Rule) error
	GetByID(ctx context.Context, id uuid.UUID) (*ruleStorage.Rule, error)
	GetByIDWithAssociations(ctx context.Context, id uuid.UUID) (*ruleStorage.Rule, []*triggerStorage.Trigger, []*actionStorage.Action, error)
	List(ctx context.Context, limit int, offset int) ([]*ruleStorage.Rule, error)
	ListAll(ctx context.Context) ([]*ruleStorage.Rule, error)
	Update(ctx context.Context, rule *ruleStorage.Rule) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetTriggersByRuleID(ctx context.Context, ruleID uuid.UUID) ([]*triggerStorage.Trigger, error)
	GetActionsByRuleID(ctx context.Context, ruleID uuid.UUID) ([]*actionStorage.Action, error)
	AddAction(ctx context.Context, ruleID, actionID uuid.UUID) error
}

// Service handles business logic for rules
type Service struct {
	ruleRepo    RuleRepository
	triggerRepo *triggerStorage.Repository
	actionRepo  *actionStorage.Repository
	redis       *redisClient.Client
}

// NewService creates a new rule service
func NewService(ruleRepo RuleRepository, triggerRepo *triggerStorage.Repository, actionRepo *actionStorage.Repository, redis *redisClient.Client) *Service {
	return &Service{
		ruleRepo:    ruleRepo,
		triggerRepo: triggerRepo,
		actionRepo:  actionRepo,
		redis:       redis,
	}
}

// Create creates a new rule
func (s *Service) Create(ctx context.Context, rule *Rule) error {
	storageRule := &ruleStorage.Rule{
		Name:      rule.Name,
		LuaScript: rule.LuaScript,
		Priority:  rule.Priority,
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

	// Invalidate caches
	s.invalidateRuleCaches(ctx, rule.ID)

	return nil
}

// GetByID retrieves a rule with its triggers and actions
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*Rule, error) {
	cacheKey := fmt.Sprintf("rule:%s", id.String())

	// Try to get from cache first
	if s.redis != nil {
		if cached, err := s.redis.Get(ctx, cacheKey); err == nil {
			var rule Rule
			if err := json.Unmarshal([]byte(cached), &rule); err == nil {
				return &rule, nil
			}
		}
	}

	// Get the rule with associations using optimized JOIN queries
	ruleStorage, triggersStorage, actionsStorage, err := s.ruleRepo.GetByIDWithAssociations(ctx, id)
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
			Type:      a.Type,
			Params:    a.Params,
			Enabled:   a.Enabled,
			CreatedAt: a.CreatedAt,
			UpdatedAt: a.UpdatedAt,
		}
		// For backward compatibility
		if a.Type == "lua_script" {
			actions[i].LuaScript = a.Params
		}
	}

	rule := &Rule{
		ID:        ruleStorage.ID,
		Name:      ruleStorage.Name,
		LuaScript: ruleStorage.LuaScript,
		Priority:  ruleStorage.Priority,
		Enabled:   ruleStorage.Enabled,
		CreatedAt: ruleStorage.CreatedAt,
		UpdatedAt: ruleStorage.UpdatedAt,
		Triggers:  triggers,
		Actions:   actions,
	}

	// Cache the result
	if s.redis != nil {
		if data, err := json.Marshal(rule); err == nil {
			if err := s.redis.Set(ctx, cacheKey, string(data), 5*time.Minute); err != nil {
				slog.Warn("Failed to cache rule", "error", err)
			}
		}
	}

	return rule, nil
}

// List retrieves rules with pagination
func (s *Service) List(ctx context.Context, limit int, offset int) ([]*Rule, error) {
	// Get current cache version for proper invalidation
	version := int64(0)
	if s.redis != nil {
		if v, err := s.redis.Get(ctx, "rules_list_version"); err == nil {
			if parsed, err := strconv.ParseInt(v, 10, 64); err == nil {
				version = parsed
			}
		}
	}
	cacheKey := fmt.Sprintf("rules:list:v%d:%d:%d", version, limit, offset)

	// Try to get from cache first
	if s.redis != nil {
		if cached, err := s.redis.Get(ctx, cacheKey); err == nil {
			var rules []*Rule
			if err := json.Unmarshal([]byte(cached), &rules); err == nil {
				return rules, nil
			}
		}
	}

	rules, err := s.ruleRepo.List(ctx, limit, offset)
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
			Priority:  rule.Priority,
			Enabled:   rule.Enabled,
			CreatedAt: rule.CreatedAt,
			UpdatedAt: rule.UpdatedAt,
		}
	}

	// Cache the result
	if s.redis != nil {
		if data, err := json.Marshal(domainRules); err == nil {
			if err := s.redis.Set(ctx, cacheKey, string(data), 2*time.Minute); err != nil {
				slog.Warn("Failed to cache rules list", "error", err)
			}
		}
	}

	return domainRules, nil
}

// ListAll retrieves all rules (for backward compatibility)
func (s *Service) ListAll(ctx context.Context) ([]*Rule, error) {
	return s.List(ctx, 1000, 0) // Default limit of 1000
}

// Update updates an existing rule
func (s *Service) Update(ctx context.Context, rule *Rule) error {
	storageRule := &ruleStorage.Rule{
		ID:        rule.ID,
		Name:      rule.Name,
		LuaScript: rule.LuaScript,
		Priority:  rule.Priority,
		Enabled:   rule.Enabled,
	}
	err := s.ruleRepo.Update(ctx, storageRule)
	if err != nil {
		return err
	}

	// Invalidate caches
	s.invalidateRuleCaches(ctx, rule.ID)

	return nil
}

// Delete deletes a rule by ID
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.ruleRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	// Invalidate caches for this specific rule
	s.invalidateRuleCaches(ctx, id)

	return nil
}

// AddAction adds an action to a rule
func (s *Service) AddAction(ctx context.Context, ruleID, actionID uuid.UUID) error {
	if err := s.ruleRepo.AddAction(ctx, ruleID, actionID); err != nil {
		return err
	}

	// Invalidate caches
	s.invalidateRuleCaches(ctx, ruleID)

	return nil
}

// invalidateRuleCaches clears rule-related caches for a specific rule
func (s *Service) invalidateRuleCaches(ctx context.Context, ruleID uuid.UUID) {
	if s.redis == nil {
		return
	}

	// Delete specific rule cache
	ruleKey := fmt.Sprintf("rule:%s", ruleID.String())
	if err := s.redis.Del(ctx, ruleKey); err != nil {
		slog.Warn("Failed to delete rule cache", "rule_id", ruleID, "error", err)
	}

	// Invalidate list caches by incrementing the version
	// This will make all existing list cache keys stale
	_, err := s.redis.Incr(ctx, "rules_list_version")
	if err != nil {
		// Log error but don't fail the operation
		slog.Warn("Failed to increment cache version", "error", err)
	}
}
