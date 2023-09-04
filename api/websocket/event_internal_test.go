package websocket

import (
	"context"
	"testing"
	"time"

	brokerMock "github.com/vsantosalmeida/browser-chat/api/websocket/mocks"
	"github.com/vsantosalmeida/browser-chat/entity"
	"github.com/vsantosalmeida/browser-chat/usecase/room"
	roomMock "github.com/vsantosalmeida/browser-chat/usecase/room/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/undefinedlabs/go-mpatch"
)

var (
	rooms = []*entity.Room{
		{
			ID: 1,
		},
	}
)

func TestSendMessageHandler(t *testing.T) {
	var (
		eventInputRaw  = `{"message":"hello world!","from":"user"}`
		eventOutputRaw = `{"message":"hello world!","from":"user","sent":"2020-01-01T00:00:00Z"}`
		event          = Event{
			Action:  SendMessageAction,
			Payload: []byte(eventInputRaw),
		}

		msg = &entity.Message{
			UserID:  10,
			RoomID:  1,
			Content: "hello world!",
		}

		expected = Event{
			Action:  SendMessageAction,
			Payload: []byte(eventOutputRaw),
		}
	)

	eventCH := make(chan Event, 1)
	roomRepo := roomMock.NewRepository(t)
	roomUseCase := room.NewService(roomRepo)

	s := &Server{
		handlers:    initEventHandlers(),
		rooms:       rooms,
		roomUseCase: roomUseCase,
		clients:     make(map[*Client]bool),
	}

	c := &Client{
		server: s,
		event:  eventCH,
		ID:     10,
		RoomID: 1,
	}

	s.joinClient(c)

	// bypass time.Now function to set a static date for sent time
	timePatch, err := mpatch.PatchMethod(time.Now, func() time.Time {
		return time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	})
	assert.NoError(t, err)
	defer timePatch.Unpatch()

	roomRepo.
		On("CreateMessage", msg).
		Return(nil).
		Maybe()

	err = SendMessageHandler(event, c)
	assert.NoError(t, err)

	got := <-eventCH
	assert.Equal(t, expected, got)
}

func TestChatRoomHandler(t *testing.T) {
	var (
		eventInputRaw = `{"roomID":1}`
		event         = Event{
			Action:  JoinRoomAction,
			Payload: []byte(eventInputRaw),
		}
	)

	s := &Server{
		handlers: initEventHandlers(),
		rooms:    rooms,
		clients:  make(map[*Client]bool),
	}

	c := &Client{
		server: s,
		ID:     10,
		RoomID: 1,
	}

	s.joinClient(c)

	err := ChatRoomHandler(event, c)
	assert.NoError(t, err)
}

func TestChatbotCommandHandler(t *testing.T) {
	var (
		eventInputRaw = `{"roomID":1,"from":"user","commandName":"stock","command":"amzn.us"}`
		event         = Event{
			Action:  SendChatbotCommandAction,
			Payload: []byte(eventInputRaw),
		}
	)

	broker := brokerMock.NewBroker(t)

	s := &Server{
		handlers: initEventHandlers(),
		rooms:    rooms,
		clients:  make(map[*Client]bool),
		broker:   broker,
	}

	c := &Client{
		server: s,
		ID:     10,
		RoomID: 1,
	}

	s.joinClient(c)

	broker.
		On("WriteMessage", context.Background(), []byte(eventInputRaw)).
		Return(nil).
		Maybe()

	err := ChatbotCommandHandler(event, c)
	assert.NoError(t, err)
}
