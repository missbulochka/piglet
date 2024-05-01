package app

import (
	"log/slog"
	grpcapp "piglet-bills-service/internal/app/grpc"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcServer string,
	grpcPort int,
) *App {
	// TODO: инициализировать хранилище

	// TODO: инициализировать сервис

	grpcApp := grpcapp.New(log, grpcServer, grpcPort)
	return &App{
		GRPCSrv: grpcApp,
	}
}
