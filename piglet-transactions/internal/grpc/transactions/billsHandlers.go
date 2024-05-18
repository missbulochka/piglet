package transactionsgrpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	transv1 "github.com/missbulochka/protos/gen/piglet-transactions"
	validation "piglet-transactions-service/internal/domain/validator"
)

func (s *serverAPI) UpdateBill(ctx context.Context, req *transv1.Bill) (*emptypb.Empty, error) {
	id, err := validation.BillValidator(req.GetId())
	if err != nil {
		fmt.Println("invalid uuid")
	}

	err = s.bills.UpdateBills(ctx, id, req.GetBillStatus(), req.GetDeletion())
	if err != nil {
		return &emptypb.Empty{}, status.Errorf(codes.Internal, "internal error")
	}

	return &emptypb.Empty{}, nil
}
