package db

import (
	"context"
	"embed"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

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

	// Create iofs source driver for embedded migrations
	sourceDriver, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance("iofs", sourceDriver, "pgx5", driver)
	if err != nil {
		return err
	}
	defer func() {
		if merr, dberr := m.Close(); merr != nil || dberr != nil {
			slog.Error("Failed to close database migration manager", "error", merr, "dberror", dberr)

		}
	}()

	err = m.Up()
	if err != nil {
		// ErrNoChange is not an error, it just means no migrations to apply
		if err.Error() != "no change" {
			return err
		}
	}
	return nil
}
