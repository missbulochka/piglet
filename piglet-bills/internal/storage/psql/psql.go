package psql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	models "piglet-bills-service/internal/domain/model"
	"piglet-bills-service/internal/storage"
)

type Storage struct {
	db *sql.DB
}

const (
	emptySumValue = 0
	openAccount   = true
)

func New(
	dbHost string,
	dbPort string,
	dbUser string,
	dbPassword string,
	dbName string,
) (*Storage, error) {
	const op = "piglet-bills | storage.psql.New"

	db, err := sql.Open(
		"postgres",
		fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			dbHost, dbPort, dbUser, dbPassword, dbName))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	fmt.Println("successfully connected to psql")

	return &Storage{db: db}, nil
}

func (s *Storage) SaveBill(
	ctx context.Context,
	billType bool,
	billName string,
	goalSum decimal.Decimal,
	date time.Time,
	monthlyPayment decimal.Decimal,
) (bill models.Bill, err error) {
	const op = "piglet-bills | storage.psql.SaveBill"

	id := uuid.New().String()
	row := s.db.QueryRowContext(ctx, storage.CreateBill, id, billName, emptySumValue, billType)
	err = row.Scan(
		&bill.ID,
		&bill.Name,
		&bill.CurrentSum,
		&bill.BillType,
	)
	if err != nil {
		return bill, fmt.Errorf("%s: %w", op, err)
	}

	if billType {
		row = s.db.QueryRowContext(ctx, storage.CreateAccount, bill.ID, openAccount)
		err = row.Scan(
			&bill.BillStatus,
		)
	} else {
		row = s.db.QueryRowContext(ctx, storage.CreateGoals, bill.ID, goalSum, date, monthlyPayment)
		err = row.Scan(
			&bill.GoalSum,
			&bill.Date,
			&bill.MonthlyPayment,
		)
	}
	if err != nil {
		// TODO: удалить запись в bills,если не удалось создать запись
		return bill, fmt.Errorf("%s: %w", op, err)
	}

	return bill, err
}

func (s *Storage) BillReturner(
	ctx context.Context,
	billId string,
	billName string,
) (bill models.Bill, err error) {
	const op = "piglet-bills | storage.psql.BillReturner"

	// HACK: обработка ошибки парсинга uuid
	row := s.db.QueryRowContext(ctx, storage.GetOneBill, billId, billName)
	err = row.Scan(
		&bill.ID,
		&bill.Name,
		&bill.CurrentSum,
		&bill.BillType,
	)
	if err != nil {
		return bill, fmt.Errorf("%s: %w", op, err)
	}

	if bill.BillType {
		row = s.db.QueryRowContext(ctx, storage.GetOneAccount, bill.ID)
		err = row.Scan(&bill.BillStatus)
	} else {
		row = s.db.QueryRowContext(ctx, storage.GetOneGoal, bill.ID)
		err = row.Scan(
			&bill.GoalSum,
			&bill.Date,
			&bill.MonthlyPayment,
		)
	}

	if err != nil {
		return bill, fmt.Errorf("%s: %w", op, err)
	}

	return bill, err
}

func (s *Storage) SomeBillsReturner(ctx context.Context, billType bool) (bills []*models.Bill, err error) {
	const op = "piglet-bills | storage.psql.SomeBillsReturner"

	// HACK: подумать о более красивом решении
	rows, err := s.db.QueryContext(ctx, storage.GetSomeBills, billType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var typedRows *sql.Rows
	if billType {
		typedRows, err = s.db.QueryContext(ctx, storage.GetAllAccounts)
	} else {
		typedRows, err = s.db.QueryContext(ctx, storage.GetAllGoals)
	}
	if err != nil {
		return nil, err
	}
	defer typedRows.Close()

	// HACK: оптимизировать
	for rows.Next() {
		typedRows.Next()
		var b models.Bill
		if err = rows.Scan(
			&b.ID,
			&b.Name,
			&b.CurrentSum,
			&b.BillType,
		); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		if billType {
			if err = typedRows.Scan(&b.BillStatus); err != nil {
				return nil, fmt.Errorf("%s: %w", op, err)
			}
		} else {
			if err = typedRows.Scan(
				&b.GoalSum,
				&b.Date,
				&b.MonthlyPayment,
			); err != nil {
				return nil, fmt.Errorf("%s: %w", op, err)
			}
		}
		bills = append(bills, &b)
	}

	if err = rows.Close(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if err = typedRows.Close(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if err = typedRows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return bills, nil
}
