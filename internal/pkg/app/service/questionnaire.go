package service

import (
	"context"
	"errors"
	"meet/internal/pkg/app/model"
	"meet/internal/pkg/app/repository"
	"time"
)

var (
	ErrQuestionnaireState = errors.New("невозможно выполнить действие с текущим состоянием анкеты")
)

type QuestionnaireService interface {
	PickUp(userID, limit, offset int) ([]*model.Questionnaire, error)
	Add(ctx context.Context, questionnaire *model.Questionnaire) error
	Update(questionnaire *model.Questionnaire) error
	UpdateName(questionnaire *model.Questionnaire, name string) error
}

type questionnaireService struct {
	questionnaireRepository repository.QuestionnaireRepository
}

func NewQuestionnaireService(questionnaireRepository repository.QuestionnaireRepository) QuestionnaireService {
	return &questionnaireService{questionnaireRepository}
}

func (qs *questionnaireService) PickUp(userID, limit, offset int) ([]*model.Questionnaire, error) {
	questionnaires := make([]*model.Questionnaire, 0, limit)

	questionnaire, err := qs.questionnaireRepository.GetByUserID(userID)
	if err != nil {
		return questionnaires, err
	}

	if !questionnaire.IsReady() {
		return questionnaires, nil
	}

	questionnaires, err = qs.questionnaireRepository.PickUp(userID, limit, offset)
	if err != nil {
		return questionnaires, err
	}

	compatibleQuestionnaires := make([]*model.Questionnaire, 0, len(questionnaires))
	for _, q := range questionnaires {
		if questionnaire.CheckCompatibilities(q, time.Now()) {
			compatibleQuestionnaires = append(compatibleQuestionnaires, q)
		}
	}

	return compatibleQuestionnaires, nil
}

func (qs *questionnaireService) Add(ctx context.Context, questionnaire *model.Questionnaire) error {
	isExists, err := qs.questionnaireRepository.HasByUserID(questionnaire.UserID)
	if err != nil {
		return err
	}
	if isExists {
		return repository.ErrDuplicate
	}

	return qs.questionnaireRepository.Add(ctx, questionnaire)
}

func (qs *questionnaireService) Update(questionnaire *model.Questionnaire) error {
	isExists, err := qs.questionnaireRepository.HasByUserID(questionnaire.UserID)
	if err != nil {
		return err
	}
	if !isExists {
		return repository.ErrNotFound
	}

	return qs.questionnaireRepository.Update(questionnaire)
}

func (qs *questionnaireService) UpdateName(questionnaire *model.Questionnaire, name string) error {
	if !questionnaire.FSM.Is(model.StateQuestionnaireFillingName) {
		// todo
		return errors.New("невозможно выполнить действие")
	}

	questionnaire.Name = name

	if err := questionnaire.FSM.Event(context.Background(), model.EventQuestionnaireFillBirthDate); err != nil {
		return err
	}

	return qs.questionnaireRepository.Update(questionnaire)
}
