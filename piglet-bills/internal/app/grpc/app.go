package grpcapp

import (
	"fmt"
	"log/slog"
	"net"

	"google.golang.org/grpc"

	accountingrpc "piglet-bills-service/internal/grpc/accounting"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	server     string
	port       string
}

func New(
	log *slog.Logger,
	accountingService accountingrpc.Accounting,
	server string,
	port string,
) *App {
	gRPCServer := grpc.NewServer()

	accountingrpc.Register(gRPCServer, accountingService)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		server:     server,
		port:       port,
	}
}

func (a *App) MustStart() {
	if err := a.start(); err != nil {
		panic(err)
	}
}

func (a *App) start() error {
	const op = "piglet-bills | grpcapp.Start"

	log := a.log.With(
		slog.String("op", op),
		slog.String("port", a.port),
	)

	l, err := net.Listen("tcp", fmt.Sprintf("%s:%s", a.server, a.port))
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
		Info("stopping grpc server", slog.String("port", a.port))

	a.gRPCServer.GracefulStop()
}
