package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	"piglet-transactions-service/internal/domain/models"
	"piglet-transactions-service/internal/storage"
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
	case 1:
		row = s.db.QueryRowContext(
			ctx,
			storage.InsertIncome,
			trans.Id,
			trans.IdCategory,
			trans.IdBillTo,
			trans.Person,
			trans.Repeat,
		)
	case 2:
		row = s.db.QueryRowContext(
			ctx,
			storage.InsertExpense,
			trans.Id,
			trans.IdCategory,
			trans.IdBillFrom,
			trans.Person,
			trans.Repeat,
		)
	case 3:
		row = s.db.QueryRowContext(
			ctx,
			storage.InsertDebt,
			trans.Id,
			trans.DebtType,
			trans.IdBillFrom,
			trans.IdBillTo,
			trans.Person,
		)
	case 4:
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

func (s *Storage) DeleteTransaction(ctx context.Context, id uuid.UUID, transType uint8) (err error) {
	const op = "piglet-transactions | storage.postgres.DeleteTransaction"

	var row *sql.Row

	switch transType {
	case 1:
		row = s.db.QueryRowContext(ctx, storage.DeleteIncome, id)
	case 2:
		row = s.db.QueryRowContext(ctx, storage.DeleteExpenses, id)
	case 3:
		row = s.db.QueryRowContext(ctx, storage.DeleteDebt, id)
	case 4:
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

	switch trans.TransType {
	case 1:
		row = s.db.QueryRowContext(ctx, storage.GetOneIncome, id)
		err = row.Scan(
			&trans.IdCategory,
			&trans.IdBillTo,
			&trans.Person,
			&trans.Repeat,
		)
	case 2:
		row = s.db.QueryRowContext(ctx, storage.GetOneExpense, id)
		err = row.Scan(
			&trans.IdCategory,
			&trans.IdBillFrom,
			&trans.Person,
			&trans.Repeat,
		)
	case 3:
		row = s.db.QueryRowContext(ctx, storage.GetOneDebt, id)
		err = row.Scan(
			&trans.DebtType,
			&trans.IdBillFrom,
			&trans.IdBillTo,
			&trans.Person,
		)
	case 4:
		row = s.db.QueryRowContext(ctx, storage.GetOneTransfer, id)
		err = row.Scan(
			&trans.IdBillFrom,
			&trans.IdBillTo,
		)
	default:
		return fmt.Errorf("%s: unknown error", op)
	}

	if err != nil {
		return fmt.Errorf("%s: unknown error", op)
	}

	return nil
}
