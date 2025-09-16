package rule

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	actionStorage "github.com/malyshevhen/rule-engine/internal/storage/action"
	triggerStorage "github.com/malyshevhen/rule-engine/internal/storage/trigger"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepository_Create(t *testing.T) {
	pool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer pool.Close()

	repo := NewRepository(pool)

	rule := &Rule{
		Name:      "Test Rule",
		LuaScript: "print('hello')",
		Enabled:   true,
	}

	expectedID := uuid.New()
	expectedCreatedAt := time.Now()
	expectedUpdatedAt := time.Now()

	pool.ExpectQuery(`INSERT INTO rules \(name, lua_script, enabled\) VALUES \(\$1, \$2, \$3\) RETURNING id, created_at, updated_at`).
		WithArgs(rule.Name, rule.LuaScript, rule.Enabled).
		WillReturnRows(pgxmock.NewRows([]string{"id", "created_at", "updated_at"}).
			AddRow(expectedID, expectedCreatedAt, expectedUpdatedAt))

	err = repo.Create(context.Background(), rule)

	assert.NoError(t, err)
	assert.Equal(t, expectedID, rule.ID)
	assert.Equal(t, expectedCreatedAt, rule.CreatedAt)
	assert.Equal(t, expectedUpdatedAt, rule.UpdatedAt)

	assert.NoError(t, pool.ExpectationsWereMet())
}

