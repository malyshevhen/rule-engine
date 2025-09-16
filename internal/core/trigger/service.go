package trigger

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	redisClient "github.com/malyshevhen/rule-engine/internal/storage/redis"
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
	repo  TriggerRepository
	redis *redisClient.Client
}

// NewService creates a new trigger service
func NewService(repo TriggerRepository, redis *redisClient.Client) *Service {
	return &Service{repo: repo, redis: redis}
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

	// Invalidate caches
	s.invalidateTriggerCaches(ctx)

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

// GetEnabledConditionalTriggers retrieves all enabled conditional triggers
func (s *Service) GetEnabledConditionalTriggers(ctx context.Context) ([]*Trigger, error) {
	cacheKey := "triggers:enabled_conditional"

	// Try to get from cache first
	if s.redis != nil {
		if cached, err := s.redis.Get(ctx, cacheKey); err == nil {
			var triggers []*Trigger
			if err := json.Unmarshal([]byte(cached), &triggers); err == nil {
				return triggers, nil
			}
		}
	}

	allTriggers, err := s.List(ctx)
	if err != nil {
		return nil, err
	}

	var conditionalTriggers []*Trigger
	for _, trigger := range allTriggers {
		if trigger.Type == Conditional && trigger.Enabled {
			conditionalTriggers = append(conditionalTriggers, trigger)
		}
	}

	// Cache the result
	if s.redis != nil {
		if data, err := json.Marshal(conditionalTriggers); err == nil {
			s.redis.Set(ctx, cacheKey, string(data), 1*time.Minute)
		}
	}

	return conditionalTriggers, nil
}

// GetEnabledScheduledTriggers retrieves all enabled scheduled (CRON) triggers
func (s *Service) GetEnabledScheduledTriggers(ctx context.Context) ([]*Trigger, error) {
	cacheKey := "triggers:enabled_scheduled"

	// Try to get from cache first
	if s.redis != nil {
		if cached, err := s.redis.Get(ctx, cacheKey); err == nil {
			var triggers []*Trigger
			if err := json.Unmarshal([]byte(cached), &triggers); err == nil {
				return triggers, nil
			}
		}
	}

	allTriggers, err := s.List(ctx)
	if err != nil {
		return nil, err
	}

	var scheduledTriggers []*Trigger
	for _, trigger := range allTriggers {
		if trigger.Type == Cron && trigger.Enabled {
			scheduledTriggers = append(scheduledTriggers, trigger)
		}
	}

	// Cache the result
	if s.redis != nil {
		if data, err := json.Marshal(scheduledTriggers); err == nil {
			s.redis.Set(ctx, cacheKey, string(data), 1*time.Minute)
		}
	}

	return scheduledTriggers, nil
}

// invalidateTriggerCaches clears all trigger-related caches
func (s *Service) invalidateTriggerCaches(ctx context.Context) {
	if s.redis == nil {
		return
	}

	// Delete trigger caches
	keys, err := s.redis.Keys(ctx, "triggers:*")
	if err == nil {
		for _, key := range keys {
			s.redis.Del(ctx, key)
		}
	}
}
