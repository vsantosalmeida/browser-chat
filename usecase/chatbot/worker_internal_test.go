package chatbot

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/vsantosalmeida/browser-chat/usecase/chatbot/mocks"
)

const (
	mockAnythingOfTypeChanByte = "chan<- []uint8"
)

var ctx = context.Background()

func TestWorker_StartAndConsume(t *testing.T) {
	var (
		msgRaw           = `{"roomID":1,"from":"user","commandName":"mock","command":"any"}`
		commandOutputRaw = `{"roomID":1,"from":"chat-bot","message":"command executed"}`
	)

	broker := mocks.NewBroker(t)
	svc := &Service{
		broker:  broker,
		workers: 1,
		handlers: map[string]CommandHandler{
			"mock": mockCommandHandler,
		},
	}
	w := Worker{
		id:  1,
		svc: svc,
	}

	broker.
		On("ReadMessage", ctx, mock.AnythingOfType(mockAnythingOfTypeChanByte)).
		Return().
		Run(func(args mock.Arguments) {
			ch := args.Get(1).(chan<- []byte)
			ch <- []byte(msgRaw)
			// closes the channel to stop waiting for messages
			close(ch)
		}).
		Once()

	broker.
		On("WriteMessage", ctx, []byte(commandOutputRaw)).
		Return(nil).
		Once()

	w.StartAndConsume(ctx)
}

func TestWorker_StartAndConsumeCommandError(t *testing.T) {
	var (
		badCommandRaw       = `{"roomID":1,"from":"user","commandName":"mock","command":"any"}`
		badCommandOutputRaw = `{"roomID":1,"from":"chat-bot","message":"could not execute command: mock error: invalid command"}`
	)

	broker := mocks.NewBroker(t)
	svc := &Service{
		broker:  broker,
		workers: 1,
	}
	w := Worker{
		id:  1,
		svc: svc,
	}

	broker.
		On("ReadMessage", ctx, mock.AnythingOfType(mockAnythingOfTypeChanByte)).
		Return().
		Run(func(args mock.Arguments) {
			ch := args.Get(1).(chan<- []byte)
			ch <- []byte(badCommandRaw)
			// closes the channel to stop waiting for messages
			close(ch)
		}).
		Once()

	broker.
		On("WriteMessage", ctx, []byte(badCommandOutputRaw)).
		Return(nil).
		Once()

	w.StartAndConsume(ctx)
}

func mockCommandHandler(_ context.Context, _ string) (string, error) {
	return "command executed", nil
}
