package grpcapp

import (
	"fmt"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	log    *slog.Logger
	Conn   *grpc.ClientConn
	Server string
	Port   string
}

func NewClientConnect(
	log *slog.Logger,
	server string,
	port string,
) (*Client, error) {
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

	return &Client{
		log:    log,
		Conn:   conn,
		Server: server,
		Port:   port,
	}, nil
}

func (cli *Client) ConnClose() {
	err := cli.Conn.Close()
	if err != nil {
		cli.log.Info("wrong connection closing")
	}
}
