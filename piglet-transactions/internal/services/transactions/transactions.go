package transactions

import (
	"context"
	"fmt"
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

// CreateTransaction create new transaction in the system and returns it
// If bill or category with given names don't exist, returns error
func (t *Transactions) CreateTransaction(
	ctx context.Context,
	trans *models.Transaction,
) (err error) {
	const op = "pigletTransactions | transactions.CreateTransaction"
	log := t.log.With(slog.String("op", op))

	trans.Id = uuid.New()

	if trans.TransType == 1 || trans.TransType == 2 {
		log.Info("verifying category")
		if _, err = t.categoryProvider.GetCategory(ctx, trans.IdCategory); err != nil {
			log.Error("failed to verify category", err)

			return fmt.Errorf("%s: %w", op, err)
		}
	}

	log.Info("saving transaction")

	if err = t.transSaver.SaveTransaction(ctx, *trans); err != nil {
		//TODO: проверка на ошибку "счет не существует
		log.Error("failed to save transaction", err)

		return err
	}

	log.Info("transaction saved")

	return nil
}

func (t *Transactions) DeleteTransaction(ctx context.Context, id uuid.UUID) error {
	const op = "pigletTransactions | transactions.DeleteTransaction"
	log := t.log.With(slog.String("op", op))

	_, transType, _, _, err := t.transProvider.DefaultTransInfo(ctx, id)
	if err != nil {
		log.Error("transaction doesn't exist", err)

		return err
	}

	log.Info("deleting transaction")

	if err = t.transSaver.DeleteTransaction(ctx, id, transType); err != nil {
		log.Error("failed to delete transaction", err)

		return err
	}

	log.Info("transaction deleted")

	return nil
}
