package presenter

import (
	"time"

	"github.com/vsantosalmeida/browser-chat/entity"
)

type CreateRoomOutput struct {
	ID int `json:"id"`
}

type Room struct {
	ID int `json:"id"`
}

type Message struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

func MapEntityToExternalRooms(rooms []*entity.Room) []*Room {
	result := make([]*Room, 0)

	for _, room := range rooms {
		result = append(
			result,
			&Room{
				ID: room.ID,
			},
		)
	}

	return result
}

func MapEntityToExternalMessages(mgs []*entity.Message) []*Message {
	result := make([]*Message, 0)

	for _, m := range mgs {
		result = append(
			result,
			&Message{
				ID:        m.ID,
				Content:   m.Content,
				CreatedAt: m.CreatedAt,
			},
		)
	}

	return result
}
