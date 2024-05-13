package transactionsgrpc

import (
	"context"

	transv1 "piglet-transactions-service/api/proto/gen"
)

func (s *serverAPI) CreateTransaction(
	ctx context.Context,
	req *transv1.CreateTransactionRequest,
) (*transv1.TransactionResponse, error) {
	// TODO: валидация данных

	// TODO: сервисный слой
	panic("implement me")
}
