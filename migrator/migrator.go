package migrator

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations() error {
	m, err := migrate.New(
		"file://migrations",
		"postgres://postgres:admin@localhost:5432/effective?sslmode=disable",
	)
	if err != nil {
		return fmt.Errorf("migrator.go: RunMigrations: migrate.New: %w", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migrator.go: RunMigrations: m.Up(): %w", err)
	}

	return nil
}
