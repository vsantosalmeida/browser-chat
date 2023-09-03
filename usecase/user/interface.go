package user

import "github.com/vsantosalmeida/browser-chat/entity"

// Reader handle the required methods to read users DB.
type Reader interface {
	FindByUsername(username string) (*entity.User, error)
	List() ([]*entity.User, error)
}

// Writer handle the required methods to write users DB.
type Writer interface {
	Create(e *entity.User) (int, error)
}

// Repository interface to bind Reader and Writer methods.
type Repository interface {
	Reader
	Writer
}

// UseCase service to handle the business rules for user context.
type UseCase interface {
	Authenticate(username, password string) (string, error)
	ListUsers() ([]*entity.User, error)
	CreateUser(username, password string) (int, error)
}
