package action

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// Pool interface for database operations
type Pool interface {
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
}

// Repository handles database operations for actions
type Repository struct {
	db Pool
}

// NewRepository creates a new action repository
func NewRepository(db Pool) *Repository {
	return &Repository{db: db}
}

// Create inserts a new action into the database
func (r *Repository) Create(ctx context.Context, action *Action) error {
	query := `INSERT INTO actions (lua_script, enabled) VALUES ($1, $2) RETURNING id, created_at, updated_at`
	return r.db.QueryRow(ctx, query, action.LuaScript, action.Enabled).Scan(&action.ID, &action.CreatedAt, &action.UpdatedAt)
}

// GetByID retrieves an action by its ID
func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*Action, error) {
	query := `SELECT id, lua_script, enabled, created_at, updated_at FROM actions WHERE id = $1`
	var action Action
	err := r.db.QueryRow(ctx, query, id).Scan(&action.ID, &action.LuaScript, &action.Enabled, &action.CreatedAt, &action.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &action, nil
}

// List retrieves all actions
func (r *Repository) List(ctx context.Context) ([]*Action, error) {
	query := `SELECT id, lua_script, enabled, created_at, updated_at FROM actions ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var actions []*Action
	for rows.Next() {
		var action Action
		err := rows.Scan(&action.ID, &action.LuaScript, &action.Enabled, &action.CreatedAt, &action.UpdatedAt)
		if err != nil {
			return nil, err
		}
		actions = append(actions, &action)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return actions, nil
}

// TODO: add action repository
