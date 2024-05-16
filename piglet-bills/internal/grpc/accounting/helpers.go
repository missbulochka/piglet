package accountingrpc

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	billsv1 "github.com/missbulochka/protos/gen/piglet-bills"
	transv1 "github.com/missbulochka/protos/gen/piglet-transactions"
	models "piglet-bills-service/internal/domain/model"
)

func syncBills(id string, cli transv1.PigletTransactionsClient, billStatus bool, del bool) {
	go func(id string, cli transv1.PigletTransactionsClient, billStatus bool) {
		ctx := context.Background()
		_, err := cli.AddBill(ctx, &transv1.Bill{Id: id, BillStatus: billStatus, Deletion: del})
		if err != nil {
			fmt.Println("service synchronization error: %w", err)
		}
	}(id, cli, billStatus)
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
