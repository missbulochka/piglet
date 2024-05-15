package transactionsgrpc

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	transv1 "github.com/missbulochka/protos/gen/piglet-transactions"
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
		"",
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
		return nil, status.Errorf(codes.InvalidArgument, "invalid creditals")
	}

	if err = VerifyBills(ctx, s.billsCli, &trans); err != nil {
		return nil, err
	}
	fmt.Println("bill verified")

	if err = s.transactions.CreateTransaction(ctx, &trans); err != nil {
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

func (s *serverAPI) UpdateTransaction(
	ctx context.Context,
	req *transv1.Transaction,
) (*transv1.TransactionResponse, error) {
	trans, err := validation.TransValidator(
		req.GetId(),
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
		return nil, status.Errorf(codes.InvalidArgument, "invalid creditals")
	}

	if err = VerifyBills(ctx, s.billsCli, &trans); err != nil {
		return nil, err
	}
	fmt.Println("bill verified")

	if err = s.transactions.UpdateTransaction(ctx, &trans); err != nil {
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

func (s *serverAPI) DeleteTransaction(
	ctx context.Context,
	req *transv1.IdRequest,
) (*transv1.SuccessResponse, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid creditals")
	}

	trans, err := s.transactions.GetTransaction(ctx, id)
	if err != nil {
		// TODO: проверка ошибки о несуществовании транзакции

		return &transv1.SuccessResponse{Success: false}, status.Errorf(codes.Internal, "internal error")
	}

	if err = s.transactions.DeleteTransaction(ctx, id); err != nil {
		return &transv1.SuccessResponse{Success: false}, status.Errorf(codes.Internal, "internal error")
	}

	// reverse the sum
	trans.Sum = trans.Sum.Neg()
	BillFixer(
		trans.IdBillTo.String(),
		trans.IdBillFrom.String(),
		trans.TransType,
		trans.DebtType,
		trans.Sum,
		s.billsCli,
	)

	return &transv1.SuccessResponse{Success: true}, nil
}

func (s *serverAPI) GetTransaction(
	ctx context.Context,
	req *transv1.IdRequest,
) (*transv1.TransactionResponse, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid creditals")
	}

	trans, err := s.transactions.GetTransaction(ctx, id)
	if err != nil {
		// TODO: проверка ошибки о несуществовании транзакции

		return nil, status.Errorf(codes.Internal, "internal error")
	}

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

func (s *serverAPI) GetAllTransactions(
	ctx context.Context,
	req *transv1.EmptyRequest,
) (*transv1.GetAllTransactionsResponse, error) {
	trans, err := s.transactions.GetLast20Transactions(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	resTrans := ReturnTransactions(trans)

	return &transv1.GetAllTransactionsResponse{Transactions: resTrans}, nil
}
