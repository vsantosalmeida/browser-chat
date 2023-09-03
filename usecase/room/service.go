package room

import (
	"github.com/vsantosalmeida/browser-chat/entity"

	"github.com/apex/log"
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

func (s *Service) ListRooms() ([]*entity.Room, error) {
	rooms, err := s.repo.ListRooms()
	if err != nil {
		log.WithError(err).Error("could not retrieve rooms list")
		return nil, errors.Wrap(err, "could not retrieve rooms list")
	}

	return rooms, nil
}

func (s *Service) ListMessages(roomID int) ([]*entity.Message, error) {
	mgs, err := s.repo.ListMessages(roomID)
	if err != nil {
		log.WithError(err).Error("could not retrieve messages list")
		return nil, errors.Wrap(err, "could not retrieve messages list")
	}

	return mgs, nil
}

func (s *Service) CreateRoom() (int, error) {
	id, err := s.repo.CreateRoom(&entity.Room{})
	if err != nil {
		log.WithError(err).Error("could not create room on DB")
		return 0, errors.Wrap(err, "could not create room on DB")
	}

	log.WithField("id", id).Info("room created")

	return id, nil
}

func (s *Service) CreateMessage(userID, roomID int, content string) error {
	logger := log.WithFields(log.Fields{
		"RoomID": roomID,
		"UserID": userID,
	})

	msg := entity.Message{
		UserID:  userID,
		RoomID:  roomID,
		Content: content,
	}

	if err := s.repo.CreateMessage(&msg); err != nil {
		logger.WithError(err).Error("could not create message on DB")
		return errors.Wrap(err, "could not create message on DB")
	}

	logger.Info("message created")

	return nil
}
