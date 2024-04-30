package service

import (
	"errors"
	"io"
	"meet/internal/pkg/app"
	"meet/internal/pkg/app/helper"
	"meet/internal/pkg/app/model"
	"meet/internal/pkg/app/repository"
	"os"
	"path/filepath"
)

const photosMax = 10

var (
	ErrPhotoUploadLimit = errors.New("превышено макисмальное количество загружаемых фотографий в анкету")
)

type PhotoService interface {
	Attach(questionnaires ...*model.Questionnaire) error
	Upload(userID int, file io.ReadSeeker) (*model.Photo, error)
	Delete(userID, photoID int) (*model.Photo, error)
}

type photoService struct {
	urlHelper               helper.URLHelper
	pathHelper              helper.PathHelper
	photoRepository         repository.PhotoRepository
	questionnaireRepository repository.QuestionnaireRepository
	fileUploaderService     FileUploaderService
}

func NewPhotoService(
	urlHelper helper.URLHelper,
	pathHelper helper.PathHelper,
	photoRepository repository.PhotoRepository,
	questionnaireRepository repository.QuestionnaireRepository,
	fileUploaderService FileUploaderService,
) PhotoService {
	return &photoService{urlHelper, pathHelper, photoRepository, questionnaireRepository, fileUploaderService}
}

func (ps *photoService) Attach(questionnaires ...*model.Questionnaire) error {
	for _, q := range questionnaires {
		photos, err := ps.photoRepository.GetByQuestionnaireID(q.ID)
		if err != nil {
			return err
		}

		for _, p := range photos {
			p.URL = ps.urlHelper.UploadURL(p.Path, app.UploadTypeImage, q.UserID)
		}

		q.Photos = append(q.Photos, photos...)
	}

	return nil
}

func (ps *photoService) Upload(userID int, fileReaderSeeker io.ReadSeeker) (*model.Photo, error) {
	q, err := ps.questionnaireRepository.GetByUserID(userID)
	if err != nil {
		return nil, err
	}
	if len(q.Photos) > photosMax {
		return nil, ErrPhotoUploadLimit
	}

	f, err := ps.fileUploaderService.Upload(fileReaderSeeker, app.UploadTypeImage, userID)
	if err != nil {
		return nil, err
	}

	_, fname := filepath.Split(f.Name())

	p := new(model.Photo)
	p.Path = fname
	p.QuestionnaireID = q.ID
	if err := ps.photoRepository.Add(p); err != nil {
		return nil, err
	}

	return p, nil
}

func (ps *photoService) Delete(userID, photoID int) (*model.Photo, error) {
	q, err := ps.questionnaireRepository.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	for _, p := range q.Photos {
		if p.ID != photoID {
			continue
		}

		if err := ps.photoRepository.Remove(p); err != nil {
			return p, err
		}

		if err := os.Remove(ps.pathHelper.UploadPath(p.Path, app.UploadTypeImage, userID)); err != nil {
			return p, err
		}

		return p, nil
	}

	return nil, repository.ErrNotFound
}
