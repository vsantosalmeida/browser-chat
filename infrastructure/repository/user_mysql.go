package repository

import (
	"github.com/vsantosalmeida/browser-chat/entity"

	"gorm.io/gorm"
)

// UserMySQL mysql repo
type UserMySQL struct {
	db *gorm.DB
}

// NewUserMySQL create new repository
func NewUserMySQL(db *gorm.DB) *UserMySQL {
	return &UserMySQL{
		db: db,
	}
}

func (u *UserMySQL) FindByUsername(username string) (*entity.User, error) {
	var user *entity.User
	if result := u.db.Find(&user, "username = ?", username); result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (u *UserMySQL) List() ([]*entity.User, error) {
	var users []*entity.User
	if result := u.db.Find(&users); result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}

func (u *UserMySQL) Create(e *entity.User) (int, error) {
	if result := u.db.Create(e); result.Error != nil {
		return 0, result.Error
	}

	return e.ID, nil
}
