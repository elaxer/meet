package repository

import (
	"context"
	"meet/internal/pkg/app/model"
)

// QuestionnaireRepository is repository for model.Questionnaire model
type QuestionnaireRepository interface {
	GetByUserID(userID int) (*model.Questionnaire, error)
	HasByUserID(userID int) (bool, error)
	Couples(userID, limit, offset int) ([]*model.Questionnaire, error)
	Suggested(userID, limit, offset int) ([]*model.Questionnaire, error)
	Assessed(userID, limit, offset int) ([]*model.Questionnaire, error)
	Add(ctx context.Context, questionnaire *model.Questionnaire) error
	Update(questionnaire *model.Questionnaire) error
}
