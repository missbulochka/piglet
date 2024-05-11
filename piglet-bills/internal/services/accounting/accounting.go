package accounting

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"strconv"
	"time"

	"github.com/shopspring/decimal"

	models "piglet-bills-service/internal/domain/model"
	"piglet-bills-service/internal/storage"
)

type Accounting struct {
	log          *slog.Logger
	billSaver    BillSaver
	billProvider BillProvider
}

type BillSaver interface {
	SaveBill(
		ctx context.Context,
		billType bool,
		billName string,
		goalSum decimal.Decimal,
		date time.Time,
		monthlyPayment decimal.Decimal,
	) (bill models.Bill, err error)
	UpdateBill(
		ctx context.Context,
		id string,
		billName string,
		currentSum decimal.Decimal,
		billStatus bool,
		goalSum decimal.Decimal,
		date time.Time,
		monthlyPayment decimal.Decimal,
	) (bill models.Bill, err error)
	DeleteBill(ctx context.Context, id string) (err error)
}

type BillProvider interface {
	ReturnBill(
		ctx context.Context, billId string) (bill models.Bill, err error)
	ReturnSomeBills(ctx context.Context, billType bool) (bills []*models.Bill, err error)
	VerifyBill(ctx context.Context, id string) (billType bool, err error)
}

var (
	ErrBillExists   = errors.New("bill already exists")
	ErrBillNotFound = errors.New("bill not found")
)

// New returns a new intarface of the Accounting service
func New(
	log *slog.Logger,
	billSaver BillSaver,
	billProvider BillProvider,
) *Accounting {
	return &Accounting{
		log:          log,
		billSaver:    billSaver,
		billProvider: billProvider,
	}
}

