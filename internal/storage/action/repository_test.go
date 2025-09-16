package action

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

	action := &Action{
		LuaScript: "send_command('device', 'on')",
		Enabled:   true,
	}

	expectedID := uuid.New()
	expectedCreatedAt := time.Now()
	expectedUpdatedAt := time.Now()

	pool.ExpectQuery(`INSERT INTO actions \(lua_script, enabled\) VALUES \(\$1, \$2\) RETURNING id, created_at, updated_at`).
		WithArgs(action.LuaScript, action.Enabled).
		WillReturnRows(pgxmock.NewRows([]string{"id", "created_at", "updated_at"}).
			AddRow(expectedID, expectedCreatedAt, expectedUpdatedAt))

	err = repo.Create(context.Background(), action)

	assert.NoError(t, err)
	assert.Equal(t, expectedID, action.ID)
	assert.Equal(t, expectedCreatedAt, action.CreatedAt)
	assert.Equal(t, expectedUpdatedAt, action.UpdatedAt)

	assert.NoError(t, pool.ExpectationsWereMet())
}

func TestRepository_GetByID(t *testing.T) {
	pool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer pool.Close()

	repo := NewRepository(pool)

	expectedAction := &Action{
		ID:        uuid.New(),
		LuaScript: "send_command('device', 'on')",
		Enabled:   true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	pool.ExpectQuery(`SELECT id, lua_script, enabled, created_at, updated_at FROM actions WHERE id = \$1`).
		WithArgs(expectedAction.ID).
		WillReturnRows(pgxmock.NewRows([]string{"id", "lua_script", "enabled", "created_at", "updated_at"}).
			AddRow(expectedAction.ID, expectedAction.LuaScript, expectedAction.Enabled, expectedAction.CreatedAt, expectedAction.UpdatedAt))

	action, err := repo.GetByID(context.Background(), expectedAction.ID)

	assert.NoError(t, err)
	assert.Equal(t, expectedAction, action)

	assert.NoError(t, pool.ExpectationsWereMet())
}
