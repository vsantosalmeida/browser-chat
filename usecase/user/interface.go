package user

import "github.com/vsantosalmeida/browser-chat/entity"

// Reader interface
type Reader interface {
	FindByUsername(username string) (*entity.User, error)
	List() ([]*entity.User, error)
}

// Writer user writer
type Writer interface {
	Create(e *entity.User) (int, error)
}

// Repository interface
type Repository interface {
	Reader
	Writer
}

// UseCase interface
type UseCase interface {
	Authenticate(username, password string) (string, error)
	ListUsers() ([]*entity.User, error)
	CreateUser(username, password string) (int, error)
}
