package room

import "github.com/vsantosalmeida/browser-chat/entity"

// Reader handle the required methods to read rooms DB.
type Reader interface {
	ListRooms() ([]*entity.Room, error)
	ListMessages(roomID int) ([]*entity.Message, error)
}

// Writer handle the required methods to write rooms DB.
type Writer interface {
	CreateRoom(e *entity.Room) (int, error)
	CreateMessage(e *entity.Message) error
}

// Repository interface to bind Reader and Writer methods.
type Repository interface {
	Reader
	Writer
}

// UseCase service to handle the business rules for room context.
type UseCase interface {
	ListRooms() ([]*entity.Room, error)
	ListMessages(roomID int) ([]*entity.Message, error)
	CreateRoom() (int, error)
	CreateMessage(userID, roomID int, content string) error
}
