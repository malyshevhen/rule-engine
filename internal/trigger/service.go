package trigger

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/malyshevhen/rule-engine/internal/storage"
	redisClient "github.com/malyshevhen/rule-engine/internal/storage/redis"
	triggerStorage "github.com/malyshevhen/rule-engine/internal/storage/trigger"
)

// TriggerRepository interface for trigger storage operations
type TriggerRepository interface {
	Create(ctx context.Context, trigger *triggerStorage.Trigger) error
	GetByID(ctx context.Context, id uuid.UUID) (*triggerStorage.Trigger, error)
	List(ctx context.Context, limit, offset int) ([]*triggerStorage.Trigger, int, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// Store interface for database operations
type Store interface {
	ExecTx(ctx context.Context, fn func(*storage.Store) error) error
	GetStore() *storage.Store
}

// Service handles business logic for triggers
type Service struct {
	store Store
	redis *redisClient.Client
}

// NewService creates a new trigger service
func NewService(store Store, redis *redisClient.Client) *Service {
	return &Service{store: store, redis: redis}
}

// Create creates a new trigger
func (s *Service) Create(ctx context.Context, trigger *Trigger) error {
	return s.store.ExecTx(ctx, func(q *storage.Store) error {
		storageTrigger := &triggerStorage.Trigger{
			RuleID:          trigger.RuleID,
			Type:            triggerStorage.TriggerType(trigger.Type),
			ConditionScript: trigger.ConditionScript,
			Enabled:         trigger.Enabled,
		}
		err := q.TriggerRepository.Create(ctx, storageTrigger)
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
	})
}

// GetByID retrieves a trigger by its ID
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*Trigger, error) {
	storageTrigger, err := s.store.GetStore().TriggerRepository.GetByID(ctx, id)
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

// List retrieves triggers with pagination
func (s *Service) List(ctx context.Context, limit, offset int) ([]*Trigger, int, error) {
	storageTriggers, total, err := s.store.GetStore().TriggerRepository.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
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

	return triggers, total, nil
}

// Delete removes a trigger
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.store.ExecTx(ctx, func(q *storage.Store) error {
		err := q.TriggerRepository.Delete(ctx, id)
		if err != nil {
			return err
		}

		// Invalidate caches
		s.invalidateTriggerCaches(ctx)

		return nil
	})
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

	allTriggers, _, err := s.List(ctx, 1000, 0) // Get all triggers for filtering
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
			if err := s.redis.Set(ctx, cacheKey, string(data), 1*time.Minute); err != nil {
				slog.Warn("Failed to cache conditional triggers", "error", err)
			}
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

	allTriggers, _, err := s.List(ctx, 1000, 0) // Get all triggers for filtering
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
			if err := s.redis.Set(ctx, cacheKey, string(data), 1*time.Minute); err != nil {
				slog.Warn("Failed to cache scheduled triggers", "error", err)
			}
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
			if err := s.redis.Del(ctx, key); err != nil {
				slog.Warn("Failed to delete trigger cache", "key", key, "error", err)
			}
		}
	}
}
