package db

import (
	"context"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

// NewPostgresPool creates a new PostgreSQL connection pool
func NewPostgresPool(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	return pool, nil
}

// RunMigrations runs database migrations
func RunMigrations(pool *pgxpool.Pool) error {
	db := stdlib.OpenDBFromPool(pool)
	driver, err := pgx.WithInstance(db, &pgx.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file://internal/storage/db/migrations", "pgx5", driver)
	if err != nil {
		return err
	}
	defer m.Close()

	return m.Up()
}
