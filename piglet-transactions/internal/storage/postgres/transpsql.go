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

	s.transMutex.Lock()
	row := s.db.QueryRowContext(
		ctx,
		storage.InsertTransaction,
		trans.Id,
		trans.Date,
		trans.TransType,
		trans.Sum,
		trans.Comment)
	s.transMutex.Unlock()
	if row.Err() != nil {
		return fmt.Errorf("%s: %w", op, row.Err())
	}

	switch trans.TransType {
	case transTypeIncome:
		s.incMutex.Lock()
		row = s.db.QueryRowContext(
			ctx,
			storage.InsertIncome,
			trans.Id,
			trans.IdCategory,
			trans.IdBillTo,
			trans.Person,
			trans.Repeat,
		)
		s.incMutex.Unlock()
	case transTypeExpense:
		s.expMutex.Lock()
		row = s.db.QueryRowContext(
			ctx,
			storage.InsertExpense,
			trans.Id,
			trans.IdCategory,
			trans.IdBillFrom,
			trans.Person,
			trans.Repeat,
		)
		s.expMutex.Unlock()
	case transTypeDebt:
		s.debtMutex.Lock()
		row = s.db.QueryRowContext(
			ctx,
			storage.InsertDebt,
			trans.Id,
			trans.DebtType,
			trans.IdBillFrom,
			trans.IdBillTo,
			trans.Person,
		)
		s.debtMutex.Unlock()
	case transTypeTransfer:
		s.transferMutex.Lock()
		row = s.db.QueryRowContext(
			ctx,
			storage.InsertTransfer,
			trans.Id,
			trans.IdBillFrom,
			trans.IdBillTo,
		)
		s.transferMutex.Unlock()
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

	s.transMutex.Lock()
	row := s.db.QueryRowContext(
		ctx,
		storage.UpdateTransaction,
		trans.Id,
		trans.Date,
		trans.TransType,
		trans.Sum,
		trans.Comment)
	s.transMutex.Unlock()
	if row.Err() != nil {
		return fmt.Errorf("%s: %w", op, row.Err())
	}

	switch trans.TransType {
	case transTypeIncome:
		s.incMutex.Lock()
		row = s.db.QueryRowContext(
			ctx,
			storage.UpdateIncome,
			trans.Id,
			trans.IdCategory,
			trans.IdBillTo,
			trans.Person,
			trans.Repeat,
		)
		s.incMutex.Unlock()
	case transTypeExpense:
		s.expMutex.Lock()
		row = s.db.QueryRowContext(
			ctx,
			storage.UpdateExpense,
			trans.Id,
			trans.IdCategory,
			trans.IdBillFrom,
			trans.Person,
			trans.Repeat,
		)
		s.expMutex.Unlock()
	case transTypeDebt:
		s.debtMutex.Lock()
		row = s.db.QueryRowContext(
			ctx,
			storage.UpdateDebt,
			trans.Id,
			trans.DebtType,
			trans.IdBillFrom,
			trans.IdBillTo,
			trans.Person,
		)
		s.debtMutex.Unlock()
	case transTypeTransfer:
		s.transferMutex.Lock()
		row = s.db.QueryRowContext(
			ctx,
			storage.UpdateTransfer,
			trans.Id,
			trans.IdBillFrom,
			trans.IdBillTo,
		)
		s.transferMutex.Unlock()
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
		s.incMutex.Lock()
		row = s.db.QueryRowContext(ctx, storage.DeleteIncome, id)
		s.incMutex.Unlock()
	case transTypeExpense:
		s.expMutex.Lock()
		row = s.db.QueryRowContext(ctx, storage.DeleteExpenses, id)
		s.expMutex.Unlock()
	case transTypeDebt:
		s.debtMutex.Lock()
		row = s.db.QueryRowContext(ctx, storage.DeleteDebt, id)
		s.debtMutex.Unlock()
	case transTypeTransfer:
		s.transferMutex.Lock()
		row = s.db.QueryRowContext(ctx, storage.DeleteTransfer, id)
		s.transferMutex.Unlock()
	default:
		return fmt.Errorf("%s: unknown error", op)
	}

	if row.Err() != nil {
		return fmt.Errorf("%s: %w", op, row.Err())
	}

	s.transMutex.Lock()
	row = s.db.QueryRowContext(ctx, storage.DeleteTransaction, id)
	s.transMutex.Unlock()

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

	s.transMutex.Lock()
	row = s.db.QueryRowContext(ctx, storage.GetOneTransaction, id)
	s.transMutex.Unlock()
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
		s.incMutex.Lock()
		row = s.db.QueryRowContext(ctx, storage.GetOneIncome, id)
		s.incMutex.Unlock()
		err = row.Scan(
			&trans.IdCategory,
			&trans.IdBillTo,
			&trans.Person,
			&trans.Repeat,
		)
	case transTypeExpense:
		s.expMutex.Lock()
		row = s.db.QueryRowContext(ctx, storage.GetOneExpense, id)
		s.expMutex.Unlock()
		err = row.Scan(
			&trans.IdCategory,
			&trans.IdBillFrom,
			&trans.Person,
			&trans.Repeat,
		)
	case transTypeDebt:
		s.debtMutex.Lock()
		row = s.db.QueryRowContext(ctx, storage.GetOneDebt, id)
		s.debtMutex.Unlock()
		err = row.Scan(
			&trans.DebtType,
			&trans.IdBillFrom,
			&trans.IdBillTo,
			&trans.Person,
		)
	case transTypeTransfer:
		s.transferMutex.Lock()
		row = s.db.QueryRowContext(ctx, storage.GetOneTransfer, id)
		s.transferMutex.Unlock()
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

	s.transMutex.Lock()
	rows, err := s.db.QueryContext(ctx, storage.GetSomeTransactions, count)
	s.transMutex.Lock()
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
