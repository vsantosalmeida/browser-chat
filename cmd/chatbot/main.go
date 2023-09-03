package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vsantosalmeida/browser-chat/config"
	"github.com/vsantosalmeida/browser-chat/infrastructure/broker"
	"github.com/vsantosalmeida/browser-chat/pkg/stooq"
	"github.com/vsantosalmeida/browser-chat/usecase/chatbot"

	"github.com/apex/log"
)

func main() {
	config.InitLogging()
	ch := config.InitRabbitMQ()

	readQueue, err := ch.QueueDeclare(
		config.GetStingEnvVarOrPanic(config.ChatbotCommandInputQueue), // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.WithError(err).Fatal("could not create read queue")
	}

	writeQeue, err := ch.QueueDeclare(
		config.GetStingEnvVarOrPanic(config.ChatbotCommandOutputQueue), // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.WithError(err).Fatal("could not create write queue")
	}

	httpClient := &http.Client{
		Timeout: 20 * time.Second,
	}

	stooqAPI := stooq.NewClient(httpClient)
	rabbitMQ := broker.NewRabbitMQ(readQueue.Name, writeQeue.Name, ch)
	svc := chatbot.NewService(rabbitMQ, stooqAPI, config.GetIntEnvVarOrDefault(config.ChatbotMaxWorkers, 1))

	ctx, cancel := context.WithCancel(context.Background())

	go svc.Start(ctx)
	log.Info("chatbot started")

	/// gracefully shutdown the workers
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)

	<-sc
	cancel()
	ch.Close()

	log.Info("chatbot stopped")
}
