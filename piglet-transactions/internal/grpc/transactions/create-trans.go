package transactionsgrpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	transv1 "piglet-transactions-service/api/proto/gen"
	validation "piglet-transactions-service/internal/domain/validator"
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

	if trans, err = s.transactions.CreateTransaction(ctx, trans); err != nil {
		// TODO: проверка ошибки о несуществовании счета или категории

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
