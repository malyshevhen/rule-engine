package action

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/malyshevhen/rule-engine/internal/storage"
	actionStorage "github.com/malyshevhen/rule-engine/internal/storage/action"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockDBTX is a mock implementation of db.DBTX
type mockDBTX struct {
	mock.Mock
}

func (m *mockDBTX) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	args := m.Called(ctx, sql, arguments)
	return args.Get(0).(pgconn.CommandTag), args.Error(1)
}

func (m *mockDBTX) Query(ctx context.Context, sql string, arguments ...any) (pgx.Rows, error) {
	args := m.Called(ctx, sql, arguments)
	return args.Get(0).(pgx.Rows), args.Error(1)
}

func (m *mockDBTX) QueryRow(ctx context.Context, sql string, arguments ...any) pgx.Row {
	args := m.Called(ctx, sql, arguments)
	return args.Get(0).(pgx.Row)
}

// mockActionRepository is a mock implementation of ActionRepository interface
type mockActionRepository struct {
	mock.Mock
}

func (m *mockActionRepository) Create(ctx context.Context, action *actionStorage.Action) error {
	args := m.Called(ctx, action)
	return args.Error(0)
}

func (m *mockActionRepository) GetByID(ctx context.Context, id uuid.UUID) (*actionStorage.Action, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*actionStorage.Action), args.Error(1)
}

func (m *mockActionRepository) List(ctx context.Context) ([]*actionStorage.Action, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*actionStorage.Action), args.Error(1)
}

// mockSQLStore is a mock implementation of Store interface for testing
type mockSQLStore struct {
	mock.Mock
	actionRepo storage.ActionRepository
}

func newMockSQLStore() *mockSQLStore {
	return &mockSQLStore{
		actionRepo: &mockActionRepository{},
	}
}

func (m *mockSQLStore) ExecTx(ctx context.Context, fn func(*storage.Store) error) error {
	store := &storage.Store{
		ActionRepository: m.actionRepo,
	}
	err := fn(store)
	m.Called(ctx, mock.AnythingOfType("func(*storage.Store) error"))
	return err
}

func (m *mockSQLStore) GetStore() *storage.Store {
	return &storage.Store{
		ActionRepository: m.actionRepo,
	}
}

func TestService_Create(t *testing.T) {
	mockStore := newMockSQLStore()
	service := NewService(mockStore)

	action := &Action{
		LuaScript: "print('action executed')",
		Enabled:   true,
	}

	mockStore.actionRepo.(*mockActionRepository).On("Create", mock.Anything, mock.MatchedBy(func(a *actionStorage.Action) bool {
		return a.Type == "lua_script" && a.Params == action.LuaScript && a.Enabled == action.Enabled
	})).Return(nil)

	mockStore.On("ExecTx", mock.Anything, mock.Anything).Return(nil)

	err := service.Create(context.Background(), action)

	assert.NoError(t, err)
	mockStore.actionRepo.(*mockActionRepository).AssertExpectations(t)
	mockStore.AssertExpectations(t)
}

func TestService_Create_Error(t *testing.T) {
	mockStore := newMockSQLStore()
	service := NewService(mockStore)

	action := &Action{
		LuaScript: "print('action executed')",
		Enabled:   true,
	}

	mockStore.actionRepo.(*mockActionRepository).On("Create", mock.Anything, mock.Anything).Return(assert.AnError)

	mockStore.On("ExecTx", mock.Anything, mock.Anything).Return(assert.AnError)

	err := service.Create(context.Background(), action)

	assert.Error(t, err)
	mockStore.actionRepo.(*mockActionRepository).AssertExpectations(t)
	mockStore.AssertExpectations(t)
}

func TestService_GetByID(t *testing.T) {
	mockStore := newMockSQLStore()
	service := NewService(mockStore)

	actionID := uuid.New()
	expectedAction := &actionStorage.Action{
		ID:      actionID,
		Type:    "lua_script",
		Params:  "print('action executed')",
		Enabled: true,
	}

	mockStore.actionRepo.(*mockActionRepository).On("GetByID", mock.Anything, actionID).Return(expectedAction, nil)

	action, err := service.GetByID(context.Background(), actionID)

	assert.NoError(t, err)
	assert.NotNil(t, action)
	assert.Equal(t, actionID, action.ID)
	assert.Equal(t, "print('action executed')", action.LuaScript)
	assert.True(t, action.Enabled)

	mockStore.actionRepo.(*mockActionRepository).AssertExpectations(t)
}

func TestService_GetByID_Error(t *testing.T) {
	mockStore := newMockSQLStore()
	service := NewService(mockStore)

	actionID := uuid.New()

	mockStore.actionRepo.(*mockActionRepository).On("GetByID", mock.Anything, actionID).Return((*actionStorage.Action)(nil), assert.AnError)

	action, err := service.GetByID(context.Background(), actionID)

	assert.Error(t, err)
	assert.Nil(t, action)
	mockStore.actionRepo.(*mockActionRepository).AssertExpectations(t)
}
