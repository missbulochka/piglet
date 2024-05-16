package transactions

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"piglet-transactions-service/internal/domain/models"
)

const (
	debtTypeImCreditor  = true
	debtTypeImDebtor    = false
	transactionsCount   = 20
	categoryTypeExpense = true
	categoryTypeIncome  = false
	transTypeIncome     = 1
	transTypeExpense    = 2
	transTypeDebt       = 3
	transTypeTransfer   = 4
)

// CreateTransaction create new transaction in the system and returns it
// If bill or category with given names don't exist, returns error
func (t *Transactions) CreateTransaction(
	ctx context.Context,
	trans *models.Transaction,
) (err error) {
	const op = "pigletTransactions | transactions.CreateTransaction"
	log := t.log.With(slog.String("op", op))

	trans.Id = uuid.New()

	log.Info("verifying bill")
	if trans.TransType == transTypeIncome ||
		(trans.TransType == transTypeDebt && trans.DebtType == debtTypeImDebtor) ||
		trans.TransType == transTypeTransfer {
		err = verifyBill(ctx, trans.IdBillTo, t)
		if err != nil {
			log.Error("%w", err)

			return fmt.Errorf("%s: %w", op, err)
		}
	}
	if trans.TransType == transTypeExpense ||
		(trans.TransType == transTypeDebt && trans.DebtType == debtTypeImCreditor) ||
		trans.TransType == transTypeTransfer {
		err = verifyBill(ctx, trans.IdBillFrom, t)
		if err != nil {
			log.Error("%w", err)

			return fmt.Errorf("%s: %w", op, err)
		}
	}

	// HACK: может быть имеет смысл вынести в отдельную функцию
	if trans.TransType == transTypeIncome || trans.TransType == transTypeExpense {
		log.Info("verifying category")
		cat, err := t.categoryProvider.GetCategory(ctx, trans.IdCategory)
		if err != nil {
			log.Error("failed to verify category", err)

			return fmt.Errorf("%s: %w", op, err)
		}

		if !((trans.TransType == transTypeIncome && cat.CategoryType == categoryTypeIncome) ||
			(trans.TransType == transTypeExpense && cat.CategoryType == categoryTypeExpense)) {
			log.Error("failed to verify category: inconsistency transaction and category types")

			return fmt.Errorf("%s: inconsistency transaction and category types", op)
		}

	}

	log.Info("saving transaction")

	if err = t.transSaver.SaveTransaction(ctx, *trans); err != nil {
		log.Error("failed to save transaction", err)

		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("transaction saved")

	return nil
}

// UpdateTransaction update exist transaction in the system and returns it
// If transaction with given id doesn't exist, returns error
func (t *Transactions) UpdateTransaction(
	ctx context.Context,
	trans *models.Transaction,
) (dif decimal.Decimal, err error) {
	const op = "pigletTransactions | transactions.UpdateTransaction"
	log := t.log.With(slog.String("op", op))

	log.Info("verifying bill")
	if trans.TransType == transTypeIncome ||
		(trans.TransType == transTypeDebt && trans.DebtType == debtTypeImDebtor) ||
		trans.TransType == transTypeTransfer {
		err = verifyBill(ctx, trans.IdBillTo, t)
		if err != nil {
			log.Error("%w", err)

			return decimal.Decimal{}, fmt.Errorf("%s: %w", op, err)
		}
	}
	if trans.TransType == transTypeExpense ||
		(trans.TransType == transTypeDebt && trans.DebtType == debtTypeImCreditor) ||
		trans.TransType == transTypeTransfer {
		err = verifyBill(ctx, trans.IdBillFrom, t)
		if err != nil {
			log.Error("%w", err)

			return decimal.Decimal{}, fmt.Errorf("%s: %w", op, err)
		}
	}

	// HACK: может быть имеет смысл вынести в отдельную функцию
	if trans.TransType == transTypeIncome || trans.TransType == transTypeExpense {
		log.Info("verifying category")
		cat, err := t.categoryProvider.GetCategory(ctx, trans.IdCategory)
		if err != nil {
			log.Error("failed to verify category", err)

			return decimal.Decimal{}, fmt.Errorf("%s: %w", op, err)
		}

		if !((trans.TransType == transTypeIncome && cat.CategoryType == categoryTypeIncome) ||
			(trans.TransType == transTypeExpense && cat.CategoryType == categoryTypeExpense)) {
			log.Error("failed to verify category: inconsistency transaction and category types")

			return decimal.Decimal{}, fmt.Errorf("%s: inconsistency transaction and category types", op)
		}

	}

	var oldtrans models.Transaction
	err = t.transProvider.GetTransaction(ctx, trans.Id, &oldtrans)
	if err != nil {
		log.Error("failed to find transaction", err)

		return decimal.Decimal{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("updating category")

	if err = t.transSaver.UpdateTransaction(ctx, *trans); err != nil {
		log.Error("failed to update transaction", err)

		return decimal.Decimal{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("category updated")

	return trans.Sum.Sub(oldtrans.Sum), nil
}

// DeleteTransaction delete transaction in the system
// If transaction with given id doesn't exist, returns error
func (t *Transactions) DeleteTransaction(ctx context.Context, id uuid.UUID) error {
	const op = "pigletTransactions | transactions.DeleteTransaction"
	log := t.log.With(slog.String("op", op))

	_, transType, _, _, err := t.transProvider.DefaultTransInfo(ctx, id)
	if err != nil {
		log.Error("transaction doesn't exist", err)

		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("deleting transaction")

	if err = t.transSaver.DeleteTransaction(ctx, id, transType); err != nil {
		log.Error("failed to delete transaction", err)

		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("transaction deleted")

	return nil
}

// GetTransaction search transaction in the system
// If transaction with given id doesn't exist, returns error
func (t *Transactions) GetTransaction(ctx context.Context, id uuid.UUID) (trans models.Transaction, err error) {
	const op = "pigletTransactions | transactions.GetTransaction"
	log := t.log.With(slog.String("op", op))

	trans.Date, trans.TransType, trans.Sum, trans.Comment, err = t.transProvider.DefaultTransInfo(ctx, id)
	if err != nil {
		log.Error("failed to search transaction", err)

		return trans, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("receiving transaction")

	if err = t.transProvider.GetTransaction(ctx, id, &trans); err != nil {
		log.Error("failed to get transaction", err)

		return trans, fmt.Errorf("%s: %w", op, err)
	}

	trans.Id = id

	log.Info("transaction received")

	return trans, nil
}

// GetLast20Transactions search last 20 transactions in the system
// If something go wrong, returns error
func (t *Transactions) GetLast20Transactions(ctx context.Context) (trans []*models.Transaction, err error) {
	const op = "pigletTransactions | transactions.GetLast20Transactions"
	log := t.log.With(slog.String("op", op))

	log.Info("receiving transactions")

	if err = t.transProvider.GetSomeTransactions(ctx, &trans, transactionsCount); err != nil {
		log.Error("failed to get transactions", err)

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	for i := 0; i < len(trans); i++ {
		if err = t.transProvider.GetTransaction(ctx, trans[i].Id, trans[i]); err != nil {
			log.Error("failed to get transaction", err)

			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	log.Info("transactions received")

	return trans, nil
}
