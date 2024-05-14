package postgres

import (
	"context"
	"piglet-transactions-service/internal/domain/models"
)

func (s *Storage) SaveTransaction(
	ctx context.Context,
	trans models.Transaction,
) (err error) {
	panic("need connect")
}
