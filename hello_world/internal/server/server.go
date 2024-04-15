package server

import (
	"context"
	"errors"
	"fmt"

	desc "hello_world/internal/proto"
	"hello_world/service/database"
	"hello_world/service/model"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type server struct {
	desc.UnsafeBillServiceServer
}

func NewServer(s *grpc.Server) {
	desc.RegisterBillServiceServer(s, &server{})
}

func (*server) CreateBill(ctx context.Context, req *desc.CreateBillRequest) (*desc.CreateBillResponse, error) {
	fmt.Println("Create bill")
	bill := req.GetBill()
	bill.Id = uuid.New().String()

	data := model.Bill{
		ID:       bill.GetId(),
		State:    bill.GetState(),
		Title:    bill.GetTitle(),
		Sum:      bill.GetSum(),
		Currency: bill.GetCurrency(),
	}

	res := database.DB.Create(&data)
	if res.RowsAffected == 0 {
		return nil, errors.New("bill creation unsuccessful")
	}
	return &desc.CreateBillResponse{
		Bill: &desc.Bill{
			Id:       bill.GetId(),
			State:    bill.GetState(),
			Title:    bill.GetTitle(),
			Sum:      bill.GetSum(),
			Currency: bill.GetCurrency(),
		},
	}, nil
}

func (*server) GetBill(ctx context.Context, req *desc.GetBillRequest) (*desc.GetBillResponse, error) {
	fmt.Println("Get bill", req.GetId())
	var bill model.Bill
	res := database.DB.Find(&bill, "id = ?", req.GetId())
	if res.RowsAffected == 0 {
		return nil, errors.New("bill not found")
	}
	return &desc.GetBillResponse{
		Bill: &desc.Bill{
			Id:       bill.ID,
			State:    bill.State,
			Title:    bill.Title,
			Sum:      bill.Sum,
			Currency: bill.Currency,
		},
	}, nil
}

func (*server) ReadBills(ctx context.Context, req *desc.ReadBillsRequest) (*desc.ReadBillsResponse, error) {
	fmt.Println("Read bills")
	var bills []*desc.Bill
	res := database.DB.Find(&bills)
	if res.RowsAffected == 0 {
		return nil, errors.New("bill not found")
	}
	return &desc.ReadBillsResponse{
		Bills: bills,
	}, nil
}

func (*server) UpdateBill(ctx context.Context, req *desc.UpdateBillRequest) (*desc.UpdateBillResponse, error) {
	fmt.Println("Update bill")
	var bill model.Bill
	reqBill := req.GetBill()

	res := database.DB.Model(&bill).Where("id=?", reqBill.Id).Updates(
		model.Bill{
			State:    reqBill.State,
			Title:    reqBill.Title,
			Sum:      reqBill.Sum,
			Currency: reqBill.Currency})

	if res.RowsAffected == 0 {
		return nil, errors.New("bill not found")
	}

	return &desc.UpdateBillResponse{
		Bill: &desc.Bill{
			Id:       bill.ID,
			State:    bill.State,
			Title:    bill.Title,
			Sum:      bill.Sum,
			Currency: bill.Currency,
		},
	}, nil
}

func (*server) DeleteBill(ctx context.Context, req *desc.DeleteBillRequest) (*desc.DeleteBillResponse, error) {
	fmt.Println("Delete bill")
	var bill model.Bill
	res := database.DB.Where("id = ?", req.GetId()).Delete(&bill)
	if res.RowsAffected == 0 {
		return nil, errors.New("bill not found")
	}

	return &desc.DeleteBillResponse{
		Success: true,
	}, nil
}