// CreateBill create new bill in the system and returns bill
// If bill with given name already exists, returns error
func (a *Accounting) CreateBill(
	ctx context.Context,
	billType bool,
	billName string,
	goalSum decimal.Decimal,
	date time.Time,
) (bill models.Bill, err error) {
	const op = "pigletBills | accounting.SaveBill"

	log := a.log.With(
		slog.String("op", op),
		// These may be things that are not profitable for business to log
		slog.String("billType", strconv.FormatBool(billType)),
		slog.String("billName", billName),
	)

	monthlyPayment := decimal.New(0, 0)
	if billType == false {
		if monthlyPayment, err = countPayment(date, goalSum); err != nil {
			log.Warn("something wrong with date", err)

			return models.Bill{}, fmt.Errorf("%s: %w", op, err)
		}
	}

	log.Info("saving bill")

	bill, err = a.billSaver.SaveBill(
		ctx,
		billType,
		billName,
		goalSum,
		date,
		monthlyPayment,
	)
	if err != nil {
		if errors.Is(err, storage.ErrBillExists) {
			log.Warn("bill already exists", err)

			return models.Bill{}, fmt.Errorf("%s: %w", op, ErrBillExists)
		}

		log.Error("failed to save bill", err)

		return models.Bill{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("saved bill")
	return bill, nil
}

// GetSomeBills retrieves bills from the system by their types and returns them.
// If an error occurs, it returns an error.
func (a *Accounting) GetSomeBills(ctx context.Context, billType bool) (bills []*models.Bill, err error) {
	const op = "pigletBills | accounting.GetSomeBills"
	log := a.log.With(slog.String("op", op))

	log.Info("searching bills")

	bills, err = a.billProvider.ReturnSomeBills(ctx, billType)
	if err != nil {
		log.Error("failed to search bills", err)

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("searched bills")
	return bills, nil
}

// GetBill retrieves a bill from the system by its name or uuid and returns it.
// If a bill with the given uuid or name does not exist, returns an error.
func (a *Accounting) GetBill(
	ctx context.Context,
	billId string,
) (bill models.Bill, err error) {
	const op = "pigletBills | accounting.GetBill"

	log := a.log.With(
		slog.String("op", op),
		// These may be things that are not profitable for business to log
		slog.String("billId", billId),
	)

	log.Info("searching bill")

	bill, err = a.billProvider.ReturnBill(ctx, billId)
	if err != nil {
		if errors.Is(err, storage.ErrBillNotFound) {
			log.Warn("bill not found", err)

			return models.Bill{}, fmt.Errorf("%s: %w", op, ErrBillExists)
		}

		log.Error("failed to search bill", err)

		return models.Bill{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("searched bill")
	return bill, nil
}

// UpdateBill update bill in the system and returns it
// If bill with given id doesn't exist, returns error
func (a *Accounting) UpdateBill(
	ctx context.Context,
	id string,
	billName string,
	currentSum decimal.Decimal,
	billStatus bool,
	goalSum decimal.Decimal,
	date time.Time,
) (bill models.Bill, err error) {
	const op = "pigletBills | accounting.UpdateBill"

	log := a.log.With(
		slog.String("op", op),
		// These may be things that are not profitable for business to log
		slog.String("bill id", id),
	)

	billType, err := a.billProvider.VerifyBill(ctx, id)
	if err != nil {
		if errors.Is(err, ErrBillNotFound) {
			log.Warn("bill doesn't exist", err)

			return models.Bill{}, fmt.Errorf("%s: %w", op, ErrBillExists)
		}
		log.Error("failed to update bill", err)

		return models.Bill{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("updating bill")

	monthlyPayment := decimal.New(0, 0)
	if billType == false {
		if monthlyPayment, err = countPayment(date, goalSum.Sub(currentSum)); err != nil {
			log.Warn("something wrong with date", err)

			return models.Bill{}, fmt.Errorf("%s: %w", op, err)
		}
	}

	bill, err = a.billSaver.UpdateBill(
		ctx,
		id,
		billName,
		currentSum,
		billStatus,
		goalSum,
		date,
		monthlyPayment,
	)
	if err != nil {
		log.Error("failed to update bill", err)

		return models.Bill{}, fmt.Errorf("%s: %w", op, err)
	}
	bill.BillType = billType
	bill.ID = id

	log.Info("bill updated")
	return bill, nil
}

// DeleteBill delete bill in the system and returns success
// If bill with given id doesn't exist, returns error
func (a *Accounting) DeleteBill(ctx context.Context, id string) (success bool, err error) {
	const op = "pigletBills | accounting.DeleteBill"

	log := a.log.With(
		slog.String("op", op),
		// These may be things that are not profitable for business to log
		slog.String("billId", id),
	)

	log.Info("deleting bill")

	err = a.billSaver.DeleteBill(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrBillNotFound) {
			log.Warn("bill not found", err)

			return false, fmt.Errorf("%s: %w", op, ErrBillExists)
		}
		log.Error("failed to delete bill", err)

		return false, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("bill deleted")
	return true, nil
}

// VerifyBill find bill in the system and returns its type
// If bill with given id doesn't exist, returns error
func (a *Accounting) VerifyBill(ctx context.Context, id string) (success bool, err error) {
	const op = "pigletBills | accounting.VerifyBill"

	log := a.log.With(
		slog.String("op", op),
		// These may be things that are not profitable for business to log
		slog.String("bill id", id),
	)

	log.Info("bill verifying")

	_, err = a.billProvider.VerifyBill(ctx, id)
	if err != nil {
		log.Error("failed to search bill", err)

		return false, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("end of process")

	return true, nil
}

// FixBillSum find bill in the system and update its sum
// If bill with given id doesn't exist, returns error
func (a *Accounting) FixBillSum(ctx context.Context, id string, sum decimal.Decimal) (err error) {
	const op = "pigletBills | accounting.FixVerifySum"

	log := a.log.With(
		slog.String("op", op),
		// These may be things that are not profitable for business to log
		slog.String("bill id", id),
	)

	// HACK: обработка ошибок в случае, если не найден счет
	bill, err := a.billProvider.ReturnBill(ctx, id)
	if err != nil {
		if errors.Is(err, ErrBillNotFound) {
			log.Info("bill doesn't exist")

			return fmt.Errorf("%s: %w", op, err)
		}
		log.Error("failed to update bill", err)

		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("updating bill")

	// HACK: проверка вычислений при отрицательных суммах
	bill.CurrentSum = bill.CurrentSum.Add(sum)
	// HACK: перерасчет ежемесячного платежа при вводе бОльшей суммы

	// HACK: обновление конкретного счета
	bill, err = a.billSaver.UpdateBill(
		ctx,
		id,
		bill.Name,
		bill.CurrentSum,
		bill.BillStatus,
		bill.GoalSum,
		bill.Date,
		bill.MonthlyPayment,
	)
	if err != nil {
		log.Info("bill doesn't updated")

		return fmt.Errorf("%s: %w", op, err)
	}
	log.Info("bill updated")

	return nil
}

func countPayment(futureDate time.Time, sum decimal.Decimal) (monthlyPayment decimal.Decimal, err error) {
	// HACK: подумать над функцией поиска (или найти библиотеку)
	months := math.Ceil(time.Until(futureDate).Hours() / 24 / 30)

	if months == 0 {
		return sum, nil
	}

	monthlyPayment = sum.Div(decimal.New(int64(months), 0))

	return monthlyPayment.Round(0), nil
}
