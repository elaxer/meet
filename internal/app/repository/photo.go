package repository

import (
	"database/sql"
	"errors"
	"meet/internal/app"
	"meet/internal/app/model"

	"github.com/huandu/go-sqlbuilder"
)

var photosTableName = "questionnaire_photos"

type PhotoRepository interface {
	Get(id int) (*model.Photo, error)
	GetByQuestionnaireID(questionnaireID int) ([]*model.Photo, error)
	Add(photo *model.Photo) error
	Remove(photo *model.Photo) error
}

type photoDBRepository struct {
	dbRepository
}

func newPhotoRepository(db *sql.DB) PhotoRepository {
	pr := new(photoDBRepository)
	pr.db = db

	return pr
}

func (pr *photoDBRepository) Get(id int) (*model.Photo, error) {
	sb := pr.dbRepository.createSelectBuilder()
	query, args := sb.Select("*").From(photosTableName).Where(sb.Equal("id", id)).Limit(1).Build()

	p := new(model.Photo)
	row := pr.db.QueryRow(query, args...)
	err := row.Scan(p.GetFieldPointers()...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
	}

	return p, err
}

func (pr *photoDBRepository) GetByQuestionnaireID(questionnaireID int) ([]*model.Photo, error) {
	sb := pr.dbRepository.createSelectBuilder()
	query, args := sb.
		Select("*").
		From(photosTableName).
		Where(sb.Equal("questionnaire_id", questionnaireID)).
		OrderBy("created_at").
		Desc().
		Build()

	photos := make([]*model.Photo, 0)

	rows, err := pr.db.Query(query, args...)
	if err != nil {
		return photos, err
	}
	defer rows.Close()

	for rows.Next() {
		p := new(model.Photo)

		if err := rows.Scan(p.GetFieldPointers()...); err != nil {
			return photos[0:0], err
		}

		photos = append(photos, p)
	}

	return photos, nil
}

func (pr *photoDBRepository) Add(photo *model.Photo) error {
	photo.BeforeAdd()

	if err := photo.Validate(); err != nil {
		return err
	}

	ib := sqlbuilder.NewInsertBuilder()
	query, args := ib.
		InsertInto(photosTableName).
		Cols("created_at", "questionnaire_id", "path").
		Values(photo.CreatedAt, photo.QuestionnaireID, photo.Path).
		BuildWithFlavor(app.SQLBuilderFlavor)

	_, err := pr.db.Exec(query, args...)

	return err
}

func (pr *photoDBRepository) Remove(photo *model.Photo) error {
	if err := photo.Validate(); err != nil {
		return err
	}

	db := sqlbuilder.NewDeleteBuilder()
	query, args := db.DeleteFrom(photosTableName).Where(db.Equal("id", photo.ID)).BuildWithFlavor(app.SQLBuilderFlavor)

	_, err := pr.db.Exec(query, args...)

	return err
}
