package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"hello_world/database"
	"hello_world/model"
	desc "hello_world/pkg/user_v1"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const grpcPort = 8080

type server struct {
	desc.UnsafeBillServiceServer
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

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	database.Init()

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterBillServiceServer(s, &server{})

	log.Printf("server listening at %v", lis.Addr())

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	log.Println("Shutting down server...")
	s.GracefulStop()
	log.Println("Server stopped")
}
