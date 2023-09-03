package config

import (
	"fmt"

	"github.com/apex/log"
	amqp "github.com/rabbitmq/amqp091-go"
)

const amqpURL = "amqp://%s:%s@%s:5672/"

func InitRabbitMQ() *amqp.Channel {
	url := fmt.Sprintf(
		amqpURL,
		GetStingEnvVarOrPanic(RabbitMQUser),
		GetStingEnvVarOrPanic(RabbitMQPass),
		GetStingEnvVarOrPanic(RabbitMQHost),
	)

	conn, err := amqp.Dial(url)
	if err != nil {
		log.WithError(err).Fatal("could not dial to rabbitmq")
	}

	ch, err := conn.Channel()
	if err != nil {
		log.WithError(err).Fatal("could not connect to rabbitmq channel")
	}

	return ch
}
