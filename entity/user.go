package entity

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User represents a User stored in the DB.
type User struct {
	ID        int    `gorm:"primaryKey"`
	Username  string `gorm:"index:idx_username,unique"`
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewUser User builder.
func NewUser(username, password string) (*User, error) {
	if err := validate(username, password); err != nil {
		return nil, ErrInvalidEntity
	}

	u := &User{
		Username: username,
	}
	pwd, err := generatePassword(password)
	if err != nil {
		return nil, err
	}
	u.Password = pwd
	return u, nil
}

// validate validate data.
func validate(username, password string) error {
	if username == "" || password == "" {
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
