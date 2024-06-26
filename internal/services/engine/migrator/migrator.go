package migrator

import (
	"embed"
	"fmt"
	"log/slog"

	"github.com/Boostport/migration"
	"github.com/Boostport/migration/driver/postgres"
)

const migrationsDir = "migrations"

// migrations source
//
//go:embed migrations
var embedFS embed.FS

type PostgresMigrator struct {
	driver migration.Driver
}

func NewPostgresMigrator(dsn string) (*PostgresMigrator, error) {
	dbDriver, err := postgres.New(dsn)
	if err != nil {
		return nil, fmt.Errorf("migrator driver initialization failed: %w", err)
	}

	return &PostgresMigrator{
		driver: dbDriver,
	}, nil
}

func (m *PostgresMigrator) Migrate() error {
	embedSource := &migration.EmbedMigrationSource{
		EmbedFS: embedFS,
		Dir:     migrationsDir,
	}

	// run all up migrations
	applied, err := migration.Migrate(m.driver, embedSource, migration.Up, 0)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	slog.Debug("migrations applied", "count", applied)
	return nil
}

func (m *PostgresMigrator) Close() error {
	return m.driver.Close()
}
