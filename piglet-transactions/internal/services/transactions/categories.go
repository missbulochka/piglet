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
	panic("implement me")
}

// UpdateCategory update exist category in the system and returns it
// If category with given id doesn't exist, returns error
func (t *Transactions) UpdateCategory(ctx context.Context, cat *models.Category) (err error) {
	panic("implement me")
}

// DeleteCategory delete exist category in the system
// If category with given id doesn't exist, returns error
func (t *Transactions) DeleteCategory(ctx context.Context, id uuid.UUID) (err error) {
	panic("implement me")
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
	panic("implement me")
}
