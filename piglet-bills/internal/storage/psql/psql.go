package psql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	models "piglet-bills-service/internal/domain/model"
	"piglet-bills-service/internal/storage"
)

type Storage struct {
	db *sql.DB
}

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
	currentSum decimal.Decimal,
	date *timestamp.Timestamp,
	monthlyPayment decimal.Decimal,
) (bill models.Bill, err error) {
	const op = "piglet-bills | storage.psql.SaveBill"

	id := uuid.New().String()
	row := s.db.QueryRowContext(ctx, storage.CreateBill, id, billName, currentSum)
	err = row.Scan(
		&bill.ID,
		&bill.Name,
		&bill.CurrentSum,
	)
	if err != nil {
		return bill, fmt.Errorf("%s: %w", op, err)
	}

	if billType == true {
		row = s.db.QueryRowContext(ctx, storage.CreateAccount, bill.ID, true)
		err = row.Scan(
			&bill.BillStatus,
		)
	} else {
		row = s.db.QueryRowContext(ctx, storage.CreateGoals, bill.ID, date, monthlyPayment)
		err = row.Scan(
			&bill.ID,
			&bill.Date,
			&bill.MonthlyPayment,
		)
	}
	if err != nil {
		return bill, fmt.Errorf("%s: %w", op, err)
	}

	return bill, err
}
