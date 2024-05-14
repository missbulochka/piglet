package postgres

import (
	"context"

	"github.com/google/uuid"

	"piglet-transactions-service/internal/domain/models"
)

func (s *Storage) GetCategory(
	ctx context.Context,
	id uuid.UUID,
) (category models.Category, err error) {
	panic("need connect")
}
