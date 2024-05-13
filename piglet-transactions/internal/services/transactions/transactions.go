package transactions

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Transactions struct {
	log        *slog.Logger
	transSaver TransactionSaver
}

type TransactionSaver interface {
	SaveTransaction(ctx context.Context)
}

func New(log *slog.Logger, transSaver TransactionSaver) *Transactions {
	return &Transactions{
		log:        log,
		transSaver: transSaver,
	}
}

func (t *Transactions) CreateTransaction(
	ctx context.Context,
	date time.Time,
	transType uint8,
	sum decimal.Decimal,
	comment string,
	idCategory uuid.UUID,
	debtType bool,
	idBillTo uuid.UUID,
	idBillFrom uuid.UUID,
	person string,
	repeat bool,
) (err error) {
	panic("implement me")
}
