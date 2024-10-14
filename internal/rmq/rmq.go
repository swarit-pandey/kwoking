package rmq

import (
	"fmt"
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Opts is rabbitmq options that are required to setup and flush data
type Opts struct {
	ExchangeName string
	RoutingKey   string
	Username     string
	Password     string
	Connection   string
	ContentType  string
	ExchangeType string
	Durable      bool
	AutoDeleted  bool
	Internal     bool
	NoWait       bool
}

// RMQPublisher returns publisher
type Publisher struct {
	connection  *amqp.Connection
	channel     *amqp.Channel
	exchange    string
	routingKey  string
	contentType string
}

func NewRMQPublisher(opts *Opts) *Publisher {
	conn, err := amqp.Dial(opts.Connection)
	if err != nil {
		slog.Error("failed to dial", err)
		return nil
	}

	ch, err := conn.Channel()
	if err != nil {
		slog.Error("failed to get the channel", err)
		conn.Close()
		return nil
	}

	err = ch.ExchangeDeclare(
		opts.ExchangeName,
		opts.ExchangeType,
		opts.Durable,
		opts.AutoDeleted,
		opts.Internal,
		opts.NoWait,
		nil,
	)
	if err != nil {
		slog.Error("failed to declare the exchange", err)
		conn.Close()
		ch.Close()
		return nil
	}

	return &Publisher{
		connection:  conn,
		channel:     ch,
		exchange:    opts.ExchangeName,
		routingKey:  opts.RoutingKey,
		contentType: opts.ContentType,
	}
}

func (rp *Publisher) Publish(data []byte) error {
	if rp.channel == nil || rp.connection == nil {
		return fmt.Errorf("Channel or connection is not initialized yet")
	}

	return rp.channel.Publish(
		rp.exchange,
		rp.routingKey,
		false,
		false,
		amqp.Publishing{
			Body: data,
		},
	)
}

func (rp *Publisher) Close() {
	if rp.channel != nil {
		rp.channel.Close()
	}
	if rp.connection != nil {
		rp.connection.Close()
	}
}
