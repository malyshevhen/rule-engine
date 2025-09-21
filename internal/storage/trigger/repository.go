package trigger

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/malyshevhen/rule-engine/internal/storage/db"
)

// ErrNotFound is returned when a trigger is not found
var ErrNotFound = errors.New("trigger not found")

// Repository handles database operations for triggers
type Repository struct {
	db db.DBTX
}

// NewRepository creates a new trigger repository
func NewRepository(db db.DBTX) *Repository {
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

// List retrieves triggers with pagination
func (r *Repository) List(ctx context.Context, limit, offset int) ([]*Trigger, int, error) {
	// First get the total count
	countQuery := `SELECT COUNT(*) FROM triggers`
	var total int
	err := r.db.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Then get the paginated results
	query := `SELECT id, rule_id, type, condition_script, enabled, created_at, updated_at FROM triggers ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var triggers []*Trigger
	for rows.Next() {
		var trigger Trigger
		err := rows.Scan(&trigger.ID, &trigger.RuleID, &trigger.Type, &trigger.ConditionScript, &trigger.Enabled, &trigger.CreatedAt, &trigger.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		triggers = append(triggers, &trigger)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return triggers, total, nil
}

// Delete removes a trigger from the database
func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM triggers WHERE id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}
