package service

import (
	"meet/internal/app/model"
	"meet/internal/app/repository"
)

type QuestionnaireService struct {
	questionnaireRepository repository.QuestionnaireRepository
}

func newQuestionnaireService(questionnaireRepository repository.QuestionnaireRepository) *QuestionnaireService {
	return &QuestionnaireService{questionnaireRepository: questionnaireRepository}
}

func (qs *QuestionnaireService) PickUp(userID, limit, offset int) ([]*model.Questionnaire, error) {
	questionnaires := make([]*model.Questionnaire, 0, limit)

	questionnaire, err := qs.questionnaireRepository.GetByUserID(userID)
	if err != nil {
		return questionnaires, err
	}

	questionnaires, err = qs.questionnaireRepository.PickUp(questionnaire, limit, offset)
	if err != nil {
		return questionnaires, err
	}

	compatibleQuestionnaires := make([]*model.Questionnaire, 0, len(questionnaires))
	for _, q := range questionnaires {
		if questionnaire.CheckCompatibility(q) && q.CheckCompatibility(questionnaire) {
			compatibleQuestionnaires = append(compatibleQuestionnaires, q)
		}
	}

	return compatibleQuestionnaires, nil
}

func (qs *QuestionnaireService) Add(questionnaire *model.Questionnaire) error {
	isExists, err := qs.questionnaireRepository.HasByUserID(questionnaire.UserID)
	if err != nil {
		return err
	}
	if isExists {
		return repository.ErrDuplicate
	}

	return qs.questionnaireRepository.Add(questionnaire)
}

func (qs *QuestionnaireService) Update(questionnaire *model.Questionnaire) error {
	isExists, err := qs.questionnaireRepository.HasByUserID(questionnaire.UserID)
	if err != nil {
		return err
	}
	if !isExists {
		return repository.ErrNotFound
	}

	return qs.questionnaireRepository.Update(questionnaire)
}
