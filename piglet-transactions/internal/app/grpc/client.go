package grpcapp

import (
	"fmt"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func ClientConnect(
	log *slog.Logger,
	server string,
	port string,
) (*grpc.ClientConn, error) {
	const op = "piglet-transactions | grpcapp.ClientConnect()"

	log = log.With(
		slog.String("op", op),
		slog.String("server", server),
		slog.String("port", port),
	)

	conn, err := grpc.Dial(
		fmt.Sprintf("%s:%s", server, port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("grpc connection for piglet-bills client is ready")

	return conn, nil
}
