package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	"piglet-transactions-service/internal/domain/models"
	"piglet-transactions-service/internal/storage"
)

const (
	transTypeIncome   = 1
	transTypeExpense  = 2
	transTypeDebt     = 3
	transTypeTransfer = 4
)

func (s *Storage) SaveTransaction(
	ctx context.Context,
	trans models.Transaction,
) error {
	const op = "piglet-transactions | storage.postgres.SaveTransaction"

	row := s.db.QueryRowContext(
		ctx,
		storage.InsertTransaction,
		trans.Id,
		trans.Date,
		trans.TransType,
		trans.Sum,
		trans.Comment)
	if row.Err() != nil {
		return fmt.Errorf("%s: %w", op, row.Err())
	}

	switch trans.TransType {
	case transTypeIncome:
		row = s.db.QueryRowContext(
			ctx,
			storage.InsertIncome,
			trans.Id,
			trans.IdCategory,
			trans.IdBillTo,
			trans.Person,
			trans.Repeat,
		)
	case transTypeExpense:
		row = s.db.QueryRowContext(
			ctx,
			storage.InsertExpense,
			trans.Id,
			trans.IdCategory,
			trans.IdBillFrom,
			trans.Person,
			trans.Repeat,
		)
	case transTypeDebt:
		row = s.db.QueryRowContext(
			ctx,
			storage.InsertDebt,
			trans.Id,
			trans.DebtType,
			trans.IdBillFrom,
			trans.IdBillTo,
			trans.Person,
		)
	case transTypeTransfer:
		row = s.db.QueryRowContext(
			ctx,
			storage.InsertTransfer,
			trans.Id,
			trans.IdBillFrom,
			trans.IdBillTo,
		)
	default:
		// TODO: удалить запись в таблице transactions
		return fmt.Errorf("%s: unknown error", op)
	}

	if row.Err() != nil {
		// TODO: удалить запись в таблице transactions
		return fmt.Errorf("%s: %w", op, row.Err())
	}

	return nil
}

func (s *Storage) UpdateTransaction(ctx context.Context, trans models.Transaction) (err error) {
	const op = "piglet-transactions | storage.postgres.UpdateTransaction"

	row := s.db.QueryRowContext(
		ctx,
		storage.UpdateTransaction,
		trans.Id,
		trans.Date,
		trans.TransType,
		trans.Sum,
		trans.Comment)
	if row.Err() != nil {
		return fmt.Errorf("%s: %w", op, row.Err())
	}

	switch trans.TransType {
	case transTypeIncome:
		row = s.db.QueryRowContext(
			ctx,
			storage.UpdateIncome,
			trans.Id,
			trans.IdCategory,
			trans.IdBillTo,
			trans.Person,
			trans.Repeat,
		)
	case transTypeExpense:
		row = s.db.QueryRowContext(
			ctx,
			storage.UpdateExpense,
			trans.Id,
			trans.IdCategory,
			trans.IdBillFrom,
			trans.Person,
			trans.Repeat,
		)
	case transTypeDebt:
		row = s.db.QueryRowContext(
			ctx,
			storage.UpdateDebt,
			trans.Id,
			trans.DebtType,
			trans.IdBillFrom,
			trans.IdBillTo,
			trans.Person,
		)
	case transTypeTransfer:
		row = s.db.QueryRowContext(
			ctx,
			storage.UpdateTransfer,
			trans.Id,
			trans.IdBillFrom,
			trans.IdBillTo,
		)
	default:
		return fmt.Errorf("%s: unknown error", op)
	}

	if row.Err() != nil {
		// TODO: удалить запись в таблице transactions
		return fmt.Errorf("%s: %w", op, row.Err())
	}

	return nil
}

func (s *Storage) DeleteTransaction(ctx context.Context, id uuid.UUID, transType uint8) (err error) {
	const op = "piglet-transactions | storage.postgres.DeleteTransaction"

	var row *sql.Row

	switch transType {
	case transTypeIncome:
		row = s.db.QueryRowContext(ctx, storage.DeleteIncome, id)
	case transTypeExpense:
		row = s.db.QueryRowContext(ctx, storage.DeleteExpenses, id)
	case transTypeDebt:
		row = s.db.QueryRowContext(ctx, storage.DeleteDebt, id)
	case transTypeTransfer:
		row = s.db.QueryRowContext(ctx, storage.DeleteTransfer, id)
	default:
		return fmt.Errorf("%s: unknown error", op)
	}

	if row.Err() != nil {
		return fmt.Errorf("%s: %w", op, row.Err())
	}

	row = s.db.QueryRowContext(ctx, storage.DeleteTransaction, id)

	if row.Err() != nil {
		return fmt.Errorf("%s: %w", op, row.Err())
	}

	return nil
}

func (s *Storage) GetTransaction(
	ctx context.Context,
	id uuid.UUID,
	trans *models.Transaction,
) (err error) {
	const op = "piglet-transactions | storage.postgres.GetTransaction"

	var row *sql.Row

	row = s.db.QueryRowContext(ctx, storage.GetOneTransaction, id)
	err = row.Scan(
		&trans.Date,
		&trans.TransType,
		&trans.Sum,
		&trans.Comment,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	switch trans.TransType {
	case transTypeIncome:
		row = s.db.QueryRowContext(ctx, storage.GetOneIncome, id)
		err = row.Scan(
			&trans.IdCategory,
			&trans.IdBillTo,
			&trans.Person,
			&trans.Repeat,
		)
	case transTypeExpense:
		row = s.db.QueryRowContext(ctx, storage.GetOneExpense, id)
		err = row.Scan(
			&trans.IdCategory,
			&trans.IdBillFrom,
			&trans.Person,
			&trans.Repeat,
		)
	case transTypeDebt:
		row = s.db.QueryRowContext(ctx, storage.GetOneDebt, id)
		err = row.Scan(
			&trans.DebtType,
			&trans.IdBillFrom,
			&trans.IdBillTo,
			&trans.Person,
		)
	case transTypeTransfer:
		row = s.db.QueryRowContext(ctx, storage.GetOneTransfer, id)
		err = row.Scan(
			&trans.IdBillFrom,
			&trans.IdBillTo,
		)
	default:
		return fmt.Errorf("%s: unknown error", op)
	}

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	trans.Id = id

	return nil
}

func (s *Storage) GetSomeTransactions(
	ctx context.Context,
	trans *[]*models.Transaction,
	count uint8) (err error) {
	const op = "piglet-transactions | storage.postgres.GetSomeTransactions"

	rows, err := s.db.QueryContext(ctx, storage.GetSomeTransactions, count)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	defer func() {
		if err = rows.Close(); err != nil {
			fmt.Printf("%s: %v", op, err)
		}
		if err = rows.Err(); err != nil {
			fmt.Printf("%s: %v", op, err)
		}
	}()

	for rows.Next() {
		var tr models.Transaction
		if err = rows.Scan(
			&tr.Id,
			&tr.Date,
			&tr.TransType,
			&tr.Sum,
			&tr.Comment,
		); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		*trans = append(*trans, &tr)
	}

	return nil
}
