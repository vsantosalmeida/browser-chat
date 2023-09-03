package broker

import (
	"context"

	"github.com/apex/log"
	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQ message broker.
type RabbitMQ struct {
	ch         *amqp.Channel
	readQueue  string
	writeQueue string
}

// NewRabbitMQ RabbitMQ builder.
func NewRabbitMQ(readQueue, writeQueue string, ch *amqp.Channel) *RabbitMQ {
	return &RabbitMQ{
		ch:         ch,
		readQueue:  readQueue,
		writeQueue: writeQueue,
	}
}

// ReadMessage loop through rabbitMQ delivery channel and send the message body to the msgCH.
// a context.Canceled error will stop the consumer.
func (r *RabbitMQ) ReadMessage(ctx context.Context, msgCH chan<- []byte) {
	msgs, err := r.ch.Consume(
		r.readQueue,
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		log.WithError(err).Error("failed to consume messages")
		return
	}

	for msg := range msgs {
		if ctx.Err() == context.Canceled {
			log.Warn("context canceled")
			break
		}

		log.WithField("Message", msg).Debug("message received")
		msgCH <- msg.Body
	}

	log.Info("consumer stopped")
}

// WriteMessage send message to the configured queue.
func (r *RabbitMQ) WriteMessage(ctx context.Context, payload []byte) error {
	if err := r.ch.PublishWithContext(ctx,
		"",           // exchange
		r.writeQueue, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        payload,
		}); err != nil {
		log.WithError(err).Error("failed to send message")
		return err
	}

	log.WithField("Message", string(payload)).Debug("message sent")
	log.Info("message sent")
	return nil
}
