package trigger

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository handles database operations for triggers
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new trigger repository
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// Create inserts a new trigger into the database
func (r *Repository) Create(ctx context.Context, trigger *Trigger) error {
	query := `INSERT INTO triggers (rule_id, type, condition_script, enabled) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`
	return r.db.QueryRow(ctx, query, trigger.RuleID, trigger.Type, trigger.ConditionScript, trigger.Enabled).Scan(&trigger.ID, &trigger.CreatedAt, &trigger.UpdatedAt)
}

// GetByID retrieves a trigger by its ID
func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*Trigger, error) {
	query := `SELECT id, rule_id, type, condition_script, enabled, created_at, updated_at FROM triggers WHERE id = $1`
	var trigger Trigger
	err := r.db.QueryRow(ctx, query, id).Scan(&trigger.ID, &trigger.RuleID, &trigger.Type, &trigger.ConditionScript, &trigger.Enabled, &trigger.CreatedAt, &trigger.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &trigger, nil
}

// TODO: add trigger repository
