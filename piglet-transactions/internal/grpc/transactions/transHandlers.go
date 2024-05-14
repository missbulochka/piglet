package transactionsgrpc

import (
	"context"
	"fmt"

	"github.com/shopspring/decimal"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	billsv1 "github.com/missbulochka/protos/gen/piglet-bills"
	transv1 "github.com/missbulochka/protos/gen/piglet-transactions"
	"piglet-transactions-service/internal/domain/models"
	validation "piglet-transactions-service/internal/domain/validator"
)

const (
	debtTypeImCreditor = true
	debtTypeImDebtor   = false
)

func (s *serverAPI) CreateTransaction(
	ctx context.Context,
	req *transv1.CreateTransactionRequest,
) (*transv1.TransactionResponse, error) {
	trans, err := validation.TransValidator(
		req.GetDate(),
		req.GetTransType(),
		req.GetSum(),
		req.GetComment(),
		req.GetIdCategory(),
		req.GetDebtType(),
		req.GetIdBillTo(),
		req.GetIdBillFrom(),
		req.GetPerson(),
		req.GetRepeat(),
	)
	if err != nil {
		return nil, err
	}

	if err = VerifyBills(ctx, s.billsCli, &trans); err != nil {
		return nil, err
	}
	fmt.Println("bill verified")

	if trans, err = s.transactions.CreateTransaction(ctx, trans); err != nil {
		// TODO: проверка ошибки о несуществовании счета или категории

		return nil, status.Errorf(codes.Internal, "internal error")
	}

	BillFixer(
		trans.IdBillTo.String(),
		trans.IdBillFrom.String(),
		trans.TransType,
		trans.DebtType,
		trans.Sum,
		s.billsCli,
	)

	// HACK: обработка ошибки
	sumFoProto, _ := trans.Sum.Float64()

	return &transv1.TransactionResponse{
		Transaction: &transv1.Transaction{
			Id:         trans.Id.String(),
			Date:       timestamppb.New(trans.Date),
			TransType:  int32(trans.TransType),
			Sum:        float32(sumFoProto),
			Comment:    trans.Comment,
			IdCategory: trans.Comment,
			DebtType:   trans.DebtType,
			IdBillTo:   trans.IdBillTo.String(),
			IdBillFrom: trans.IdBillFrom.String(),
			Person:     trans.Person,
			Repeat:     trans.Repeat,
		},
	}, nil
}

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

	if transType == 1 || (transType == 3 && debtType == debtTypeImDebtor) || transType == 4 {
		// HACK: обработка ошибок
		float64Sum, _ := sum.Float64()
		floatSum = float32(float64Sum)
		id = idTo
	}

	if transType == 2 || (transType == 4 && debtType == debtTypeImCreditor) || transType == 4 {
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
	case 1:
		if existTo.Success == true {
			return nil
		}
	case 2:
		if existFrom.Success == true {
			return nil
		}
	case 3:
		if existFrom.Success == true || existTo.Success == true {
			return nil
		}
	case 4:
		if existTo.Success == true && existFrom.Success == true {
			return nil
		}
	}

	return fmt.Errorf("bill varification failed")
}
