package trigger

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepository_Create(t *testing.T) {
	pool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer pool.Close()

	repo := NewRepository(pool)

	trigger := &Trigger{
		RuleID:          uuid.New(),
		Type:            "conditional",
		ConditionScript: "return true",
		Enabled:         true,
	}

	expectedID := uuid.New()
	expectedCreatedAt := time.Now()
	expectedUpdatedAt := time.Now()

	pool.ExpectQuery(`INSERT INTO triggers \(rule_id, type, condition_script, enabled\) VALUES \(\$1, \$2, \$3, \$4\) RETURNING id, created_at, updated_at`).
		WithArgs(trigger.RuleID, trigger.Type, trigger.ConditionScript, trigger.Enabled).
		WillReturnRows(pgxmock.NewRows([]string{"id", "created_at", "updated_at"}).
			AddRow(expectedID, expectedCreatedAt, expectedUpdatedAt))

	err = repo.Create(context.Background(), trigger)

	assert.NoError(t, err)
	assert.Equal(t, expectedID, trigger.ID)
	assert.Equal(t, expectedCreatedAt, trigger.CreatedAt)
	assert.Equal(t, expectedUpdatedAt, trigger.UpdatedAt)

	assert.NoError(t, pool.ExpectationsWereMet())
}

func TestRepository_GetByID(t *testing.T) {
	pool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer pool.Close()

	repo := NewRepository(pool)

	expectedTrigger := &Trigger{
		ID:              uuid.New(),
		RuleID:          uuid.New(),
		Type:            "conditional",
		ConditionScript: "return true",
		Enabled:         true,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	pool.ExpectQuery(`SELECT id, rule_id, type, condition_script, enabled, created_at, updated_at FROM triggers WHERE id = \$1`).
		WithArgs(expectedTrigger.ID).
		WillReturnRows(pgxmock.NewRows([]string{"id", "rule_id", "type", "condition_script", "enabled", "created_at", "updated_at"}).
			AddRow(expectedTrigger.ID, expectedTrigger.RuleID, expectedTrigger.Type, expectedTrigger.ConditionScript, expectedTrigger.Enabled, expectedTrigger.CreatedAt, expectedTrigger.UpdatedAt))

	trigger, err := repo.GetByID(context.Background(), expectedTrigger.ID)

	assert.NoError(t, err)
	assert.Equal(t, expectedTrigger, trigger)

	assert.NoError(t, pool.ExpectationsWereMet())
}

func TestRepository_List(t *testing.T) {
	pool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer pool.Close()

	repo := NewRepository(pool)

	expectedTriggers := []*Trigger{
		{
			ID:              uuid.New(),
			RuleID:          uuid.New(),
			Type:            "conditional",
			ConditionScript: "return true",
			Enabled:         true,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		},
		{
			ID:              uuid.New(),
			RuleID:          uuid.New(),
			Type:            "scheduled",
			ConditionScript: "@daily",
			Enabled:         false,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		},
	}

	rows := pgxmock.NewRows([]string{"id", "rule_id", "type", "condition_script", "enabled", "created_at", "updated_at"})
	for _, trigger := range expectedTriggers {
		rows.AddRow(trigger.ID, trigger.RuleID, trigger.Type, trigger.ConditionScript, trigger.Enabled, trigger.CreatedAt, trigger.UpdatedAt)
	}

	pool.ExpectQuery(`SELECT id, rule_id, type, condition_script, enabled, created_at, updated_at FROM triggers ORDER BY created_at DESC`).
		WillReturnRows(rows)

	triggers, err := repo.List(context.Background())

	assert.NoError(t, err)
	assert.Len(t, triggers, 2)
	assert.Equal(t, expectedTriggers[0], triggers[0])
	assert.Equal(t, expectedTriggers[1], triggers[1])

	assert.NoError(t, pool.ExpectationsWereMet())
}

func TestRepository_List_Error(t *testing.T) {
	pool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer pool.Close()

	repo := NewRepository(pool)

	pool.ExpectQuery(`SELECT id, rule_id, type, condition_script, enabled, created_at, updated_at FROM triggers ORDER BY created_at DESC`).
		WillReturnError(assert.AnError)

	triggers, err := repo.List(context.Background())

	assert.Error(t, err)
	assert.Nil(t, triggers)

	assert.NoError(t, pool.ExpectationsWereMet())
}

func TestRepository_GetByID_NotFound(t *testing.T) {
	pool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer pool.Close()

	repo := NewRepository(pool)

	id := uuid.New()

	pool.ExpectQuery(`SELECT id, rule_id, type, condition_script, enabled, created_at, updated_at FROM triggers WHERE id = \$1`).
		WithArgs(id).
		WillReturnError(pgx.ErrNoRows)

	trigger, err := repo.GetByID(context.Background(), id)

	assert.Error(t, err)
	assert.Nil(t, trigger)

	assert.NoError(t, pool.ExpectationsWereMet())
}

func TestRepository_Create_Error(t *testing.T) {
	pool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer pool.Close()

	repo := NewRepository(pool)

	trigger := &Trigger{
		RuleID:          uuid.New(),
		Type:            "conditional",
		ConditionScript: "return true",
		Enabled:         true,
	}

	pool.ExpectQuery(`INSERT INTO triggers \(rule_id, type, condition_script, enabled\) VALUES \(\$1, \$2, \$3, \$4\) RETURNING id, created_at, updated_at`).
		WithArgs(trigger.RuleID, trigger.Type, trigger.ConditionScript, trigger.Enabled).
		WillReturnError(assert.AnError)

	err = repo.Create(context.Background(), trigger)

	assert.Error(t, err)

	assert.NoError(t, pool.ExpectationsWereMet())
}
