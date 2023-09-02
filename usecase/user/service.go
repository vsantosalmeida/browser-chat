package user

import (
	"github.com/vsantosalmeida/browser-chat/entity"
	"github.com/vsantosalmeida/browser-chat/pkg/auth"

	"github.com/pkg/errors"
)

// Service
type Service struct {
	repo Repository
}

// NewService
func NewService(r Repository) *Service {
	return &Service{
		repo: r,
	}
}

func (s *Service) Authenticate(username, password string) (string, error) {
	user, err := s.repo.FindByUsername(username)
	if err != nil {
		return "", errors.Wrap(err, "could not find user")
	}

	if err = user.ValidatePassword(password); err != nil {
		return "", errors.Wrap(err, "could not validate user")
	}

	token, err := auth.CreateJWTToken(user)
	if err != nil {
		return "", errors.Wrap(err, "could not generate user token")
	}

	return token, nil

}

func (s *Service) ListUsers() ([]*entity.User, error) {
	users, err := s.repo.List()
	if err != nil {
		return nil, errors.Wrap(err, "could not retrieve users list")
	}

	return users, nil
}

func (s *Service) CreateUser(username, password string) (int, error) {
	user, err := entity.NewUser(username, password)
	if err != nil {
		return 0, errors.Wrap(err, "could not create an user object")
	}

	id, err := s.repo.Create(user)
	if err != nil {
		return 0, errors.Wrap(err, "could not create user on DB")
	}

	return id, nil
}
