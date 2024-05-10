package rmq

import (
	"fmt"
	"log/slog"

	"github.com/streadway/amqp"

	"piglet-bills-service/internal/config"
)

type Connection struct {
	log        *slog.Logger
	connection *amqp.Connection
}

// InitConnection inits connections and in-out channels with RMQ with provided config
func InitConnection(log *slog.Logger, cfg *config.RabbitMQConfig) (*Connection, error) {
	const op = "piglet-bills | rmq.InitConnection"

	log = log.With(
		slog.String("op", op),
	)

	rmqConn, err := amqp.Dial(fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		cfg.RMQUser,
		cfg.RMQPassword,
		cfg.Server,
		cfg.Port,
	))
	if err != nil {
		return nil, fmt.Errorf("init connection error: %w", err)
	}

	log.Info("successfully init RabbitMQ")

	return &Connection{
		log:        log,
		connection: rmqConn,
	}, nil
}

// Close closes connection to rmq, to graceful shutdown
func (c *Connection) Close() error {
	if err := c.connection.Close(); err != nil {
		return err
	}
	return nil
}
