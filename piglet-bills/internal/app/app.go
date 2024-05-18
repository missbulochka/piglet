package app

import (
	"fmt"
	"log/slog"

	grpcapp "piglet-bills-service/internal/app/grpc"
	"piglet-bills-service/internal/config"
	"piglet-bills-service/internal/services/accounting"
	pgmigration "piglet-bills-service/internal/storage/pg-migration"
	"piglet-bills-service/internal/storage/psql"
)

type App struct {
	GRPCSrv  *grpcapp.App
	TransCli *grpcapp.Client
}

func New(
	log *slog.Logger,
	cfg *config.Config,
) *App {
	if err := pgmigration.RunMigration(
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
		cfg.DB.DBHost,
		cfg.DB.DBPort,
		cfg.DB.UserName,
		cfg.DB.Password,
		cfg.DB.DBName,
	)
	if err != nil {
		panic(err)
	}

	accountingService := accounting.New(
		log,
		storage,
		storage,
	)

	grpcTransCli, err := grpcapp.NewClientConnect(log, cfg.GRPC.GRPCTransCliServer, cfg.GRPC.GRPCTransCliPort)
	if err != nil {
		panic(err)
	}

	grpcApp := grpcapp.New(
		log,
		accountingService,
		grpcTransCli.Conn,
		cfg.GRPC.GRPCServer,
		cfg.GRPC.GRPCPort,
	)
	return &App{
		GRPCSrv:  grpcApp,
		TransCli: grpcTransCli,
	}
}
