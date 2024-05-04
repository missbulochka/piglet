package app

import (
	"log/slog"
	grpcapp "piglet-bills-service/internal/app/grpc"
	psqlapp "piglet-bills-service/internal/app/postgres"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcServer string,
	grpcPort int,
	storagePath string,
	migrationPath string,
	dbPort string,
	dbName string,
) *App {
	// TODO: инициализировать хранилище

	migrationApp := psqlapp.New(
		log,
		migrationPath,
		dbPort,
		dbName,
	)
	migrationApp.MustRunMigrations()

	// TODO: инициализировать сервис

	grpcApp := grpcapp.New(log, grpcServer, grpcPort)
	return &App{
		GRPCSrv: grpcApp,
	}
}
