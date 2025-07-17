package migrate

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Up performs database migrations using the provided migrationsDir and databaseURL.
// It returns an error if any migration operation fails, or nil if the migrations are successful.
func Up(migrationsDir, databaseURL string) error {
	m, err := migrate.New("file://"+migrationsDir, databaseURL)
	if err != nil {
		return fmt.Errorf("failed to create a new migrate instance: %w", err)
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to perform migrations: %w", err)
	}

	return nil
}
