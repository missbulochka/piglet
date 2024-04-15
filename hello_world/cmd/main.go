package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"hello_world/internal/server"
	"hello_world/service/database"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const grpcPort = 8080

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	database.Init()

	s := grpc.NewServer()
	reflection.Register(s)
	server.NewServer(s)

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
