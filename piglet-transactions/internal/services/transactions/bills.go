package transactions

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"log/slog"
)

func (t *Transactions) UpdateBills(
	ctx context.Context,
	id uuid.UUID,
	billStatus bool,
	del bool,
) (err error) {
	const op = "pigletTransactions | transactions.UpdateBills"
	log := t.log.With(slog.String("op", op))

	if _, err = t.billUpdater.GetBill(ctx, id); err != nil {
		log.Info("adding bill")
		if err = t.billUpdater.SaveBill(ctx, id, billStatus); err != nil {
			log.Error("failed to save bill", err)

			return fmt.Errorf("%s: %w", op, err)
		}
		log.Info("bill added")
	} else {
		if del {
			log.Info("deleting bill")
			if err = t.billUpdater.DeleteBill(ctx, id); err != nil {
				log.Error("failed to delete bill", err)

				return fmt.Errorf("%s: %w", op, err)
			}
			log.Info("bill deleted")
		} else {
			log.Info("updating bill")
			if err = t.billUpdater.UpdateBill(ctx, id, billStatus); err != nil {
				log.Error("failed to update bill", err)

				return fmt.Errorf("%s: %w", op, err)
			}
			log.Info("bill updated")
		}
	}

	return nil
}
