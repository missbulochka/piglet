package transactionsgrpc

import (
	"context"
	"google.golang.org/grpc"
	"piglet-transactions-service/internal/domain/models"

	billsv1 "github.com/missbulochka/protos/gen/piglet-bills"
	transv1 "github.com/missbulochka/protos/gen/piglet-transactions"
)

type serverAPI struct {
	transv1.UnimplementedPigletTransactionsServer
	billsCli     billsv1.PigletBillsClient
	transactions Transactions
}

type Transactions interface {
	CreateTransaction(ctx context.Context, trans *models.Transaction) (err error)
}

func Register(gRPCServer *grpc.Server, conn *grpc.ClientConn, transactions Transactions) {
	transv1.RegisterPigletTransactionsServer(
		gRPCServer,
		&serverAPI{
			transactions: transactions,
			billsCli:     billsv1.NewPigletBillsClient(conn),
		})
}
