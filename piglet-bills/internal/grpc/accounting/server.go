package accounting

import (
	"context"
	"google.golang.org/grpc"
	billsv1 "piglet-bills-service/api/proto/gen"
)

type serverAPI struct {
	billsv1.UnimplementedPigletBillsServer
	//accounting Accounting
}

//type Accounting interface {
//	createBill(ctx context.Context,
//		billType bool,
//		billStatus bool,
//		billName string,
//		currentSum float32,
//		date time.Time,
//		monthlyPayment float32,
//	) (bill models.Bill, err error)
//	// TODO: несколько счетов, getSomeBills получает массив с id
//	getSomeBills(ctx context.Context) (bills []models.Bill, err error)
//	getBill(ctx context.Context, ID string) (bill models.Bill, err error)
//	updateBill(ctx context.Context,
//		billType bool,
//		billStatus bool,
//		billName string,
//		currentSum float32,
//		date time.Time,
//		monthlyPayment float32,
//	) (bill models.Bill, err error)
//	deleteBill(ctx context.Context, ID string) (success bool, err error)
//}

func Register(gRPCServer *grpc.Server) {
	billsv1.RegisterPigletBillsServer(gRPCServer, &serverAPI{})
}

func (s *serverAPI) CreateBill(
	ctx context.Context,
	req *billsv1.CreateBillRequest,
) (*billsv1.CreateBillResponse, error) {
	// TODO: setup validation

	// TODO: setup logic

	return &billsv1.CreateBillResponse{
		Bill: &billsv1.Bill{
			Id:             "",
			BillType:       false,
			BillStatus:     false,
			BillName:       "",
			CurrentSum:     0,
			Date:           nil,
			MonthlyPayment: 0,
		},
	}, nil
}

func (s *serverAPI) GetSomeBills(
	ctx context.Context,
	req *billsv1.GetSomeBillsRequest,
) (*billsv1.GetSomeBillsResponse, error) {
	// TODO: setup validation

	// TODO: setup logic
	var bills []*billsv1.Bill

	return &billsv1.GetSomeBillsResponse{
		Bills: bills,
	}, nil
}

func (s *serverAPI) GetBill(
	ctx context.Context,
	req *billsv1.GetBillRequest,
) (*billsv1.GetBillResponse, error) {
	// TODO: setup validation

	// TODO: setup logic

	return &billsv1.GetBillResponse{
		Bill: &billsv1.Bill{
			Id:             "",
			BillType:       false,
			BillStatus:     false,
			BillName:       "",
			CurrentSum:     0,
			Date:           nil,
			MonthlyPayment: 0,
		},
	}, nil
}

func (s *serverAPI) UpdateBill(
	ctx context.Context,
	req *billsv1.UpdateBillRequest,
) (*billsv1.UpdateBillResponse, error) {
	// TODO: setup validation

	// TODO: setup logic

	return &billsv1.UpdateBillResponse{
		Bill: &billsv1.Bill{
			Id:             "",
			BillType:       false,
			BillStatus:     false,
			BillName:       "",
			CurrentSum:     0,
			Date:           nil,
			MonthlyPayment: 0,
		},
	}, nil
}

func (s *serverAPI) DeleteBill(
	ctx context.Context,
	req *billsv1.DeleteBillRequest,
) (*billsv1.DeleteBillResponse, error) {
	// TODO: setup validation

	// TODO: setup logic

	return &billsv1.DeleteBillResponse{
		Success: true,
	}, nil
}
