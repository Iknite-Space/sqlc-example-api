package repo

import (
	"errors"
	"log"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	//"github.com/golang-migrate/migrate/v4/database"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // Postgres driver
	_ "github.com/golang-migrate/migrate/v4/source/file"       // File source for migrations
)

// Migrate function applies migrations to the database.
func Migrate(dbURL string, migrationsPath string) error {
	absPath, err := filepath.Abs(migrationsPath)
	if err != nil {
		return err
	}

	// Create a new migration instance with the absolute path
	m, err := migrate.New(
		"file://"+absPath,
		dbURL,
	)

	if err != nil {
		return err
	}
	defer func() {
    result, err := m.Close()
    if err != nil {
        // handle error, e.g. log or panic
        log.Printf("error closing: %v (result: %v)", err, result)
    }
}()


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

	// Create a new migration instance with the absolute path
	m, err := migrate.New(
		"file://"+absPath,
		dbURL,
	)
	if err != nil {
		return err
	}
	defer func() {
    result, err := m.Close()
    if err != nil {
        // handle error, e.g. log or panic
        log.Printf("error closing: %v (result: %v)", err, result)
    }
}()


	// Apply migrations
	err = m.Down()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
