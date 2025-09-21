package storage

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	actionStorage "github.com/malyshevhen/rule-engine/internal/storage/action"
	ruleStorage "github.com/malyshevhen/rule-engine/internal/storage/rule"
	triggerStorage "github.com/malyshevhen/rule-engine/internal/storage/trigger"
)

// ActionRepository interface for action storage operations
type ActionRepository interface {
	Create(ctx context.Context, action *actionStorage.Action) error
	GetByID(ctx context.Context, id uuid.UUID) (*actionStorage.Action, error)
	List(ctx context.Context, limit, offset int) ([]*actionStorage.Action, int, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// TriggerRepository interface for trigger storage operations
type TriggerRepository interface {
	Create(ctx context.Context, trigger *triggerStorage.Trigger) error
	GetByID(ctx context.Context, id uuid.UUID) (*triggerStorage.Trigger, error)
	List(ctx context.Context, limit, offset int) ([]*triggerStorage.Trigger, int, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// RuleRepository interface for rule storage operations
type RuleRepository interface {
	Create(ctx context.Context, rule *ruleStorage.Rule) error
	GetByID(ctx context.Context, id uuid.UUID) (*ruleStorage.Rule, error)
	GetByIDWithAssociations(ctx context.Context, id uuid.UUID) (*ruleStorage.Rule, []*triggerStorage.Trigger, []*actionStorage.Action, error)
	List(ctx context.Context, limit int, offset int) ([]*ruleStorage.Rule, int, error)
	ListAll(ctx context.Context) ([]*ruleStorage.Rule, error)
	Update(ctx context.Context, rule *ruleStorage.Rule) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetTriggersByRuleID(ctx context.Context, ruleID uuid.UUID) ([]*triggerStorage.Trigger, error)
	GetActionsByRuleID(ctx context.Context, ruleID uuid.UUID) ([]*actionStorage.Action, error)
	AddAction(ctx context.Context, ruleID, actionID uuid.UUID) error
}

// Store provides all functions to execute db queries and transactions
type Store struct {
	RuleRepository    RuleRepository
	TriggerRepository TriggerRepository
	ActionRepository  ActionRepository
}

// SQLStore provides all functions to execute SQL queries and transactions
type SQLStore struct {
	*Store
	pool *pgxpool.Pool
}

// NewSQLStore creates a new SQLStore
func NewSQLStore(pool *pgxpool.Pool) *SQLStore {
	return &SQLStore{
		pool: pool,
		Store: &Store{
			RuleRepository:    ruleStorage.NewRepository(pool),
			TriggerRepository: triggerStorage.NewRepository(pool),
			ActionRepository:  actionStorage.NewRepository(pool),
		},
	}
}

// ExecTx executes a function within a database transaction
func (s *SQLStore) ExecTx(ctx context.Context, fn func(*Store) error) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}

	store := &Store{
		RuleRepository:    ruleStorage.NewRepository(tx),
		TriggerRepository: triggerStorage.NewRepository(tx),
		ActionRepository:  actionStorage.NewRepository(tx),
	}

	if err := fn(store); err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return rbErr
		}
		return err
	}

	return tx.Commit(ctx)
}

// GetStore returns the Store instance
func (s *SQLStore) GetStore() *Store {
	return s.Store
}
