package app

import (
	"fmt"
	"log/slog"

	grpcapp "piglet-transactions-service/internal/app/grpc"
	"piglet-transactions-service/internal/config"
	"piglet-transactions-service/internal/services/transactions"
	migrator "piglet-transactions-service/internal/storage/pg-migration"
	psql "piglet-transactions-service/internal/storage/postgres"
)

type App struct {
	GRPCSrv  *grpcapp.App
	BillsCli *grpcapp.Client
}

func New(log *slog.Logger, cfg *config.Config) *App {
	if err := migrator.RunMigration(
		log,
		"file://"+cfg.DB.MigrationPath,
		fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
			cfg.DB.UserName,
			cfg.DB.Password,
			cfg.DB.DBHost,
			cfg.DB.DBPort,
			cfg.DB.DBName,
		),
	); err != nil {
		panic(err)
	}

	storage, err := psql.New(
		log,
		fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			cfg.DB.DBHost,
			cfg.DB.DBPort,
			cfg.DB.UserName,
			cfg.DB.Password,
			cfg.DB.DBName,
		),
	)
	if err != nil {
		panic(err)
	}

	transService := transactions.New(log, storage, storage, storage, storage, storage)

	grpcBillsCli, err := grpcapp.NewClientConnect(log, cfg.GRPC.GRPCBillsCliServer, cfg.GRPC.GRPCBillsCliPort)
	if err != nil {
		panic(err)
	}

	grpcApp := grpcapp.New(
		log,
		transService,
		transService,
		transService,
		grpcBillsCli.Conn,
		cfg.GRPC.GRPCServer,
		cfg.GRPC.GRPCPort,
	)

	return &App{
		GRPCSrv:  grpcApp,
		BillsCli: grpcBillsCli,
	}
}
