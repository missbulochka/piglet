package postgres

import (
	"errors"
	"fmt"

	// Библиотека миграций
	"github.com/golang-migrate/migrate/v4"
	// Драйвер для выполнения миграций postgres
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	// Драйвер для получения миграция из файлов
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigration(sourceURL string, databaseURL string) error {
	const op = "piglet-bills | psqlapp.RunMigration"

	migration, err := migrate.New(
		sourceURL,
		databaseURL,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err = migration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("%s: %w", op, err)
	}

	fmt.Println("migrated successfully")
	return nil
}
