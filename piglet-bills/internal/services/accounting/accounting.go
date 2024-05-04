package accountingrpc

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	models "piglet-bills-service/internal/domain/model"
	"piglet-bills-service/internal/storage"
	"strconv"
)

type Accounting struct {
	log       *slog.Logger
	billSaver BillSaver
}

type BillSaver interface {
	SaveBill(ctx context.Context,
		billType bool,
		billName string,
		date string,
		monthlyPayment float32,
	) (bill models.Bill, err error)
}

// TODO: дописать интерфейс BillProvider

var (
	ErrUserExists = errors.New("user already exists")
)

// New returns a new intarface of the Accounting service
func New(
	log *slog.Logger,
	billSaver BillSaver,
) *Accounting {
	return &Accounting{
		billSaver: billSaver,
		log:       log,
	}
}

// CreateBill create new bill in the system and returns bill
// If bill with given name already exists, returns error
func (a *Accounting) CreateBill(
	ctx context.Context,
	billType bool,
	billName string,
	date string,
) (models.Bill, error) {
	const op = "pigletBills | accounting.saveBill"

	log := a.log.With(
		slog.String("op", op),
		// These may be things that are not profitable for business to log
		slog.String("billType", strconv.FormatBool(billType)),
		slog.String("billName", billName),
		slog.String("date", date),
	)

	log.Info("saving bill")

	var monthlyPayment float32
	if billType == true {
		monthlyPayment = 0
	} else {
		// TODO: вычисление помесячного платежа
		monthlyPayment = 0
	}

	bill, err := a.billSaver.SaveBill(ctx, billType, billName, date, monthlyPayment)
	if err != nil {
		if errors.Is(err, storage.ErrBillExists) {
			log.Warn("bill already exists", err)

			return models.Bill{}, fmt.Errorf("%s: %w", op, ErrUserExists)
		}

		log.Error("failed to save bill", err)

		return models.Bill{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("saved bill")
	return bill, nil
}
