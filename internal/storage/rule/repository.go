package rule

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	actionStorage "github.com/malyshevhen/rule-engine/internal/storage/action"
	triggerStorage "github.com/malyshevhen/rule-engine/internal/storage/trigger"
)

// Pool interface for database operations
type Pool interface {
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
}

// Repository handles database operations for rules
type Repository struct {
	db Pool
}

// NewRepository creates a new rule repository
func NewRepository(db Pool) *Repository {
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

// GetTriggersByRuleID retrieves all triggers associated with a rule
func (r *Repository) GetTriggersByRuleID(ctx context.Context, ruleID uuid.UUID) ([]*triggerStorage.Trigger, error) {
	query := `
		SELECT t.id, t.rule_id, t.type, t.condition_script, t.enabled, t.created_at, t.updated_at
		FROM triggers t
		JOIN rule_triggers rt ON t.id = rt.trigger_id
		WHERE rt.rule_id = $1
	`
	rows, err := r.db.Query(ctx, query, ruleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var triggers []*triggerStorage.Trigger
	for rows.Next() {
		var t triggerStorage.Trigger
		err := rows.Scan(&t.ID, &t.RuleID, &t.Type, &t.ConditionScript, &t.Enabled, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			return nil, err
		}
		triggers = append(triggers, &t)
	}
	return triggers, nil
}

// GetActionsByRuleID retrieves all actions associated with a rule
func (r *Repository) GetActionsByRuleID(ctx context.Context, ruleID uuid.UUID) ([]*actionStorage.Action, error) {
	query := `
		SELECT a.id, a.lua_script, a.enabled, a.created_at, a.updated_at
		FROM actions a
		JOIN rule_actions ra ON a.id = ra.action_id
		WHERE ra.rule_id = $1
	`
	rows, err := r.db.Query(ctx, query, ruleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var actions []*actionStorage.Action
	for rows.Next() {
		var a actionStorage.Action
		err := rows.Scan(&a.ID, &a.LuaScript, &a.Enabled, &a.CreatedAt, &a.UpdatedAt)
		if err != nil {
			return nil, err
		}
		actions = append(actions, &a)
	}
	return actions, nil
}

// TODO: add rule repository
