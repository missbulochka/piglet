package postgres

import (
	"database/sql"
	"fmt"
	"log/slog"
)

type Storage struct {
	db *sql.DB
}

func New(log *slog.Logger, dataSourceName string) (*Storage, error) {
	const op = "piglet-transactions | storage.postgres.New"

	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("successfully connected to psql")

	return &Storage{db: db}, nil
}
