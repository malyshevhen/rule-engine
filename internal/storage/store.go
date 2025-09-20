package storage

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	actionStorage "github.com/malyshevhen/rule-engine/internal/storage/action"
	ruleStorage "github.com/malyshevhen/rule-engine/internal/storage/rule"
	triggerStorage "github.com/malyshevhen/rule-engine/internal/storage/trigger"
)

// Store provides all functions to execute db queries and transactions
type Store struct {
	RuleRepository    *ruleStorage.Repository
	TriggerRepository *triggerStorage.Repository
	ActionRepository  *actionStorage.Repository
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
