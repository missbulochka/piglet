package transactionsgrpc

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc"

	billsv1 "github.com/missbulochka/protos/gen/piglet-bills"
	transv1 "github.com/missbulochka/protos/gen/piglet-transactions"
	"piglet-transactions-service/internal/domain/models"
)

type serverAPI struct {
	transv1.UnimplementedPigletTransactionsServer
	billsCli     billsv1.PigletBillsClient
	transactions Transactions
	categories   Categories
}

type Transactions interface {
	CreateTransaction(ctx context.Context, trans *models.Transaction) (err error)
	DeleteTransaction(ctx context.Context, id uuid.UUID) (err error)
	GetTransaction(ctx context.Context, id uuid.UUID) (trans models.Transaction, err error)
	GetLast20Transactions(ctx context.Context) (trans []*models.Transaction, err error)
}

type Categories interface {
	CreateCategory(ctx context.Context, cat *models.Category) (err error)
	UpdateCategory(ctx context.Context, cat *models.Category) (err error)
	DeleteCategory(ctx context.Context, id uuid.UUID) (err error)
	GetCategory(ctx context.Context, id uuid.UUID) (cat models.Category, err error)
	GetAllCategories(ctx context.Context) (cat []*models.Category, err error)
}

func Register(gRPCServer *grpc.Server, conn *grpc.ClientConn, transactions Transactions) {
	transv1.RegisterPigletTransactionsServer(
		gRPCServer,
		&serverAPI{
			transactions: transactions,
			billsCli:     billsv1.NewPigletBillsClient(conn),
		})
}
