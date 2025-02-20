package rabbitmq

import (
	"fmt"
	"multitenant/pkg/config"

	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	config     *config.Config
	Connection *amqp.Connection
}

func NewRabbitMQ(cfg *config.Config) (*RabbitMQ, error) {
	uri := fmt.Sprintf("%s://%s:%s@%s:%s/",
		cfg.RabbitMQ.Protocol,
		cfg.RabbitMQ.Username,
		cfg.RabbitMQ.Password,
		cfg.RabbitMQ.Host,
		cfg.RabbitMQ.Port,
	)
	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, err
	}

	return &RabbitMQ{
		config:     cfg,
		Connection: conn,
	}, nil
}

func (r *RabbitMQ) Close() error {
	return r.Connection.Close()
}
