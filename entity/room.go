package entity

import "time"

type Room struct {
	ID        int `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

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
