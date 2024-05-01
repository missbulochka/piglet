package main

import (
	"log/slog"
	"os"
	"piglet-bills/internal/config"
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

	// TODO: инициализировать приложение

	// TODO: запустить gRPC-сервер приложения

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
