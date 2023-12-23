package service

import (
	"errors"
	"io"
	"meet/internal/app"
	"meet/internal/app/model"
	"meet/internal/app/repository"
	"os"
)

const dirUploadPhoto = "photos"
const photosMax = 10

var ErrPhotoUploadLimit = errors.New("превышено макисмальное количество загружаемых фотографий в анкету")

type PhotoService struct {
	config                  *app.Config
	photoRepository         repository.PhotoRepository
	questionnaireRepository repository.QuestionnaireRepository
	fileService             *FileService
}

func newPhotoService(
	config *app.Config,
	photoRepository repository.PhotoRepository,
	questionnaireRepository repository.QuestionnaireRepository,
	fileService *FileService,
) *PhotoService {
	return &PhotoService{
		config:                  config,
		photoRepository:         photoRepository,
		questionnaireRepository: questionnaireRepository,
		fileService:             fileService,
	}
}

func (ps *PhotoService) Upload(userID int, file io.ReadSeeker) (*model.Photo, error) {
	q, err := ps.questionnaireRepository.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	if len(q.Photos) > photosMax {
		return nil, ErrPhotoUploadLimit
	}

	_, fname, err := ps.fileService.Upload(file, []string{"image"}, dirUploadPhoto)
	if err != nil {
		return nil, err
	}

	p := new(model.Photo)
	p.Path = fname
	p.QuestionnaireID = q.ID
	if err := ps.photoRepository.Add(p); err != nil {
		return nil, err
	}

	return p, nil
}

func (ps *PhotoService) GetPath(userID int, photoID int) (string, error) {
	q, err := ps.questionnaireRepository.GetByUserID(userID)
	if err != nil {
		return "", err
	}

	for _, p := range q.Photos {
		if p.ID == photoID {
			return ps.fileService.FullPath(p.Path, dirUploadPhoto), nil
		}
	}

	return "", repository.ErrNotFound
}

func (ps *PhotoService) Delete(userID int, photoID int) (*model.Photo, error) {
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

		if err := os.Remove(ps.fileService.FullPath(p.Path, dirUploadPhoto)); err != nil {
			return p, err
		}

		return p, nil
	}

	return nil, repository.ErrNotFound
}
