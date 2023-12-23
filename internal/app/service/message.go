package service

import (
	"errors"
	"meet/internal/app/model"
	"meet/internal/app/repository"
)

var ErrUsersNotCoupled = errors.New("пользователи не находятся в паре")

type MessageService struct {
	assessmentRepository repository.AssessmentRepository
	messageRepository    repository.MessageRepository
	assessmentService    *AssessmentService
}

func newMessageService(
	assessmentRepository repository.AssessmentRepository,
	messageRepository repository.MessageRepository,
	assessmentService *AssessmentService,
) *MessageService {
	return &MessageService{assessmentRepository, messageRepository, assessmentService}
}

func (ms *MessageService) Text(message *model.Message) error {
	if err := message.Validate(); err != nil {
		return err
	}

	isCouple, err := ms.assessmentService.IsCouple(message.UsersDirection)
	if err != nil {
		return err
	}
	if !isCouple {
		return ErrUsersNotCoupled
	}

	err = ms.messageRepository.Add(message)

	return err
}

func (ms *MessageService) Read(userID, messageID int) (*model.Message, error) {
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
