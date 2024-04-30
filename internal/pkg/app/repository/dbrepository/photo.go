package dbrepository

import (
	"meet/internal/pkg/app/database"
	"meet/internal/pkg/app/model"
	"meet/internal/pkg/app/repository"
)

var photosTableName = "questionnaire_photos"

type photoRepository struct {
	conn database.Connection
}

func NewPhotoRepository(conn database.Connection) repository.PhotoRepository {
	return &photoRepository{conn}
}

func (pr *photoRepository) GetByQuestionnaireID(questionnaireID int) ([]*model.Photo, error) {
	sb := newSelectBuilder()
	query, args := sb.
		Select("*").
		From(photosTableName).
		Where(sb.Equal("questionnaire_id", questionnaireID)).
		OrderBy("created_at").
		Desc().
		Build()

	photos := make([]*model.Photo, 0)

	rows, err := pr.conn.Query(query, args...)
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

func (pr *photoRepository) Add(photo *model.Photo) error {
	photo.BeforeAdd()

	if err := photo.Validate(); err != nil {
		return err
	}

	query, args := newInsertBuilder().
		InsertInto(photosTableName).
		Cols("created_at", "questionnaire_id", "path").
		Values(photo.CreatedAt, photo.QuestionnaireID, photo.Path).
		SQL("RETURNING id").
		Build()

	var id int
	row := pr.conn.QueryRow(query, args...)
	if err := row.Scan(&id); err != nil {
		return err
	}

	photo.ID = id

	return nil
}

func (pr *photoRepository) Remove(photo *model.Photo) error {
	if err := photo.Validate(); err != nil {
		return err
	}

	db := newDeleteBuilder()
	query, args := db.DeleteFrom(photosTableName).Where(db.Equal("id", photo.ID)).Build()

	_, err := pr.conn.Exec(query, args...)

	return err
}
