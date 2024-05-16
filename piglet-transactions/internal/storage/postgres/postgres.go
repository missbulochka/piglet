package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"piglet-transactions-service/internal/storage"
)

type Storage struct {
	db            *sql.DB
	billsMutex    sync.Mutex
	catMutex      sync.Mutex
	transMutex    sync.Mutex
	incMutex      sync.Mutex
	expMutex      sync.Mutex
	debtMutex     sync.Mutex
	transferMutex sync.Mutex
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

func (s *Storage) DefaultTransInfo(ctx context.Context, id uuid.UUID) (
	date time.Time,
	transType uint8,
	sum decimal.Decimal,
	comment string,
	err error) {
	const op = "piglet-transactions | storage.postgres.defaultTransInfo"

	row := s.db.QueryRowContext(
		ctx,
		storage.GetOneTransaction,
		id)
	err = row.Scan(
		&date,
		&transType,
		&sum,
		&comment,
	)
	if row.Err() != nil {
		return date, transType, sum, comment, fmt.Errorf("%s: %w", op, row.Err())
	}

	return date, transType, sum, comment, nil
}
