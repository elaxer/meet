package repository

import "meet/internal/pkg/app/model"

type questionnaireRepository struct {
	collectionRepository[*model.Questionnaire]
}

func NewQuestionnaireRepository() QuestionnaireRepository {
	return &questionnaireRepository{}
}

func (qr *questionnaireRepository) GetByUserID(userID int) (*model.Questionnaire, error) {
	return nil, nil
}

func (qr *questionnaireRepository) HasByUserID(userID int) (bool, error) {
	return false, nil
}

func (qr *questionnaireRepository) Couples(userID, limit, offset int) ([]*model.Questionnaire, error) {
	return nil, nil
}

func (qr *questionnaireRepository) PickUp(userID, limit, offset int) ([]*model.Questionnaire, error) {
	return nil, nil
}
