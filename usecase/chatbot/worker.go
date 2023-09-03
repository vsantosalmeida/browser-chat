package chatbot

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/apex/log"
)

// Worker handle the commands sent by a user.
type Worker struct {
	id  int
	svc *Service
}

// StartAndConsume loop through a message channel to receive user commands and process it.
// sends an output to another queue with the command result.
func (w *Worker) StartAndConsume(ctx context.Context) {
	logger := log.WithField("WorkerID", w.id)

	msgCH := make(chan []byte)
	go w.svc.broker.ReadMessage(ctx, msgCH)

	for msg := range msgCH {
		if ctx.Err() == context.Canceled {
			logger.Warn("context canceled")
			break
		}

		var input CommandInput

		if err := json.Unmarshal(msg, &input); err != nil {
			logger.WithError(err).Error("could not decode message body")
			continue
		}

		logger.WithField("CommandInput", input).Info("command received")

		var message string
		if result, err := w.svc.ExecuteCommand(ctx, input.CommandName, input.Command); err == nil {
			message = result
		} else {
			message = fmt.Sprintf("could not execute command: %s error: %q", input.CommandName, err)
		}

		output := CommandOutput{
			RoomID:  input.RoomID,
			From:    "chat-bot",
			Message: message,
		}

		b, err := json.Marshal(output)
		if err != nil {
			logger.WithError(err).Error("could not encode message body")
			continue
		}

		if err = w.svc.broker.WriteMessage(ctx, b); err != nil {
			continue
		}
	}

	logger.Info("worker stopped")
}
