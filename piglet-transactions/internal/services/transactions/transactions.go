package transactions

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"

	"piglet-transactions-service/internal/domain/models"
)

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
		log.Error("failed to save transaction", err)

		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("transaction saved")

	return nil
}

// DeleteTransaction delete transaction in the system
// If transaction with given id doesn't exist, returns error
func (t *Transactions) DeleteTransaction(ctx context.Context, id uuid.UUID) error {
	const op = "pigletTransactions | transactions.DeleteTransaction"
	log := t.log.With(slog.String("op", op))

	_, transType, _, _, err := t.transProvider.DefaultTransInfo(ctx, id)
	if err != nil {
		log.Error("transaction doesn't exist", err)

		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("deleting transaction")

	if err = t.transSaver.DeleteTransaction(ctx, id, transType); err != nil {
		log.Error("failed to delete transaction", err)

		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("transaction deleted")

	return nil
}

// GetTransaction search transaction in the system
// If transaction with given id doesn't exist, returns error
func (t *Transactions) GetTransaction(ctx context.Context, id uuid.UUID) (trans models.Transaction, err error) {
	const op = "pigletTransactions | transactions.GetTransaction"
	log := t.log.With(slog.String("op", op))

	trans.Date, trans.TransType, trans.Sum, trans.Comment, err = t.transProvider.DefaultTransInfo(ctx, id)
	if err != nil {
		log.Error("failed to search transaction", err)

		return trans, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("receiving transaction")

	if err = t.transProvider.GetTransaction(ctx, id, &trans); err != nil {
		log.Error("failed to get transaction", err)

		return trans, fmt.Errorf("%s: %w", op, err)
	}

	trans.Id = id

	log.Info("transaction received")

	return trans, nil
}
