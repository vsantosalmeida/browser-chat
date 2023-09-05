package websocket

import (
	"context"
	"testing"
	"time"

	brokerMock "github.com/vsantosalmeida/browser-chat/api/websocket/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/undefinedlabs/go-mpatch"
)

const mockAnythingOfTypeChanByte = "chan<- []uint8"

func TestServerListenChatbotMessages(t *testing.T) {
	var (
		msgRaw         = `{"roomID":1,"from":"chat-bot","message":"command executed"}`
		eventOutputRaw = `{"message":"command executed","from":"chat-bot","sent":"2020-01-01T00:00:00Z"}`

		expected = Event{
			Action:  MessageReceivedAction,
			Payload: []byte(eventOutputRaw),
		}
	)

	eventCH := make(chan Event, 1)

	broker := brokerMock.NewBroker(t)

	s := &Server{
		handlers: initEventHandlers(),
		rooms:    rooms,
		clients:  make(map[*Client]bool),
		broker:   broker,
	}

	c := &Client{
		server: s,
		event:  eventCH,
		ID:     10,
		RoomID: 1,
	}

	s.joinClient(c)

	ctx := context.Background()

	// bypass time.Now function to set a static date for sent time
	timePatch, err := mpatch.PatchMethod(time.Now, func() time.Time {
		return time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	})
	assert.NoError(t, err)
	defer timePatch.Unpatch()

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

	go s.listenChatbotMessages(ctx)
	got := <-eventCH
	assert.Equal(t, expected, got)
}
