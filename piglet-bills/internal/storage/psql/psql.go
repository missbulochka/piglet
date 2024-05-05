package psql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/shopspring/decimal"

	models "piglet-bills-service/internal/domain/model"
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

func (s *Storage) SaveBill(ctx context.Context,
	billType bool,
	billName string,
	currentSum decimal.Decimal,
	date string,
	monthlyPayment decimal.Decimal,
) (bill models.Bill, err error) {
	const op = "piglet-bills | storage.psql.SaveBill"

	// TODO: работа с базой данных

	return models.Bill{}, nil
}
