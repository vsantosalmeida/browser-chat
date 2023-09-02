package entity

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int    `gorm:"primaryKey"`
	Username  string `gorm:"index:idx_username,unique"`
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewUser create a new user
func NewUser(userName, password string) (*User, error) {
	u := &User{
		Username: userName,
	}
	pwd, err := generatePassword(password)
	if err != nil {
		return nil, err
	}
	u.Password = pwd

	if err = u.Validate(); err != nil {
		return nil, ErrInvalidEntity
	}
	return u, nil
}

// Validate validate data
func (u *User) Validate() error {
	if u.Username == "" || u.Password == "" {
		return ErrInvalidEntity
	}

	return nil
}

// ValidatePassword validate user password
func (u *User) ValidatePassword(p string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(p))
	if err != nil {
		return ErrInvalidPassword
	}
	return nil
}

func generatePassword(raw string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(raw), 10)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
