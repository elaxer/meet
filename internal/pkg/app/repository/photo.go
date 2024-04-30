package repository

import "meet/internal/pkg/app/model"

type PhotoRepository interface {
	GetByQuestionnaireID(questionnaireID int) ([]*model.Photo, error)
	Add(photo *model.Photo) error
	Remove(photo *model.Photo) error
}
