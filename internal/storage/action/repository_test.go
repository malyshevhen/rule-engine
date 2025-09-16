package action

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

	action := &Action{
		Type:    "lua_script",
		Params:  "send_command('device', 'on')",
		Enabled: true,
	}

	expectedID := uuid.New()
	expectedCreatedAt := time.Now()
	expectedUpdatedAt := time.Now()

	pool.ExpectQuery(`INSERT INTO actions \(type, params, enabled\) VALUES \(\$1, \$2, \$3\) RETURNING id, created_at, updated_at`).
		WithArgs(action.Type, action.Params, action.Enabled).
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
		Type:      "lua_script",
		Params:    "send_command('device', 'on')",
		Enabled:   true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	pool.ExpectQuery(`SELECT id, type, params, enabled, created_at, updated_at FROM actions WHERE id = \$1`).
		WithArgs(expectedAction.ID).
		WillReturnRows(pgxmock.NewRows([]string{"id", "type", "params", "enabled", "created_at", "updated_at"}).
			AddRow(expectedAction.ID, expectedAction.Type, expectedAction.Params, expectedAction.Enabled, expectedAction.CreatedAt, expectedAction.UpdatedAt))

	action, err := repo.GetByID(context.Background(), expectedAction.ID)

	assert.NoError(t, err)
	assert.Equal(t, expectedAction, action)

	assert.NoError(t, pool.ExpectationsWereMet())
}

func TestRepository_List(t *testing.T) {
	pool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer pool.Close()

	repo := NewRepository(pool)

	expectedActions := []*Action{
		{
			ID:        uuid.New(),
			Type:      "lua_script",
			Params:    "send_command('device1', 'on')",
			Enabled:   true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			Type:      "lua_script",
			Params:    "send_command('device2', 'off')",
			Enabled:   false,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	rows := pgxmock.NewRows([]string{"id", "type", "params", "enabled", "created_at", "updated_at"})
	for _, action := range expectedActions {
		rows.AddRow(action.ID, action.Type, action.Params, action.Enabled, action.CreatedAt, action.UpdatedAt)
	}

	pool.ExpectQuery(`SELECT id, type, params, enabled, created_at, updated_at FROM actions ORDER BY created_at DESC`).
		WillReturnRows(rows)

	actions, err := repo.List(context.Background())

	assert.NoError(t, err)
	assert.Len(t, actions, 2)
	assert.Equal(t, expectedActions[0], actions[0])
	assert.Equal(t, expectedActions[1], actions[1])

	assert.NoError(t, pool.ExpectationsWereMet())
}

func TestRepository_List_Error(t *testing.T) {
	pool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer pool.Close()

	repo := NewRepository(pool)

	pool.ExpectQuery(`SELECT id, type, params, enabled, created_at, updated_at FROM actions ORDER BY created_at DESC`).
		WillReturnError(assert.AnError)

	actions, err := repo.List(context.Background())

	assert.Error(t, err)
	assert.Nil(t, actions)

	assert.NoError(t, pool.ExpectationsWereMet())
}

func TestRepository_GetByID_NotFound(t *testing.T) {
	pool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer pool.Close()

	repo := NewRepository(pool)

	id := uuid.New()

	pool.ExpectQuery(`SELECT id, type, params, enabled, created_at, updated_at FROM actions WHERE id = \$1`).
		WithArgs(id).
		WillReturnError(pgx.ErrNoRows)

	action, err := repo.GetByID(context.Background(), id)

	assert.Error(t, err)
	assert.Nil(t, action)

	assert.NoError(t, pool.ExpectationsWereMet())
}

func TestRepository_Create_Error(t *testing.T) {
	pool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer pool.Close()

	repo := NewRepository(pool)

	action := &Action{
		Type:    "lua_script",
		Params:  "send_command('device', 'on')",
		Enabled: true,
	}

	pool.ExpectQuery(`INSERT INTO actions \(type, params, enabled\) VALUES \(\$1, \$2, \$3\) RETURNING id, created_at, updated_at`).
		WithArgs(action.Type, action.Params, action.Enabled).
		WillReturnError(assert.AnError)

	err = repo.Create(context.Background(), action)

	assert.Error(t, err)

	assert.NoError(t, pool.ExpectationsWereMet())
}
