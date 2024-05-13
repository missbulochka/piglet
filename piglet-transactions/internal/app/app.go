package app

import (
	"log/slog"

	grpcapp "piglet-transactions-service/internal/app/grpc"
	"piglet-transactions-service/internal/config"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(log *slog.Logger, cfg *config.Config) *App {
	// TODO: run migration

	// TODO: connect data base

	// TODO: add service

	// TODO: setup grpc-server

	grpcApp := grpcapp.New(log, cfg.GRPC.GRPCServer, cfg.GRPC.GRPCPort)

	return &App{
		GRPCSrv: grpcApp,
	}
}
