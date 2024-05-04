package psqlapp

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log/slog"
)

type MigrationApp struct {
	log           *slog.Logger
	migrationPath string
	DBHost        string
	DBPort        string
	DBName        string
}

func New(
	log *slog.Logger,
	migrationPath string,
	DBPort string,
	DBName string,
) *MigrationApp {
	return &MigrationApp{
		log:           log,
		migrationPath: migrationPath,
		DBHost:        "bills_psql",
		DBPort:        DBPort,
		DBName:        DBName,
	}
}

func (m *MigrationApp) MustRunMigrations() {
	if err := m.runMigration(); err != nil {
		panic(err)
	}
}

func (m *MigrationApp) runMigration() error {
	const op = "piglet-bills | psql.RunMigration"

	m.log = m.log.With(
		slog.String("op", op),
		slog.String("migration path", m.migrationPath),
		slog.String("database name", m.DBName),
	)

	migration, err := migrate.New(
		"file://"+m.migrationPath,
		fmt.Sprintf("postgres://%s:%s/%s?sslmode=enable", m.DBHost, m.DBPort, m.DBName),
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err = migration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("%s: %w", op, err)
	}

	m.log.Info("migrated successfully")
	return nil
}
