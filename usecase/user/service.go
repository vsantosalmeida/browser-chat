package user

import (
	"github.com/vsantosalmeida/browser-chat/entity"
	"github.com/vsantosalmeida/browser-chat/pkg/auth"

	"github.com/apex/log"
	"github.com/pkg/errors"
)

// Service implements UseCase interface.
type Service struct {
	repo Repository
}

// NewService Service builder.
func NewService(r Repository) *Service {
	return &Service{
		repo: r,
	}
}

// Authenticate retrieve user from DB and validate the password, if no error occurs a JWT token will be generated
// and returned.
func (s *Service) Authenticate(username, password string) (string, error) {
	user, err := s.repo.FindByUsername(username)
	if err != nil {
		log.WithError(err).Error("user not found")
		return "", errors.Wrap(err, "could not find user")
	}

	if err = user.ValidatePassword(password); err != nil {
		log.WithError(err).Error("could not validate user")
		return "", errors.Wrap(err, "could not validate user")
	}

	token, err := auth.CreateJWTToken(user)
	if err != nil {
		log.WithError(err).Error("could not generate user token")
		return "", errors.Wrap(err, "could not generate user token")
	}

	log.WithField("username", username).Info("user authenticated")

	return token, nil

}

// ListUsers retrieve all users from DB.
func (s *Service) ListUsers() ([]*entity.User, error) {
	users, err := s.repo.List()
	if err != nil {
		log.WithError(err).Error("could not retrieve users list")
		return nil, errors.Wrap(err, "could not retrieve users list")
	}

	return users, nil
}

// CreateUser validate user input and create it in the DB.
func (s *Service) CreateUser(username, password string) (int, error) {
	user, err := entity.NewUser(username, password)
	if err != nil {
		log.WithError(err).Error("could not create an user object")
		return 0, errors.Wrap(err, "could not create an user object")
	}

	id, err := s.repo.Create(user)
	if err != nil {
		log.WithError(err).Error("could not create user on DB")
		return 0, errors.Wrap(err, "could not create user on DB")
	}

	log.WithField("id", id).Info("user created")

	return id, nil
}
