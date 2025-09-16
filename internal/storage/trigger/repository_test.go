package trigger

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
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
