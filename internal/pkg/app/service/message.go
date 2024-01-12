package service

import (
	"errors"
	"meet/internal/pkg/app/model"
	"meet/internal/pkg/app/repository"
)

var (
	ErrUsersNotCoupled = errors.New("пользователи не находятся в паре")
)

type MessageService interface {
	Send(message *model.Message) error
	Read(userID, messageID int) (*model.Message, error)
}

type messageService struct {
	messageRepository repository.MessageRepository
	coupleRepository  repository.CoupleRepository
}

func NewMessageService(messageRepository repository.MessageRepository, coupleRepository repository.CoupleRepository) MessageService {
	return &messageService{messageRepository, coupleRepository}
}

func (ms *messageService) Send(message *model.Message) error {
	if err := message.Validate(); err != nil {
		return err
	}

	isCouple, err := ms.coupleRepository.Has(message.UsersDirection)
	if err != nil {
		return err
	}
	if !isCouple {
		return ErrUsersNotCoupled
	}

	err = ms.messageRepository.Add(message)

	return err
}

func (ms *messageService) Read(userID, messageID int) (*model.Message, error) {
	m, err := ms.messageRepository.Get(messageID)
	if err != nil {
		return nil, err
	}

	if m.UsersDirection.ToID != userID {
		return nil, repository.ErrNotFound
	}

	if m.IsReaded {
		return m, nil
	}

	m.Read()
	err = ms.messageRepository.Update(m)

	return m, err
}
