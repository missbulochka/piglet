package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"piglet-transactions-service/internal/domain/models"
	"piglet-transactions-service/internal/storage"
)

func (s *Storage) SaveCategory(ctx context.Context, cat models.Category) (err error) {
	const op = "piglet-transactions | storage.postgres.SaveCategory"

	row := s.db.QueryRowContext(
		ctx,
		storage.InsertCategory,
		cat.Id,
		cat.CategoryType,
		cat.Name,
		cat.Mandatory,
	)
	if row.Err() != nil {
		return fmt.Errorf("%s: %w", op, row.Err())
	}

	return nil
}

func (s *Storage) GetCategory(
	ctx context.Context,
	id uuid.UUID,
) (category models.Category, err error) {
	const op = "piglet-transactions | storage.postgres.GetCategory"

	row := s.db.QueryRowContext(ctx, storage.GetCategory, id)
	if err = row.Scan(
		&category.Id,
		&category.CategoryType,
		&category.Name,
		&category.Mandatory,
	); err != nil {
		return category, fmt.Errorf("%s: %w", op, err)
	}

	return category, nil
}

func (s *Storage) DeleteCategory(ctx context.Context, id uuid.UUID) (err error) {
	const op = "piglet-transactions | storage.postgres.DeleteCategory"

	row := s.db.QueryRowContext(ctx, storage.DeleteCategory, id)

	if row.Err() != nil {
		return fmt.Errorf("%s: %w", op, row.Err())
	}

	return nil
}
