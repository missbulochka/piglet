package accountingrpc

import (
	"context"
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	billsv1 "piglet-bills-service/api/proto/gen"
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
	//GetSomeBills(ctx context.Context) (bills []models.Bill, err error)
	//GetBill(ctx context.Context, ID string) (bill models.Bill, err error)
	//UpdateBill(ctx context.Context,
	//	billStatus bool,
	//	billName string,
	//	currentSum float32,
	//	date string,
	//) (bill models.Bill, err error)
	//DeleteBill(ctx context.Context, ID string) (success bool, err error)
}

func Register(gRPCServer *grpc.Server, accounting Accounting) {
	billsv1.RegisterPigletBillsServer(gRPCServer, &serverAPI{accounting: accounting})
}

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
		if errors.Is(err, accounting.ErrUserExists) {
			return nil, status.Error(codes.InvalidArgument, "invalid email or password")
		}

		return nil, status.Errorf(codes.Internal, "internal error")
	}

	// HACK: обработка ошибок
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

// TODO: возврат нескольких счетов, getSomeBills получает массив с id (?)
func (s *serverAPI) GetSomeBills(
	ctx context.Context,
	req *billsv1.GetSomeBillsRequest,
) (*billsv1.GetSomeBillsResponse, error) {
	var bills []*billsv1.Bill

	return &billsv1.GetSomeBillsResponse{
		Bills: bills,
	}, nil
}

func (s *serverAPI) GetBill(
	ctx context.Context,
	req *billsv1.GetBillRequest,
) (*billsv1.BillResponse, error) {
	// TODO: setup validation

	// TODO: setup logic

	return &billsv1.BillResponse{
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
) (*billsv1.BillResponse, error) {
	// TODO: setup validation

	// TODO: setup logic

	return &billsv1.BillResponse{
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

type validateData struct {
	billType bool   `validate:"boolean"`
	billName string `validate:"required"`
}
