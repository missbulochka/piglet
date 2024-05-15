package transactions

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"log/slog"

	"piglet-transactions-service/internal/domain/models"
)

// CreateCategory create new category in the system and returns it
// If category with given names don't exist, returns error
func (t *Transactions) CreateCategory(ctx context.Context, cat *models.Category) (err error) {
	const op = "pigletTransactions | transactions.CreateCategory"
	log := t.log.With(slog.String("op", op))

	cat.Id = uuid.New()

	log.Info("saving category")

	if err = t.categorySaver.SaveCategory(ctx, *cat); err != nil {
		log.Error("failed to save transaction", err)

		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("category saved")

	return nil
}

// UpdateCategory update exist category in the system and returns it
// If category with given id doesn't exist, returns error
func (t *Transactions) UpdateCategory(ctx context.Context, cat *models.Category) (err error) {
	const op = "pigletBills | accounting.UpdateCategory"
	log := t.log.With(slog.String("op", op))

	log.Info("updating category")

	if err = t.categorySaver.UpdateCategory(ctx, *cat); err != nil {
		log.Error("failed to update bill", err)

		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("category updated")

	return nil
}

// DeleteCategory delete exist category in the system
// If category with given id doesn't exist, returns error
func (t *Transactions) DeleteCategory(ctx context.Context, id uuid.UUID) (err error) {
	const op = "pigletTransactions | transactions.DeleteCategory"
	log := t.log.With(slog.String("op", op))

	log.Info("deleting category")

	if err = t.categorySaver.DeleteCategory(ctx, id); err != nil {
		log.Error("failed to delete transaction", err)

		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("category deleted")

	return nil
}

// GetCategory return exist category in the system
// If category with given id doesn't exist, returns error
func (t *Transactions) GetCategory(ctx context.Context, id uuid.UUID) (cat models.Category, err error) {
	const op = "pigletTransactions | transactions.GetCategory"
	log := t.log.With(slog.String("op", op))

	log.Info("receiving category")

	cat, err = t.categoryProvider.GetCategory(ctx, id)
	if err != nil {
		log.Error("failed to get category", err)

		return cat, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("category received")

	return cat, nil
}

// GetAllCategories search all categories in the system
// If something go wrong, returns error
func (t *Transactions) GetAllCategories(ctx context.Context) (cat []*models.Category, err error) {
	const op = "pigletTransactions | transactions.GetAllCategories"
	log := t.log.With(slog.String("op", op))

	log.Info("receiving categories")

	if err = t.categoryProvider.GetAllCategories(ctx, &cat); err != nil {
		log.Error("failed to get transactions", err)

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("categories received")

	return cat, nil
}
