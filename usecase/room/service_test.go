package room_test

import (
	"testing"

	"github.com/vsantosalmeida/browser-chat/entity"
	"github.com/vsantosalmeida/browser-chat/usecase/room"
	"github.com/vsantosalmeida/browser-chat/usecase/room/mocks"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

var errDB = errors.New("db error")

func TestService_CreateMessage(t *testing.T) {
	var (
		message = &entity.Message{
			UserID:  1,
			RoomID:  3,
			Content: "hello world!",
		}
	)

	repository := mocks.NewRepository(t)
	svc := room.NewService(repository)

	repository.
		On("CreateMessage", message).
		Return(nil).
		Once()

	err := svc.CreateMessage(1, 3, "hello world!")
	assert.NoError(t, err)
}

func TestService_CreateMessageError(t *testing.T) {
	var (
		message = &entity.Message{
			UserID:  1,
			RoomID:  3,
			Content: "hello world!",
		}

		expected = "could not create message on DB: db error"
	)

	repository := mocks.NewRepository(t)
	svc := room.NewService(repository)

	repository.
		On("CreateMessage", message).
		Return(errDB).
		Once()

	err := svc.CreateMessage(1, 3, "hello world!")
	assert.EqualError(t, err, expected)
}

func TestService_CreateRoom(t *testing.T) {
	var expected = 1

	repository := mocks.NewRepository(t)
	svc := room.NewService(repository)

	repository.
		On("CreateRoom", &entity.Room{}).
		Return(1, nil).
		Once()

	id, err := svc.CreateRoom()
	assert.NoError(t, err)
	assert.Equal(t, expected, id)
}

func TestService_CreateRoomError(t *testing.T) {
	var expected = "could not create room on DB: db error"

	repository := mocks.NewRepository(t)
	svc := room.NewService(repository)

	repository.
		On("CreateRoom", &entity.Room{}).
		Return(0, errDB).
		Once()

	id, err := svc.CreateRoom()
	assert.EqualError(t, err, expected)
	assert.Empty(t, id)
}

func TestService_ListMessages(t *testing.T) {
	var (
		messagesList = []*entity.Message{
			{
				ID:      1,
				UserID:  2,
				RoomID:  1,
				Content: "hello!",
			},
			{
				ID:      2,
				UserID:  5,
				RoomID:  1,
				Content: "hello again",
			},
		}

		expected = []*entity.Message{
			{
				ID:      1,
				UserID:  2,
				RoomID:  1,
				Content: "hello!",
			},
			{
				ID:      2,
				UserID:  5,
				RoomID:  1,
				Content: "hello again",
			},
		}
	)

	repository := mocks.NewRepository(t)
	svc := room.NewService(repository)

	repository.
		On("ListMessages", 1).
		Return(messagesList, nil).
		Once()

	messages, err := svc.ListMessages(1)
	assert.NoError(t, err)
	assert.Equal(t, expected, messages)
}

func TestService_ListMessagesError(t *testing.T) {
	var expected = "could not retrieve messages list: db error"

	repository := mocks.NewRepository(t)
	svc := room.NewService(repository)

	repository.
		On("ListMessages", 1).
		Return(nil, errDB).
		Once()

	messages, err := svc.ListMessages(1)
	assert.EqualError(t, err, expected)
	assert.Empty(t, messages)
}

func TestService_ListRooms(t *testing.T) {
	var (
		roomsList = []*entity.Room{
			{
				ID: 1,
			},
			{
				ID: 2,
			},
		}

		expected = []*entity.Room{
			{
				ID: 1,
			},
			{
				ID: 2,
			},
		}
	)

	repository := mocks.NewRepository(t)
	svc := room.NewService(repository)

	repository.
		On("ListRooms").
		Return(roomsList, nil).
		Once()

	rooms, err := svc.ListRooms()
	assert.NoError(t, err)
	assert.Equal(t, expected, rooms)
}

func TestService_ListRoomsError(t *testing.T) {
	var expected = "could not retrieve rooms list: db error"

	repository := mocks.NewRepository(t)
	svc := room.NewService(repository)

	repository.
		On("ListRooms").
		Return(nil, errDB).
		Once()

	rooms, err := svc.ListRooms()
	assert.EqualError(t, err, expected)
	assert.Empty(t, rooms)
}
