package rule

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository handles database operations for rules
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new rule repository
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// Create inserts a new rule into the database
func (r *Repository) Create(ctx context.Context, rule *Rule) error {
	query := `INSERT INTO rules (name, lua_script, enabled) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`
	return r.db.QueryRow(ctx, query, rule.Name, rule.LuaScript, rule.Enabled).Scan(&rule.ID, &rule.CreatedAt, &rule.UpdatedAt)
}

// GetByID retrieves a rule by its ID
func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*Rule, error) {
	query := `SELECT id, name, lua_script, enabled, created_at, updated_at FROM rules WHERE id = $1`
	var rule Rule
	err := r.db.QueryRow(ctx, query, id).Scan(&rule.ID, &rule.Name, &rule.LuaScript, &rule.Enabled, &rule.CreatedAt, &rule.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &rule, nil
}

// TODO: add rule repository