func TestRepository_GetByID(t *testing.T) {
	pool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer pool.Close()

	repo := NewRepository(pool)

	expectedRule := &Rule{
		ID:        uuid.New(),
		Name:      "Test Rule",
		LuaScript: "print('hello')",
		Enabled:   true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	pool.ExpectQuery(`SELECT id, name, lua_script, enabled, created_at, updated_at FROM rules WHERE id = \$1`).
		WithArgs(expectedRule.ID).
		WillReturnRows(pgxmock.NewRows([]string{"id", "name", "lua_script", "enabled", "created_at", "updated_at"}).
			AddRow(expectedRule.ID, expectedRule.Name, expectedRule.LuaScript, expectedRule.Enabled, expectedRule.CreatedAt, expectedRule.UpdatedAt))

	rule, err := repo.GetByID(context.Background(), expectedRule.ID)

	assert.NoError(t, err)
	assert.Equal(t, expectedRule, rule)

	assert.NoError(t, pool.ExpectationsWereMet())
}

func TestRepository_GetTriggersByRuleID(t *testing.T) {
	pool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer pool.Close()

	repo := NewRepository(pool)

	ruleID := uuid.New()
	expectedTriggers := []*triggerStorage.Trigger{
		{
			ID:              uuid.New(),
			RuleID:          ruleID,
			Type:            "conditional",
			ConditionScript: "return true",
			Enabled:         true,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		},
	}

	rows := pgxmock.NewRows([]string{"id", "rule_id", "type", "condition_script", "enabled", "created_at", "updated_at"})
	for _, trigger := range expectedTriggers {
		rows.AddRow(trigger.ID, trigger.RuleID, trigger.Type, trigger.ConditionScript, trigger.Enabled, trigger.CreatedAt, trigger.UpdatedAt)
	}

	pool.ExpectQuery(`SELECT t\.id, t\.rule_id, t\.type, t\.condition_script, t\.enabled, t\.created_at, t\.updated_at FROM triggers t JOIN rule_triggers rt ON t\.id = rt\.trigger_id WHERE rt\.rule_id = \$1`).
		WithArgs(ruleID).
		WillReturnRows(rows)

	triggers, err := repo.GetTriggersByRuleID(context.Background(), ruleID)

	assert.NoError(t, err)
	assert.Len(t, triggers, len(expectedTriggers))
	assert.Equal(t, expectedTriggers, triggers)

	assert.NoError(t, pool.ExpectationsWereMet())
}

func TestRepository_GetActionsByRuleID(t *testing.T) {
	pool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer pool.Close()

	repo := NewRepository(pool)

	ruleID := uuid.New()
	expectedActions := []*actionStorage.Action{
		{
			ID:        uuid.New(),
			LuaScript: "send_command('device', 'on')",
			Enabled:   true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	rows := pgxmock.NewRows([]string{"id", "lua_script", "enabled", "created_at", "updated_at"})
	for _, action := range expectedActions {
		rows.AddRow(action.ID, action.LuaScript, action.Enabled, action.CreatedAt, action.UpdatedAt)
	}

	pool.ExpectQuery(`SELECT a\.id, a\.lua_script, a\.enabled, a\.created_at, a\.updated_at FROM actions a JOIN rule_actions ra ON a\.id = ra\.action_id WHERE ra\.rule_id = \$1`).
		WithArgs(ruleID).
		WillReturnRows(rows)

	actions, err := repo.GetActionsByRuleID(context.Background(), ruleID)

	assert.NoError(t, err)
	assert.Len(t, actions, len(expectedActions))
	assert.Equal(t, expectedActions, actions)

	assert.NoError(t, pool.ExpectationsWereMet())
}

func TestRepository_List(t *testing.T) {
	pool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer pool.Close()

	repo := NewRepository(pool)

	expectedRules := []*Rule{
		{
			ID:        uuid.New(),
			Name:      "Rule 1",
			LuaScript: "print('rule1')",
			Enabled:   true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			Name:      "Rule 2",
			LuaScript: "print('rule2')",
			Enabled:   false,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	rows := pgxmock.NewRows([]string{"id", "name", "lua_script", "enabled", "created_at", "updated_at"})
	for _, rule := range expectedRules {
		rows.AddRow(rule.ID, rule.Name, rule.LuaScript, rule.Enabled, rule.CreatedAt, rule.UpdatedAt)
	}

	pool.ExpectQuery(`SELECT id, name, lua_script, enabled, created_at, updated_at FROM rules ORDER BY created_at DESC LIMIT \$1 OFFSET \$2`).
		WithArgs(10, 0).
		WillReturnRows(rows)

	rules, err := repo.List(context.Background(), 10, 0)

	assert.NoError(t, err)
	assert.Len(t, rules, 2)
	assert.Equal(t, expectedRules[0], rules[0])
	assert.Equal(t, expectedRules[1], rules[1])

	assert.NoError(t, pool.ExpectationsWereMet())
}

func TestRepository_Update(t *testing.T) {
	pool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer pool.Close()

	repo := NewRepository(pool)

	rule := &Rule{
		ID:        uuid.New(),
		Name:      "Updated Rule",
		LuaScript: "print('updated')",
		Enabled:   false,
	}

	pool.ExpectQuery(`UPDATE rules SET name = \$1, lua_script = \$2, enabled = \$3, updated_at = NOW\(\) WHERE id = \$4`).
		WithArgs(rule.Name, rule.LuaScript, rule.Enabled, rule.ID).
		WillReturnRows(pgxmock.NewRows([]string{}))

	err = repo.Update(context.Background(), rule)

	assert.NoError(t, err)

	assert.NoError(t, pool.ExpectationsWereMet())
}

func TestRepository_Delete(t *testing.T) {
	pool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer pool.Close()

	repo := NewRepository(pool)

	id := uuid.New()

	pool.ExpectQuery(`DELETE FROM rules WHERE id = \$1`).
		WithArgs(id).
		WillReturnRows(pgxmock.NewRows([]string{}))

	err = repo.Delete(context.Background(), id)

	assert.NoError(t, err)

	assert.NoError(t, pool.ExpectationsWereMet())
}

func TestRepository_GetByID_NotFound(t *testing.T) {
	pool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer pool.Close()

	repo := NewRepository(pool)

	id := uuid.New()

	pool.ExpectQuery(`SELECT id, name, lua_script, enabled, created_at, updated_at FROM rules WHERE id = \$1`).
		WithArgs(id).
		WillReturnError(pgx.ErrNoRows)

	rule, err := repo.GetByID(context.Background(), id)

	assert.Error(t, err)
	assert.Nil(t, rule)

	assert.NoError(t, pool.ExpectationsWereMet())
}

func TestRepository_Update_Error(t *testing.T) {
	pool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer pool.Close()

	repo := NewRepository(pool)

	rule := &Rule{
		ID:        uuid.New(),
		Name:      "Updated Rule",
		LuaScript: "print('updated')",
		Enabled:   false,
	}

	pool.ExpectQuery(`UPDATE rules SET name = \$1, lua_script = \$2, enabled = \$3, updated_at = NOW\(\) WHERE id = \$4`).
		WithArgs(rule.Name, rule.LuaScript, rule.Enabled, rule.ID).
		WillReturnError(assert.AnError)

	err = repo.Update(context.Background(), rule)

	assert.Error(t, err)

	assert.NoError(t, pool.ExpectationsWereMet())
}

func TestRepository_Delete_Error(t *testing.T) {
	pool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer pool.Close()

	repo := NewRepository(pool)

	id := uuid.New()

	pool.ExpectQuery(`DELETE FROM rules WHERE id = \$1`).
		WithArgs(id).
		WillReturnError(assert.AnError)

	err = repo.Delete(context.Background(), id)

	assert.Error(t, err)

	assert.NoError(t, pool.ExpectationsWereMet())
}
