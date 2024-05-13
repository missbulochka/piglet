package transactionsgrpc

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"google.golang.org/grpc"

	transv1 "piglet-transactions-service/api/proto/gen"
)

type serverAPI struct {
	transv1.UnimplementedPigletTransactionsServer
	transactions Transactions
}

type Transactions interface {
	CreateTransaction(
		ctx context.Context,
		date time.Time,
		transType uint8,
		sum decimal.Decimal,
		comment string,
		idCategory uuid.UUID,
		debtType bool,
		idBillTo uuid.UUID,
		idBillFrom uuid.UUID,
		person string,
		repeat bool,
	) (err error)
}

func Register(gRPCServer *grpc.Server, transactions Transactions) {
	transv1.RegisterPigletTransactionsServer(gRPCServer, &serverAPI{transactions: transactions})
}
