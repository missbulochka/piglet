package postgres

import (
	"errors"
	"fmt"
	"log/slog"

	// Библиотека миграций
	"github.com/golang-migrate/migrate/v4"
	// Драйвер для выполнения миграций postgres
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	// Драйвер для получения миграция из файлов
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigration(log *slog.Logger, sourceURL string, databaseURL string) error {
	const op = "piglet-transactions | postgres.RunMigration"

	log = log.With(slog.String("op", op))

	migration, err := migrate.New(sourceURL, databaseURL)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err = migration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("migrated successfully")

	return nil
}
