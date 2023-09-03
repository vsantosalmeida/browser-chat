package entity

import "time"

// Room represents a Room stored in the DB.
type Room struct {
	ID        int `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Message represents a Message stored in the DB.
type Message struct {
	ID        int `gorm:"primaryKey"`
	UserID    int
	RoomID    int
	User      User
	Room      Room
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
