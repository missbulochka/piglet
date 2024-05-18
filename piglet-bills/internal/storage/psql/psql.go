package psql

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/shopspring/decimal"

	models "piglet-bills-service/internal/domain/model"
	"piglet-bills-service/internal/storage"
)

type Storage struct {
	db            *sql.DB
	billsMutex    sync.Mutex
	accountsMutex sync.Mutex
	goalsMutex    sync.Mutex
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
	s.billsMutex.Lock()
	row := s.db.QueryRowContext(ctx, storage.CreateBill, id, billName, emptySumValue, billType)
	err = row.Scan(
		&bill.ID,
		&bill.Name,
		&bill.CurrentSum,
		&bill.BillType,
	)
	s.billsMutex.Unlock()
	if err != nil {
		return bill, fmt.Errorf("%s: %w", op, err)
	}

	if billType {
		s.accountsMutex.Lock()
		row = s.db.QueryRowContext(ctx, storage.CreateAccount, bill.ID, openAccount)
		s.accountsMutex.Unlock()
		err = row.Scan(
			&bill.BillStatus,
		)
	} else {
		s.goalsMutex.Lock()
		row = s.db.QueryRowContext(ctx, storage.CreateGoals, bill.ID, goalSum, date, monthlyPayment)
		s.goalsMutex.Unlock()
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

func (s *Storage) ReturnBill(
	ctx context.Context,
	billId string,
) (bill models.Bill, err error) {
	const op = "piglet-bills | storage.psql.ReturnBill"

	// HACK: обработка ошибки парсинга uuid
	s.billsMutex.Lock()
	row := s.db.QueryRowContext(ctx, storage.GetOneBill, billId)
	s.billsMutex.Unlock()
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
		s.accountsMutex.Lock()
		row = s.db.QueryRowContext(ctx, storage.GetOneAccount, bill.ID)
		s.accountsMutex.Unlock()
		err = row.Scan(&bill.BillStatus)
	} else {
		s.goalsMutex.Lock()
		row = s.db.QueryRowContext(ctx, storage.GetOneGoal, bill.ID)
		s.goalsMutex.Unlock()
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

func (s *Storage) ReturnSomeBills(ctx context.Context, billType bool) (bills []*models.Bill, err error) {
	const op = "piglet-bills | storage.psql.ReturnSomeBills"

	// HACK: подумать о более красивом решении
	s.billsMutex.Lock()
	rows, err := s.db.QueryContext(ctx, storage.GetSomeBills, billType)
	s.billsMutex.Unlock()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var typedRows *sql.Rows
	if billType {
		s.accountsMutex.Lock()
		typedRows, err = s.db.QueryContext(ctx, storage.GetAllAccounts)
		s.accountsMutex.Unlock()
	} else {
		s.goalsMutex.Lock()
		typedRows, err = s.db.QueryContext(ctx, storage.GetAllGoals)
		s.goalsMutex.Unlock()
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

func (s *Storage) UpdateBill(
	ctx context.Context,
	id string,
	billName string,
	currentSum decimal.Decimal,
	billStatus bool,
	goalSum decimal.Decimal,
	date time.Time,
	monthlyPayment decimal.Decimal,
) (bill models.Bill, err error) {
	const op = "piglet-bills | storage.psql.UpdateBill"

	s.billsMutex.Lock()
	row := s.db.QueryRowContext(
		ctx,
		storage.UpdateBill,
		id,
		billName,
		currentSum,
	)
	s.billsMutex.Unlock()
	err = row.Scan(
		&bill.Name,
		&bill.CurrentSum,
		&bill.BillType,
	)
	if err != nil {
		return bill, fmt.Errorf("%s: %w", op, err)
	}

	if bill.BillType {
		s.accountsMutex.Lock()
		row = s.db.QueryRowContext(
			ctx,
			storage.UpdateAccount,
			id,
			billStatus,
		)
		s.accountsMutex.Unlock()
		err = row.Scan(
			&bill.BillStatus,
		)
	} else {
		s.goalsMutex.Lock()
		row = s.db.QueryRowContext(
			ctx,
			storage.UpdateGoal,
			id,
			goalSum,
			date,
			monthlyPayment,
		)
		s.goalsMutex.Unlock()
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

func (s *Storage) DeleteBill(ctx context.Context, id string) (err error) {
	const op = "piglet-bills | storage.psql.DeleteBill"

	// HACK: восстановление строки в случае, если не удалось удалить записи из зависимых таблиц
	var billType bool
	s.billsMutex.Lock()
	row := s.db.QueryRowContext(ctx, storage.VerifyBill, id)
	s.billsMutex.Unlock()
	err = row.Scan(&billType)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if billType {
		s.accountsMutex.Lock()
		row = s.db.QueryRowContext(ctx, storage.DeleteAccount, id)
		s.accountsMutex.Unlock()
	} else {
		s.goalsMutex.Lock()
		row = s.db.QueryRowContext(ctx, storage.DeleteGoal, id)
		s.goalsMutex.Unlock()
	}
	if row.Err() != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	s.billsMutex.Lock()
	row = s.db.QueryRowContext(ctx, storage.DeleteBill, id)
	s.billsMutex.Unlock()
	if row.Err() != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) VerifyBill(
	ctx context.Context,
	id string,
) (billType bool, err error) {
	const op = "piglet-bills | storage.psql.VerifyBill"

	s.billsMutex.Lock()
	row := s.db.QueryRowContext(ctx, storage.VerifyBill, id)
	s.billsMutex.Unlock()
	err = row.Scan(&billType)
	// HACK: в случае ошибки возвращает false (что является в БД типом goals)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return billType, nil
}
