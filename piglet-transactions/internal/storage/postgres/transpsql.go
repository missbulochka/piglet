package postgres

import (
	"context"
	"fmt"

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

	fmt.Println("Inserted")

	return nil
}
