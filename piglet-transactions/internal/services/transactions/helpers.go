package transactions

import (
	"context"
	"fmt"
	"github.com/google/uuid"
)

func verifyBill(ctx context.Context, id uuid.UUID, t *Transactions) (err error) {
	billStatus, err := t.billUpdater.GetBill(ctx, id)
	if err != nil {
		return fmt.Errorf("bill doesn't exist")
	}
	if !billStatus {

		return fmt.Errorf("bill is closed")
	}
	return nil
}
