package accountingrpc

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	billsv1 "github.com/missbulochka/protos/gen/piglet-bills"
	models "piglet-bills-service/internal/domain/model"
	"piglet-bills-service/internal/services/accounting"
)

type serverAPI struct {
	billsv1.UnimplementedPigletBillsServer
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

func Register(gRPCServer *grpc.Server, accounting Accounting) {
	billsv1.RegisterPigletBillsServer(gRPCServer, &serverAPI{accounting: accounting})
}

const (
	accountType = true
	goalType    = false
)

// HACK: стоит разбить на CreateAccount и CreateGoal
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

	// HACK: обработка ошибок
	goalSum, _ := decimal.NewFromString(strconv.FormatUint(uint64(req.GetGoalSum()), 10))

	bill, err := s.accounting.CreateBill(
		ctx,
		req.GetBillType(),
		req.GetBillName(),
		goalSum,
		req.GetDate().AsTime(),
	)
	if err != nil {
		if errors.Is(err, accounting.ErrBillExists) {
			return nil, status.Error(codes.InvalidArgument, "invalid credits")
		}

		return nil, status.Errorf(codes.Internal, "internal error")
	}

	// HACK: поработать над преобразованиями и обработкой ошибок
	currentSum, _ := bill.CurrentSum.Float64()
	newGoalSum, _ := bill.GoalSum.Float64()
	monthlyPayment := uint32(int32(bill.MonthlyPayment.IntPart()))

	return &billsv1.BillResponse{
		Bill: &billsv1.Bill{
			Id:             bill.ID,
			BillType:       bill.BillType,
			BillStatus:     bill.BillStatus,
			BillName:       bill.Name,
			CurrentSum:     float32(currentSum),
			GoalSum:        float32(newGoalSum),
			Date:           timestamppb.New(bill.Date),
			MonthlyPayment: monthlyPayment,
		},
	}, nil
}

// HACK: возврат нескольких счетов, getSomeBills получает массив с id (?)
func (s *serverAPI) GetAllAccounts(
	ctx context.Context,
	req *billsv1.GetSomeBillsRequest,
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
	req *billsv1.GetSomeBillsRequest,
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
	monthlyPayment := uint32(int32(bill.MonthlyPayment.IntPart()))
	return &billsv1.BillResponse{
		Bill: &billsv1.Bill{
			Id:             bill.ID,
			BillType:       bill.BillType,
			BillStatus:     bill.BillStatus,
			BillName:       bill.Name,
			CurrentSum:     float32(currentSum),
			GoalSum:        float32(newGoalSum),
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

	// HACK: обработка ошибок
	currentSum, _ := decimal.NewFromString(strconv.FormatUint(uint64(req.GetCurrentSum()), 10))
	goalSum, _ := decimal.NewFromString(strconv.FormatUint(uint64(req.GetGoalSum()), 10))

	bill, err := s.accounting.UpdateBill(
		ctx,
		req.GetId(),
		req.GetBillName(),
		currentSum,
		req.GetBillStatus(),
		goalSum,
		req.GetDate().AsTime(),
	)
	if err != nil {
		if errors.Is(err, accounting.ErrBillNotFound) {
			return nil, status.Error(codes.InvalidArgument, "invalid credits")
		}

		return nil, status.Errorf(codes.Internal, "internal error")
	}

	// HACK: поработать над преобразованиями и обработкой ошибок
	newCurrentSum, _ := bill.CurrentSum.Float64()
	newGoalSum, _ := bill.GoalSum.Float64()
	monthlyPayment := uint32(int32(bill.MonthlyPayment.IntPart()))

	return &billsv1.BillResponse{
		Bill: &billsv1.Bill{
			Id:             bill.ID,
			BillType:       bill.BillType,
			BillStatus:     bill.BillStatus,
			BillName:       bill.Name,
			CurrentSum:     float32(newCurrentSum),
			GoalSum:        float32(newGoalSum),
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

	success, err := s.accounting.DeleteBill(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, accounting.ErrBillNotFound) {
			return nil, status.Error(codes.InvalidArgument, "invalid uuid")
		}
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	return &billsv1.SuccessResponse{
		Success: success,
	}, nil
}

func (s *serverAPI) VerifyBill(
	ctx context.Context,
	req *billsv1.IdRequest,
) (*billsv1.SuccessResponse, error) {
	// HACK: улучшить валидацию
	if err := orValidation(req.GetId(), ""); err != nil {
		return nil, err
	}

	success, err := s.accounting.VerifyBill(ctx, req.GetId())
	if err != nil {
		success = false
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

	sum, _ := decimal.NewFromString(strconv.FormatUint(uint64(req.GetSum()), 10))

	err := s.accounting.FixBillSum(ctx, req.GetId(), sum)
	if err != nil {
		return &emptypb.Empty{}, status.Errorf(codes.Internal, "internal error")
	}

	return &emptypb.Empty{}, nil
}

func validation(
	vd validateData,
) error {
	val := validator.New(validator.WithRequiredStructEnabled())

	if err := val.Struct(vd); err != nil {
		var validationErr validator.ValidationErrors
		if errors.As(err, &validationErr) {
			log.Println("Validation errors:")
			for _, err := range validationErr {
				log.Printf("- Namespace: %s, Field: %s, Tag: %s, ActualTag: %s, Value: %v, Param: %s",
					err.Namespace(), err.Field(), err.Tag(), err.ActualTag(), err.Value(), err.Param())
			}
			return status.Errorf(codes.InvalidArgument, "invalid bill: %v", validationErr)
		}
		log.Printf("Validation error: %v", err)
		return status.Errorf(codes.Internal, "internal error: %v", err)
	}
	return nil
}

func orValidation(uuid string, name string) error {
	val := validator.New()

	if err := val.Var(uuid, "required"); err != nil {
		if err2 := val.Var(name, "required"); err2 != nil {
			return status.Errorf(codes.InvalidArgument, err2.Error())
		}
	}

	return nil
}

type validateData struct {
	billType bool   `validate:"boolean"`
	billName string `validate:"required"`
}

func billsConversion(bills []*models.Bill) []*billsv1.Bill {
	var resBills []*billsv1.Bill
	for _, bill := range bills {
		currentSum, _ := bill.CurrentSum.Float64()
		goalSum, _ := bill.GoalSum.Float64()
		monthlyPayment := uint32(int32(bill.MonthlyPayment.IntPart()))

		resBill := &billsv1.Bill{
			Id:             bill.ID,
			BillType:       bill.BillType,
			BillStatus:     bill.BillStatus,
			BillName:       bill.Name,
			CurrentSum:     float32(currentSum),
			GoalSum:        float32(goalSum),
			Date:           timestamppb.New(bill.Date),
			MonthlyPayment: monthlyPayment,
		}
		resBills = append(resBills, resBill)
	}

	return resBills
}
