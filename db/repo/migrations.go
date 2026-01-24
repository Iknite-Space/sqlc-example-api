package repo

import (
	"errors"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // Postgres driver
	_ "github.com/golang-migrate/migrate/v4/source/file"       // File source for migrations
)

// Migrate function applies migrations to the database.
func Migrate(dbURL string, migrationsPath string) error {
	absPath, err := filepath.Abs(migrationsPath)
	if err != nil {
		return err
	}

	// Convert Windows path to file URL format
	fileURL := "file://" + filepath.ToSlash(absPath)

	// Create a new migration instance with the file URL
	m, err := migrate.New(
		fileURL,
		dbURL,
	)

	if err != nil {
		return err
	}
	defer m.Close()

	// Apply migrations
	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

// MigrateDown function rolls back migrations from the database.
func MigrateDown(dbURL string, migrationsPath string) error {
	absPath, err := filepath.Abs(migrationsPath)
	if err != nil {
		return err
	}

	// Convert Windows path to file URL format
	fileURL := "file://" + filepath.ToSlash(absPath)

	// Create a new migration instance with the file URL
	m, err := migrate.New(
		fileURL,
		dbURL,
	)
	if err != nil {
		return err
	}
	defer m.Close()

	// Apply migrations
	err = m.Down()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
