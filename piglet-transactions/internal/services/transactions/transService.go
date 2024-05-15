package transactions

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"piglet-transactions-service/internal/domain/models"
)

type Transactions struct {
	log              *slog.Logger
	transSaver       TransactionSaver
	transProvider    TransactionProvider
	categoryProvider CategoryProvider
}

type TransactionSaver interface {
	SaveTransaction(ctx context.Context, trans models.Transaction) (err error)
	DeleteTransaction(ctx context.Context, id uuid.UUID, transType uint8) (err error)
}

type TransactionProvider interface {
	DefaultTransInfo(ctx context.Context, id uuid.UUID) (
		date time.Time,
		transType uint8,
		sum decimal.Decimal,
		comment string,
		err error)
	GetTransaction(ctx context.Context, id uuid.UUID, trans *models.Transaction) (err error)
}

type CategoryProvider interface {
	GetCategory(ctx context.Context, id uuid.UUID) (category models.Category, err error)
}

// New returns a new intarface of the Transactions service
func New(
	log *slog.Logger,
	transSaver TransactionSaver,
	transProvider TransactionProvider,
	categoryProvider CategoryProvider,
) *Transactions {
	return &Transactions{
		log:              log,
		transSaver:       transSaver,
		transProvider:    transProvider,
		categoryProvider: categoryProvider,
	}
}
