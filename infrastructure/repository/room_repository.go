package repository

import (
	"github.com/vsantosalmeida/browser-chat/entity"

	"gorm.io/gorm"
)

const maxMessages = 50

// RoomMySQL mysql repo
type RoomMySQL struct {
	db *gorm.DB
}

// NewRoomMySQL create new repository
func NewRoomMySQL(db *gorm.DB) *RoomMySQL {
	return &RoomMySQL{
		db: db,
	}
}

func (r *RoomMySQL) ListRooms() ([]*entity.Room, error) {
	var rooms []*entity.Room
	if result := r.db.Find(&rooms); result.Error != nil {
		return nil, result.Error
	}

	return rooms, nil
}

func (r *RoomMySQL) ListMessages(roomID int) ([]*entity.Message, error) {
	var mgs []*entity.Message
	if result := r.db.Where("room_id = ?", roomID).Limit(maxMessages).Order("created_at desc").Find(&mgs); result.Error != nil {
		return nil, result.Error
	}

	return mgs, nil
}

func (r *RoomMySQL) CreateRoom(e *entity.Room) (int, error) {
	if result := r.db.Create(e); result.Error != nil {
		return 0, result.Error
	}

	return e.ID, nil
}

func (r *RoomMySQL) CreateMessage(e *entity.Message) error {
	if result := r.db.Create(e); result.Error != nil {
		return result.Error
	}

	return nil
}
