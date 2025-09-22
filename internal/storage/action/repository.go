package action

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/malyshevhen/rule-engine/internal/storage/db"
)

// ErrNotFound is returned when an action is not found
var ErrNotFound = errors.New("action not found")

// Repository handles database operations for actions
type Repository struct {
	db db.DBTX
}

// NewRepository creates a new action repository
func NewRepository(db db.DBTX) *Repository {
	return &Repository{db: db}
}

// Create inserts a new action into the database
func (r *Repository) Create(ctx context.Context, action *Action) error {
	query := `INSERT INTO actions (name, type, params, enabled) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`
	return r.db.QueryRow(ctx, query, action.Name, action.Type, action.Params, action.Enabled).Scan(&action.ID, &action.CreatedAt, &action.UpdatedAt)
}

// GetByID retrieves an action by its ID
func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*Action, error) {
	query := `SELECT id, name, type, params, enabled, created_at, updated_at FROM actions WHERE id = $1`
	var action Action
	err := r.db.QueryRow(ctx, query, id).Scan(&action.ID, &action.Name, &action.Type, &action.Params, &action.Enabled, &action.CreatedAt, &action.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &action, nil
}

// List retrieves actions with pagination
func (r *Repository) List(ctx context.Context, limit, offset int) ([]*Action, int, error) {
	// First get the total count
	countQuery := `SELECT COUNT(*) FROM actions`
	var total int
	err := r.db.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Then get the paginated results
	query := `SELECT id, name, type, params, enabled, created_at, updated_at FROM actions ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var actions []*Action
	for rows.Next() {
		var action Action
		err := rows.Scan(&action.ID, &action.Name, &action.Type, &action.Params, &action.Enabled, &action.CreatedAt, &action.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		actions = append(actions, &action)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return actions, total, nil
}

// Update modifies an existing action in the database
func (r *Repository) Update(ctx context.Context, action *Action) error {
	query := `UPDATE actions SET name = $1, type = $2, params = $3, enabled = $4, updated_at = NOW() WHERE id = $5`
	result, err := r.db.Exec(ctx, query, action.Name, action.Type, action.Params, action.Enabled, action.ID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

// Delete removes an action from the database
func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM actions WHERE id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}
