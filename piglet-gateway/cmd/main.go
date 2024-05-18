package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	billsv1 "github.com/missbulochka/protos/gen/piglet-bills"
	transv1 "github.com/missbulochka/protos/gen/piglet-transactions"
	"piglet-manager-service/internal/config"
)

func main() {
	cfg := config.MustLoadConfig()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Listen for OS signals to gracefully shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan
		cancel()
	}()

	mux := runtime.NewServeMux()

	// Listen accounting
	accOpts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := billsv1.RegisterPigletBillsHandlerFromEndpoint(
		ctx,
		mux,
		fmt.Sprintf("%s:%s", cfg.GRPC.BillsServer, cfg.GRPC.BillsPort),
		accOpts)
	if err != nil {
		log.Fatal(err)
	}

	// Listen transactions
	transOpts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err = transv1.RegisterPigletTransactionsHandlerFromEndpoint(
		ctx,
		mux,
		fmt.Sprintf("%s:%s", cfg.GRPC.TransServer, cfg.GRPC.TransPort),
		transOpts)
	if err != nil {
		log.Fatal(err)
	}

	// Start HTTP server
	log.Println("starting HTTP server...")
	if err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.GRPC.GatewayPort), mux); err != nil {
		log.Fatal(err)
	}
	log.Println("server started")
}
