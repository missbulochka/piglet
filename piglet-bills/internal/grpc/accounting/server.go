package accountingrpc

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	billsv1 "github.com/missbulochka/protos/gen/piglet-bills"
	transv1 "github.com/missbulochka/protos/gen/piglet-transactions"
	models "piglet-bills-service/internal/domain/model"
	"piglet-bills-service/internal/services/accounting"
)

type serverAPI struct {
	billsv1.UnimplementedPigletBillsServer
	transCli   transv1.PigletTransactionsClient
	accounting Accounting
}

type Accounting interface {
	CreateBill(
		ctx context.Context,
		billType bool,
		billName string,
		goalSum decimal.Decimal,
		date time.Time,
	) (bill models.Bill, err error)
	GetSomeBills(ctx context.Context, billType bool) (bills []*models.Bill, err error)
	GetBill(
		ctx context.Context,
		billId string,
	) (bill models.Bill, err error)
	UpdateBill(
		ctx context.Context,
		id string,
		billName string,
		currentSum decimal.Decimal,
		billStatus bool,
		goalSum decimal.Decimal,
		date time.Time,
	) (bill models.Bill, err error)
	DeleteBill(ctx context.Context, id string) (success bool, err error)
	VerifyBill(ctx context.Context, id string) (success bool, err error)
	FixBillSum(ctx context.Context, id string, sum decimal.Decimal) (err error)
}

func Register(gRPCServer *grpc.Server, conn *grpc.ClientConn, accounting Accounting) {
	billsv1.RegisterPigletBillsServer(
		gRPCServer,
		&serverAPI{
			accounting: accounting,
			transCli:   transv1.NewPigletTransactionsClient(conn),
		},
	)
}

const (
	accountType = true
	goalType    = false
)

func (s *serverAPI) CreateBill(
	ctx context.Context,
	req *billsv1.CreateBillRequest,
) (*billsv1.BillResponse, error) {
	// HACK: улучшить валидацию (не передавать структуру в целом)
	// HACK: проверка даты для цели на корректность (не в прошлом и реальная в будущем)
	if err := validation(
		validateData{
			billType: req.GetBillType(),
			billName: req.GetBillName(),
		},
	); err != nil {
		return nil, err
	}

	bill, err := s.accounting.CreateBill(
		ctx,
		req.GetBillType(),
		req.GetBillName(),
		decimal.NewFromFloat(req.GetGoalSum()),
		req.GetDate().AsTime(),
	)
	if err != nil {
		if errors.Is(err, accounting.ErrBillExists) {
			return nil, status.Error(codes.InvalidArgument, "invalid credits")
		}

		return nil, status.Errorf(codes.Internal, "internal error")
	}

	syncBills(bill.ID, s.transCli, bill.BillStatus, false)

	// HACK: поработать над преобразованиями и обработкой ошибок
	currentSum, _ := bill.CurrentSum.Float64()
	newGoalSum, _ := bill.GoalSum.Float64()
	monthlyPayment := uint32(bill.MonthlyPayment.IntPart())

	return &billsv1.BillResponse{
		Bill: &billsv1.Bill{
			Id:             bill.ID,
			BillType:       bill.BillType,
			BillStatus:     bill.BillStatus,
			BillName:       bill.Name,
			CurrentSum:     currentSum,
			GoalSum:        newGoalSum,
			Date:           timestamppb.New(bill.Date),
			MonthlyPayment: monthlyPayment,
		},
	}, nil
}

// HACK: возврат нескольких счетов, getSomeBills получает массив с id (?)
func (s *serverAPI) GetAllAccounts(
	ctx context.Context,
	req *emptypb.Empty,
) (*billsv1.GetSomeBillsResponse, error) {
	bills, err := s.accounting.GetSomeBills(ctx, accountType)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "something go wrong")
	}

	// HACK: аккуратный возврат всех счетов
	resBills := billsConversion(bills)

	return &billsv1.GetSomeBillsResponse{
		Bills: resBills,
	}, nil
}

func (s *serverAPI) GetAllGoals(
	ctx context.Context,
	req *emptypb.Empty,
) (*billsv1.GetSomeBillsResponse, error) {
	bills, err := s.accounting.GetSomeBills(ctx, goalType)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "something go wrong")
	}

	// HACK: аккуратный возврат всех счетов
	resBills := billsConversion(bills)

	return &billsv1.GetSomeBillsResponse{
		Bills: resBills,
	}, nil
}

