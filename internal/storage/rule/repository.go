package rule

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	actionStorage "github.com/malyshevhen/rule-engine/internal/storage/action"
	"github.com/malyshevhen/rule-engine/internal/storage/db"
	triggerStorage "github.com/malyshevhen/rule-engine/internal/storage/trigger"
)

// ErrNotFound is returned when a rule is not found
var ErrNotFound = errors.New("rule not found")

// Repository handles database operations for rules
type Repository struct {
	db db.DBTX
}

// NewRepository creates a new rule repository
func NewRepository(db db.DBTX) *Repository {
	return &Repository{db: db}
}

// Create inserts a new rule into the database
func (r *Repository) Create(ctx context.Context, rule *Rule) error {
	query := `INSERT INTO rules (name, lua_script, priority, enabled) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`
	return r.db.QueryRow(ctx, query, rule.Name, rule.LuaScript, rule.Priority, rule.Enabled).Scan(&rule.ID, &rule.CreatedAt, &rule.UpdatedAt)
}

// GetByID retrieves a rule by its ID
func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*Rule, error) {
	query := `SELECT id, name, lua_script, priority, enabled, created_at, updated_at FROM rules WHERE id = $1`
	var rule Rule
	err := r.db.QueryRow(ctx, query, id).Scan(&rule.ID, &rule.Name, &rule.LuaScript, &rule.Priority, &rule.Enabled, &rule.CreatedAt, &rule.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &rule, nil
}

// GetByIDWithAssociations retrieves a rule with its triggers and actions using JOINs
func (r *Repository) GetByIDWithAssociations(ctx context.Context, id uuid.UUID) (*Rule, []*triggerStorage.Trigger, []*actionStorage.Action, error) {
	// Get the rule
	ruleQuery := `SELECT id, name, lua_script, priority, enabled, created_at, updated_at FROM rules WHERE id = $1`
	var rule Rule
	err := r.db.QueryRow(ctx, ruleQuery, id).Scan(&rule.ID, &rule.Name, &rule.LuaScript, &rule.Priority, &rule.Enabled, &rule.CreatedAt, &rule.UpdatedAt)
	if err != nil {
		return nil, nil, nil, err
	}

	// Get triggers directly
	triggersQuery := `
		SELECT t.id, t.rule_id, t.type, t.condition_script, t.enabled, t.created_at, t.updated_at
		FROM triggers t
		WHERE t.rule_id = $1
		ORDER BY t.created_at
	`
	triggersRows, err := r.db.Query(ctx, triggersQuery, id)
	if err != nil {
		return nil, nil, nil, err
	}
	defer triggersRows.Close()

	var triggers []*triggerStorage.Trigger
	for triggersRows.Next() {
		var t triggerStorage.Trigger
		err := triggersRows.Scan(&t.ID, &t.RuleID, &t.Type, &t.ConditionScript, &t.Enabled, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			return nil, nil, nil, err
		}
		triggers = append(triggers, &t)
	}

	// Get actions using JOIN
	actionsQuery := `
		SELECT a.id, a.type, a.params, a.enabled, a.created_at, a.updated_at
		FROM actions a
		INNER JOIN rule_actions ra ON a.id = ra.action_id
		WHERE ra.rule_id = $1
		ORDER BY a.created_at
	`
	actionsRows, err := r.db.Query(ctx, actionsQuery, id)
	if err != nil {
		return nil, nil, nil, err
	}
	defer actionsRows.Close()

	var actions []*actionStorage.Action
	for actionsRows.Next() {
		var a actionStorage.Action
		err := actionsRows.Scan(&a.ID, &a.Type, &a.Params, &a.Enabled, &a.CreatedAt, &a.UpdatedAt)
		if err != nil {
			return nil, nil, nil, err
		}
		actions = append(actions, &a)
	}

	return &rule, triggers, actions, nil
}

// GetTriggersByRuleID retrieves all triggers associated with a rule
func (r *Repository) GetTriggersByRuleID(ctx context.Context, ruleID uuid.UUID) ([]*triggerStorage.Trigger, error) {
	query := `
		SELECT t.id, t.rule_id, t.type, t.condition_script, t.enabled, t.created_at, t.updated_at
		FROM triggers t
		WHERE t.rule_id = $1
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
		SELECT a.id, a.type, a.params, a.enabled, a.created_at, a.updated_at
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
		err := rows.Scan(&a.ID, &a.Type, &a.Params, &a.Enabled, &a.CreatedAt, &a.UpdatedAt)
		if err != nil {
			return nil, err
		}
		actions = append(actions, &a)
	}
	return actions, nil
}

// List retrieves all rules with pagination
func (r *Repository) List(ctx context.Context, limit int, offset int) ([]*Rule, error) {
	query := `SELECT id, name, lua_script, priority, enabled, created_at, updated_at FROM rules ORDER BY priority DESC, created_at DESC LIMIT $1 OFFSET $2`
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []*Rule
	for rows.Next() {
		var rule Rule
		err := rows.Scan(&rule.ID, &rule.Name, &rule.LuaScript, &rule.Priority, &rule.Enabled, &rule.CreatedAt, &rule.UpdatedAt)
		if err != nil {
			return nil, err
		}
		rules = append(rules, &rule)
	}
	return rules, nil
}

// ListAll retrieves all rules (for backward compatibility)
func (r *Repository) ListAll(ctx context.Context) ([]*Rule, error) {
	return r.List(ctx, 1000, 0) // Default limit of 1000
}

// Update updates an existing rule
func (r *Repository) Update(ctx context.Context, rule *Rule) error {
	query := `UPDATE rules SET name = $1, lua_script = $2, priority = $3, enabled = $4, updated_at = NOW() WHERE id = $5`
	_, err := r.db.Query(ctx, query, rule.Name, rule.LuaScript, rule.Priority, rule.Enabled, rule.ID)
	return err
}

// Delete deletes a rule by ID
func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM rules WHERE id = $1 RETURNING id`
	row := r.db.QueryRow(ctx, query, id)
	var deletedID uuid.UUID
	err := row.Scan(&deletedID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}
	return nil
}

// AddAction associates an action with a rule
func (r *Repository) AddAction(ctx context.Context, ruleID, actionID uuid.UUID) error {
	query := `INSERT INTO rule_actions (rule_id, action_id) VALUES ($1, $2)`
	_, err := r.db.Exec(ctx, query, ruleID, actionID)
	return err
}
