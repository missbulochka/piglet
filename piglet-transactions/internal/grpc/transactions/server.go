package transactionsgrpc

import (
	"context"
	"google.golang.org/grpc"
	"piglet-transactions-service/internal/domain/models"

	transv1 "piglet-transactions-service/api/proto/gen"
)

type serverAPI struct {
	transv1.UnimplementedPigletTransactionsServer
	transactions Transactions
}

type Transactions interface {
	CreateTransaction(ctx context.Context, trans models.Transaction) (savedTrans models.Transaction, err error)
}

func Register(gRPCServer *grpc.Server, transactions Transactions) {
	transv1.RegisterPigletTransactionsServer(gRPCServer, &serverAPI{transactions: transactions})
}
