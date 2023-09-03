package config

import (
	"github.com/apex/log"
	amqp "github.com/rabbitmq/amqp091-go"
)

func InitRabbitMQ() *amqp.Channel {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.WithError(err).Fatal("could not dial to rabbitmq")
	}

	ch, err := conn.Channel()
	if err != nil {
		log.WithError(err).Fatal("could not connect to rabbitmq channel")
	}

	return ch
}
