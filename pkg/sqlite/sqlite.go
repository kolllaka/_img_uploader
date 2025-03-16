package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

var ErrNotMigrationsToApply = errors.New("no migrations to apply")

func NewDB(storagePath string) (*sql.DB, error) {
	const op = "storage.sqlite.NewDB"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return db, nil
}

func NewMigrate(storagePath string, migrationsPath string, migrationsTable string) (*migrate.Migrate, error) {
	const op = "storage.sqlite.NewMigrate"

	databaseURL := fmt.Sprintf("sqlite3://%s", storagePath)
	if migrationsTable != "" {
		databaseURL = fmt.Sprintf("%s?x-migrations-table=%s", databaseURL, migrationsTable)
	}

	migr, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsPath),
		databaseURL,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return migr, nil
}

func Up(migr *migrate.Migrate) error {
	const op = "storage.sqlite.Up"

	if err := migr.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("%s: %w", op, err)
		}

		return fmt.Errorf("%s: %w", op, ErrNotMigrationsToApply)
	}

	return nil
}

func Down(migr *migrate.Migrate) error {
	const op = "storage.sqlite.Down"

	if err := migr.Down(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("%s: %w", op, err)
		}

		return fmt.Errorf("%s: %w", op, ErrNotMigrationsToApply)
	}

	return nil
}
