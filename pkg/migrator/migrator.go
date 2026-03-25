// Package migrator provides database migration for Postgres.
package migrator

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // драйвер для Postgres
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// MigratePostgres performs migrations for Postgresql.
// The db URL should be: "postgres://user:pass@host:port/dbname?ssmode=disable"
func MigratePostgres(dbURL string, migrationsPath string) error {
	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsPath),
		dbURL,
	)
	if err != nil {
		return fmt.Errorf("creation of a migrator: %w", err)
	}
	defer func() { _, _ = m.Close() }()

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("error of migration: %w", err)
	}

	return nil
}
