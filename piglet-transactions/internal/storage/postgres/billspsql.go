package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"piglet-transactions-service/internal/storage"
)

func (s *Storage) SaveBill(ctx context.Context, id uuid.UUID, billStatus bool) (err error) {
	const op = "piglet-bills | storage.psql.SaveBill"

	row := s.db.QueryRowContext(ctx, storage.InsertBill, id, billStatus)
	if row.Err() != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) UpdateBill(ctx context.Context, id uuid.UUID, billStatus bool) (err error) {
	const op = "piglet-bills | storage.psql.UpdateBill"

	row := s.db.QueryRowContext(ctx, storage.UpdateBill, id, billStatus)
	if row.Err() != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetBill(ctx context.Context, id uuid.UUID) (status bool, err error) {
	const op = "piglet-transactions | storage.postgres.GetBill"

	row := s.db.QueryRowContext(ctx, storage.GetBill, id)
	if err = row.Scan(&status); err != nil {
		return status, fmt.Errorf("%s: %w", op, err)
	}

	return status, nil
}

func (s *Storage) DeleteBill(ctx context.Context, id uuid.UUID) (err error) {
	const op = "piglet-transactions | storage.postgres.DeleteBill"

	row := s.db.QueryRowContext(ctx, storage.DeleteBill, id)

	if row.Err() != nil {
		return fmt.Errorf("%s: %w", op, row.Err())
	}

	return nil
}
