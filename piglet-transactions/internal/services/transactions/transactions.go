package transactions

import (
	"context"
	"log/slog"
	"piglet-transactions-service/internal/domain/models"
)

type Transactions struct {
	log        *slog.Logger
	transSaver TransactionSaver
}

type TransactionSaver interface {
	SaveTransaction(ctx context.Context)
}

// New returns a new intarface of the Transactions service
func New(log *slog.Logger, transSaver TransactionSaver) *Transactions {
	return &Transactions{
		log:        log,
		transSaver: transSaver,
	}
}

// CreateTransaction create new transaction in the system and returns it
// If bill or category with given names don't exist, returns error
func (t *Transactions) CreateTransaction(
	ctx context.Context,
	trans models.Transaction,
) (savedTrans models.Transaction, err error) {
	const op = "pigletTransactions | transactions.CreateTransaction"
	log := t.log.With(slog.String("op", op))

	log.Info("Saving transaction")

	// TODO: VerifyBill()

	// TODO: SaveTransaction()

	// TODO: FixBillSum()

	log.Info("Transaction saved")

	return trans, nil
}
