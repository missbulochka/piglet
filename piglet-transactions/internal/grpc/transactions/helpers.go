package transactionsgrpc

import (
	"context"
	"fmt"

	"github.com/shopspring/decimal"
	"google.golang.org/protobuf/types/known/timestamppb"

	billsv1 "github.com/missbulochka/protos/gen/piglet-bills"
	transv1 "github.com/missbulochka/protos/gen/piglet-transactions"
	"piglet-transactions-service/internal/domain/models"
)

const (
	transTypeIncome   = 1
	transTypeExpense  = 2
	transTypeDebt     = 3
	transTypeTransfer = 4
)

func BillFixer(
	idTo string,
	idFrom string,
	transType uint8,
	debtType bool,
	sum decimal.Decimal,
	cli billsv1.PigletBillsClient,
) {
	ctx := context.Background()
	var id string
	var floatSum float32

	if transType == transTypeIncome ||
		(transType == transTypeDebt && debtType == debtTypeImDebtor) ||
		transType == transTypeTransfer {
		// HACK: обработка ошибок
		float64Sum, _ := sum.Float64()
		floatSum = float32(float64Sum)
		id = idTo
	}

	if transType == transTypeExpense ||
		(transType == transTypeDebt && debtType == debtTypeImCreditor) ||
		transType == transTypeTransfer {
		// HACK: обработка ошибок
		float64Sum, _ := sum.Neg().Float64()
		floatSum = float32(float64Sum)
		id = idFrom
	}

	go func(ctx context.Context, sum float32, id string, cli billsv1.PigletBillsClient) {
		_, err := cli.FixBillSum(
			ctx,
			&billsv1.FixBillSumRequest{
				Id:  id,
				Sum: sum,
			},
		)
		if err != nil {
			fmt.Println("service synchronization error: %w", err)
		}
	}(ctx, floatSum, id, cli)
}

func VerifyBills(
	ctx context.Context,
	cli billsv1.PigletBillsClient,
	trans *models.Transaction,
) (err error) {
	existTo, err := cli.VerifyBill(
		ctx,
		&billsv1.IdRequest{
			Id: trans.IdBillTo.String(),
		},
	)
	if err != nil {
		return err
	}

	existFrom, err := cli.VerifyBill(
		ctx,
		&billsv1.IdRequest{
			Id: trans.IdBillFrom.String(),
		},
	)
	if err != nil {
		return err
	}

	switch trans.TransType {
	case transTypeIncome:
		if existTo.Success == true {
			return nil
		}
	case transTypeExpense:
		if existFrom.Success == true {
			return nil
		}
	case transTypeDebt:
		if existFrom.Success == true || existTo.Success == true {
			return nil
		}
	case transTypeTransfer:
		if existTo.Success == true && existFrom.Success == true {
			return nil
		}
	}

	return fmt.Errorf("bill varification failed")
}

func ReturnTransactions(trans []*models.Transaction) (resTrans []*transv1.Transaction) {
	for _, tr := range trans {
		// HACK: обработка ошибок
		sumFoProto, _ := tr.Sum.Float64()

		node := &transv1.Transaction{
			Id:         tr.Id.String(),
			Date:       timestamppb.New(tr.Date),
			TransType:  int32(tr.TransType),
			Sum:        float32(sumFoProto),
			Comment:    tr.Comment,
			IdCategory: tr.Comment,
			DebtType:   tr.DebtType,
			IdBillTo:   tr.IdBillTo.String(),
			IdBillFrom: tr.IdBillFrom.String(),
			Person:     tr.Person,
			Repeat:     tr.Repeat,
		}
		resTrans = append(resTrans, node)
	}

	return resTrans
}

func ReturnCategories(cat []*models.Category) (resCat []*transv1.Category) {
	for _, c := range cat {
		// HACK: обработка ошибок

		node := &transv1.Category{
			Id:           c.Id.String(),
			Type:         c.CategoryType,
			CategoryName: c.Name,
			Mandatory:    c.Mandatory,
		}
		resCat = append(resCat, node)
	}

	return resCat
}