func (s *serverAPI) GetBill(
	ctx context.Context,
	req *billsv1.IdRequest,
) (*billsv1.BillResponse, error) {
	// HACK: улучшить валидацию
	if err := orValidation(req.GetId(), ""); err != nil {
		return nil, err
	}

	bill, err := s.accounting.GetBill(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, accounting.ErrBillNotFound) {
			return nil, status.Error(codes.InvalidArgument, "invalid uuid")
		}
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	// HACK: поработать над преобразованиями и обработкой ошибок
	currentSum, _ := bill.CurrentSum.Float64()
	newGoalSum, _ := bill.GoalSum.Float64()
	monthlyPayment := uint32(bill.MonthlyPayment.IntPart())
	return &billsv1.BillResponse{
		Bill: &billsv1.Bill{
			Id:             bill.ID,
			BillType:       bill.BillType,
			BillStatus:     bill.BillStatus,
			BillName:       bill.Name,
			CurrentSum:     currentSum,
			GoalSum:        newGoalSum,
			Date:           timestamppb.New(bill.Date),
			MonthlyPayment: monthlyPayment,
		},
	}, nil
}

func (s *serverAPI) UpdateBill(
	ctx context.Context,
	req *billsv1.UpdateBillRequest,
) (*billsv1.BillResponse, error) {
	// HACK: улучшить валидацию
	if err := orValidation(req.GetId(), req.GetBillName()); err != nil {
		return nil, err
	}

	bill, err := s.accounting.UpdateBill(
		ctx,
		req.GetId(),
		req.GetBillName(),
		decimal.NewFromFloat(req.GetCurrentSum()),
		req.GetBillStatus(),
		decimal.NewFromFloat(req.GetGoalSum()),
		req.GetDate().AsTime(),
	)
	if err != nil {
		if errors.Is(err, accounting.ErrBillNotFound) {
			return nil, status.Error(codes.InvalidArgument, "invalid credits")
		}

		return nil, status.Errorf(codes.Internal, "internal error")
	}

	// HACK: проверка на изменение статуса
	syncBills(bill.ID, s.transCli, bill.BillStatus, false)

	// HACK: поработать над преобразованиями и обработкой ошибок
	newCurrentSum, _ := bill.CurrentSum.Float64()
	newGoalSum, _ := bill.GoalSum.Float64()
	monthlyPayment := uint32(bill.MonthlyPayment.IntPart())

	return &billsv1.BillResponse{
		Bill: &billsv1.Bill{
			Id:             bill.ID,
			BillType:       bill.BillType,
			BillStatus:     bill.BillStatus,
			BillName:       bill.Name,
			CurrentSum:     newCurrentSum,
			GoalSum:        newGoalSum,
			Date:           timestamppb.New(bill.Date),
			MonthlyPayment: monthlyPayment,
		},
	}, nil
}

func (s *serverAPI) DeleteBill(
	ctx context.Context,
	req *billsv1.IdRequest,
) (*billsv1.SuccessResponse, error) {
	// HACK: улучшить валидацию
	if err := orValidation(req.GetId(), ""); err != nil {
		return nil, err
	}

	bill, err := s.accounting.GetBill(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, accounting.ErrBillNotFound) {
			return nil, status.Error(codes.InvalidArgument, "invalid uuid")
		}
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	success, err := s.accounting.DeleteBill(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, accounting.ErrBillNotFound) {
			return nil, status.Error(codes.InvalidArgument, "invalid uuid")
		}
		return nil, status.Errorf(codes.Internal, "internal error")
	} else {
		syncBills(bill.ID, s.transCli, bill.BillStatus, true)
	}

	return &billsv1.SuccessResponse{
		Success: success,
	}, nil
}

func (s *serverAPI) FixBillSum(
	ctx context.Context,
	req *billsv1.FixBillSumRequest,
) (*emptypb.Empty, error) {
	// HACK: улучшить валидацию
	if err := orValidation(req.GetId(), ""); err != nil {
		fmt.Println("invalid uuid")
	}

	sum := decimal.NewFromFloat(req.GetSum())

	err := s.accounting.FixBillSum(ctx, req.GetId(), sum)
	if err != nil {
		return &emptypb.Empty{}, status.Errorf(codes.Internal, "internal error")
	}

	return &emptypb.Empty{}, nil
}
