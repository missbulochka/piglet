package grpcapp

import (
	"fmt"
	"google.golang.org/grpc"
	"log/slog"
	"net"
	accountingrpc "piglet-bills-service/internal/grpc/accounting"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	server     string
	port       int
}

func New(log *slog.Logger, server string, port int) *App {
	gRPCServer := grpc.NewServer()

	accountingrpc.Register(gRPCServer)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		server:     server,
		port:       port,
	}
}

func (a *App) MustStart() {
	if err := a.Start(); err != nil {
		panic(err)
	}
}

func (a *App) Start() error {
	const op = "piglet-bills | grpcapp.Start"

	log := a.log.With(
		slog.String("op", op),
		slog.Int("port", a.port),
	)

	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", a.server, a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("grpc server is running", slog.String("addr", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (a *App) Stop() {
	const op = "piglet-bills | grpcapp.Stop"

	a.log.With(slog.String("op", op)).
		Info("stopping grpc server", slog.Int("port", a.port))

	a.gRPCServer.GracefulStop()
}
