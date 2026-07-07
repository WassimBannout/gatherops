package database

import (
	"errors"
	"fmt"
	"net/url"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type MigrationDirection string

const (
	MigrationUp   MigrationDirection = "up"
	MigrationDown MigrationDirection = "down"
)

func MigrationSourceURL(path string) (string, error) {
	if path == "" {
		return "", errors.New("migration path is required")
	}

	if filepath.IsAbs(path) {
		return (&url.URL{Scheme: "file", Path: filepath.ToSlash(path)}).String(), nil
	}

	return "file://" + filepath.ToSlash(path), nil
}

func RunMigrations(databaseURL, migrationsPath string, direction MigrationDirection, steps int) error {
	if databaseURL == "" {
		return errors.New("database URL is required")
	}

	sourceURL, err := MigrationSourceURL(migrationsPath)
	if err != nil {
		return err
	}

	migrator, err := migrate.New(sourceURL, databaseURL)
	if err != nil {
		return fmt.Errorf("create migrator: %w", err)
	}
	defer migrator.Close()

	switch direction {
	case MigrationUp:
		if steps > 0 {
			err = migrator.Steps(steps)
		} else {
			err = migrator.Up()
		}
	case MigrationDown:
		if steps > 0 {
			err = migrator.Steps(-steps)
		} else {
			err = migrator.Down()
		}
	default:
		return fmt.Errorf("unsupported migration direction %q", direction)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("run %s migrations: %w", direction, err)
	}

	return nil
}
