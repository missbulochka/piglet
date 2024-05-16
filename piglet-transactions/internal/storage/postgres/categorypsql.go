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

	s.catMutex.Lock()
	row := s.db.QueryRowContext(
		ctx,
		storage.InsertCategory,
		cat.Id,
		cat.CategoryType,
		cat.Name,
		cat.Mandatory,
	)
	s.catMutex.Unlock()
	if row.Err() != nil {
		return fmt.Errorf("%s: %w", op, row.Err())
	}

	return nil
}

func (s *Storage) UpdateCategory(ctx context.Context, cat models.Category) (err error) {
	const op = "piglet-bills | storage.psql.UpdateCategory"

	s.catMutex.Lock()
	row := s.db.QueryRowContext(
		ctx,
		storage.UpdateCategory,
		cat.Id,
		cat.CategoryType,
		cat.Name,
		cat.Mandatory,
	)
	s.catMutex.Unlock()
	if row.Err() != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetCategory(
	ctx context.Context,
	search interface{},
) (category models.Category, err error) {
	const op = "piglet-transactions | storage.postgres.GetCategory"

	s.catMutex.Lock()
	row := s.db.QueryRowContext(ctx, storage.GetCategory, search)
	s.catMutex.Unlock()
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

func (s *Storage) GetAllCategories(ctx context.Context, cat *[]*models.Category) (err error) {
	const op = "piglet-transactions | storage.postgres.GetSomeCategories"

	s.catMutex.Lock()
	rows, err := s.db.QueryContext(ctx, storage.GetAllCategories)
	s.catMutex.Unlock()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	defer func() {
		if err = rows.Close(); err != nil {
			fmt.Printf("%s: %v", op, err)
		}
		if err = rows.Err(); err != nil {
			fmt.Printf("%s: %v", op, err)
		}
	}()

	for rows.Next() {
		var c models.Category
		if err = rows.Scan(
			&c.Id,
			&c.CategoryType,
			&c.Name,
			&c.Mandatory,
		); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		*cat = append(*cat, &c)
	}

	return nil
}

func (s *Storage) DeleteCategory(ctx context.Context, id uuid.UUID) (err error) {
	const op = "piglet-transactions | storage.postgres.DeleteCategory"

	s.catMutex.Lock()
	row := s.db.QueryRowContext(ctx, storage.DeleteCategory, id)
	s.catMutex.Unlock()

	if row.Err() != nil {
		return fmt.Errorf("%s: %w", op, row.Err())
	}

	return nil
}
