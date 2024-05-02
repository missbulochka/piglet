package main

import (
	"log/slog"
	"os"
	"os/signal"
	"piglet-bills-service/internal/app"
	"piglet-bills-service/internal/config"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// TODO: иницилизировать объект конфига (+RabbitMQ)
	cfg := config.MustLoadConfig()

	log := setupLogger(cfg.Env)
	log.Info("starting piglet-bills service", slog.Any("config", cfg))

	application := app.New(log, cfg.GRPC.Server, cfg.GRPC.Port)

	go func() {
		application.GRPCSrv.MustStart()
	}()

	// Graceful shutdown

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop
	log.Info("shutting down piglet-bills service", slog.String("signal", sign.String()))
	application.GRPCSrv.Stop()
	log.Info("piglet-bills service stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}