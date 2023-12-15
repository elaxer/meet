package service

import (
	"errors"
	"meet/internal/app/model"
	"meet/internal/app/repository"
)

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

func (ms *MessageService) Text(usersDirection model.Direction, text string) (*model.Message, error) {
	isCouple, err := ms.assessmentService.IsCouple(usersDirection)
	if err != nil {
		return nil, err
	}
	if !isCouple {
		return nil, errors.New("невозможно отправить сообщение, так как пользователи не находятся в паре")
	}

	m := &model.Message{
		UsersDirection: usersDirection,
		Text:           text,
	}
	if err := ms.messageRepository.Add(m); err != nil {
		return nil, err
	}

	return m, nil
}

func (ms *MessageService) Read(userID, messageID int) (*model.Message, error) {
	m, err := ms.messageRepository.Get(messageID)
	if err != nil {
		return nil, err
	}

	if m.UsersDirection.ToID != userID {
		return nil, errors.New("невозможно прочитать сообщение, так как оно не было отправлено данному пользователю")
	}

	if m.IsReaded {
		return m, nil
	}

	m.Read()
	err = ms.messageRepository.Update(m)

	return m, err
}
