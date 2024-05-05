package accounting

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/shopspring/decimal"

	models "piglet-bills-service/internal/domain/model"
	"piglet-bills-service/internal/storage"
)

type Accounting struct {
	log       *slog.Logger
	billSaver BillSaver
}

type BillSaver interface {
	SaveBill(
		ctx context.Context,
		billType bool,
		billName string,
		currentSum decimal.Decimal,
		date *timestamp.Timestamp,
		monthlyPayment decimal.Decimal,
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
	currentSum decimal.Decimal,
	date *timestamp.Timestamp,
) (bill models.Bill, err error) {
	const op = "pigletBills | accounting.saveBill"

	log := a.log.With(
		slog.String("op", op),
		// These may be things that are not profitable for business to log
		slog.String("billType", strconv.FormatBool(billType)),
		slog.String("billName", billName),
	)

	monthlyPayment := decimal.New(0, 0)
	if billType == false {
		if monthlyPayment, err = countPayment(date, currentSum); err != nil {
			log.Warn("something wrong with date", err)

			return models.Bill{}, fmt.Errorf("%s: %w", op, err)
		}
	}

	log.Info("saving bill")

	bill, err = a.billSaver.SaveBill(
		ctx,
		billType,
		billName,
		currentSum,
		date,
		monthlyPayment,
	)
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

func countPayment(futureDate *timestamp.Timestamp, sum decimal.Decimal) (monthlyPayment decimal.Decimal, err error) {
	date := time.Unix(futureDate.GetSeconds(), int64(futureDate.GetNanos())).UTC()

	// HACK: подумать над функцией поиска (или найти библиотеку)
	months := int(time.Until(date).Hours() / 24 / 30)

	if months == 0 {
		return sum, nil
	}

	monthlyPayment = sum.Div(decimal.New(int64(months), 0))

	return monthlyPayment.Round(0), nil
}
