package transactions

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"

	"piglet-transactions-service/internal/domain/models"
)

type Transactions struct {
	log              *slog.Logger
	transSaver       TransactionSaver
	categoryProvider CategoryProvider
}

type TransactionSaver interface {
	SaveTransaction(ctx context.Context, trans models.Transaction) (err error)
}

type CategoryProvider interface {
	GetCategory(ctx context.Context, id uuid.UUID) (category models.Category, err error)
}

// New returns a new intarface of the Transactions service
func New(log *slog.Logger, transSaver TransactionSaver, categoryProvider CategoryProvider) *Transactions {
	return &Transactions{
		log:        log,
		transSaver: transSaver,
		categoryProvider: categoryProvider,
	}
}

// CreateTransaction create new transaction in the system and returns it
// If bill or category with given names don't exist, returns error
func (t *Transactions) CreateTransaction(
	ctx context.Context,
	trans *models.Transaction,
) (err error) {
	const op = "pigletTransactions | transactions.CreateTransaction"
	log := t.log.With(slog.String("op", op))

	trans.Id = uuid.New()

	if trans.Id != uuid.Nil {
		if _, err = t.categoryProvider.GetCategory(ctx, trans.IdCategory); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	} else {
		trans.Id = uuid.Nil
	}

	log.Info("Saving transaction")

	if err = t.transSaver.SaveTransaction(ctx, *trans); err != nil {
		//TODO: проверка на ошибку "счет не существует"
	}

	log.Info("Transaction saved")

	return nil
}
