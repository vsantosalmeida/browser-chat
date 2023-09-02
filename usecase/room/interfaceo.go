package room

import "github.com/vsantosalmeida/browser-chat/entity"

type Reader interface {
	ListRooms() ([]*entity.Room, error)
	ListMessages(roomID int) ([]*entity.Message, error)
}

type Writer interface {
	CreateRoom(e *entity.Room) (int, error)
	CreateMessage(e *entity.Message) error
}

type Repository interface {
	Reader
	Writer
}

type UseCase interface {
	ListRooms() ([]*entity.Room, error)
	ListMessages(roomID int) ([]*entity.Message, error)
	CreateRoom() (int, error)
	CreateMessage(userID, roomID int, content string) error
}
