package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"meet/internal/app"
	"meet/internal/app/model"

	"github.com/huandu/go-sqlbuilder"
)

// questionnaireTableName represents the name of the questionnaires table
const questionnaireTableName = "questionnaires"

// QuestionnaireRepository is repository for model.Questionnaire model
type QuestionnaireRepository interface {
	GetByUserID(userID int) (*model.Questionnaire, error)
	HasByUserID(userID int) (bool, error)
	GetCouples(userID int, limit, offset int) ([]*model.Questionnaire, error)
	PickUp(questionnaire *model.Questionnaire, limit, offset int) ([]*model.Questionnaire, error)
	Add(questionnaire *model.Questionnaire) error
	Update(questionnaire *model.Questionnaire) error
}

type questionnaireDBRepository struct {
	dbRepository
	photoRepository PhotoRepository
}

// newQuestionnaireRepository is the constructor of the questionnaireDBRepository
func newQuestionnaireRepository(db *sql.DB, photoRepository PhotoRepository) QuestionnaireRepository {
	qdr := new(questionnaireDBRepository)
	qdr.db = db
	qdr.photoRepository = photoRepository

	return qdr
}

// GetByUserID implements QuestionnaireRepository interface
func (qr *questionnaireDBRepository) GetByUserID(userID int) (*model.Questionnaire, error) {
	sb := qr.createSelectBuilder()
	query, args := sb.Select("*").
		From(questionnaireTableName).
		Where(sb.Equal("user_id", userID)).
		Limit(1).
		Build()

	q := model.NewQuestionnaire()

	row := qr.db.QueryRow(query, args...)
	err := row.Scan(q.GetFieldPointers()...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	q.Photos, err = qr.photoRepository.GetByQuestionnaireID(q.ID)
	if err != nil {
		return nil, err
	}

	return q, err
}

func (qr *questionnaireDBRepository) HasByUserID(userID int) (bool, error) {
	sb := qr.createSelectBuilder()
	sql, args := sb.Select("1").
		From(questionnaireTableName).
		Where(sb.Equal("user_id", userID)).
		Limit(1).
		Build()

	res, err := qr.db.Exec(sql, args...)
	if err != nil {
		return false, err
	}

	ra, err := res.RowsAffected()

	return ra > 0, err
}

// GetCouples реализует интерфейс UserRepository
// GetCouples возвращает список анкет, с которыми пользователь состоит в паре
func (qr *questionnaireDBRepository) GetCouples(questionnaireID int, limit, offset int) ([]*model.Questionnaire, error) {
	sb := qr.createSelectBuilder()
	query, args := sb.
		Select(fmt.Sprintf("%s.*", questionnaireTableName)).
		From(sb.As("assessments", "a1")).
		Join(sb.As("assessments", "a2"), sb.Equal("a2.to_questionnaire_id", questionnaireID), "a1.to_questionnaire_id = a2.from_questionnaire_id").
		Join(questionnaireTableName, "questionnaires.id = a1.to_questionnaire_id").
		Join(userTableName, "users.id = questionnaires.user_id").
		Where("users.is_blocked = false").
		Limit(limit).
		Offset(offset).
		Build()

	questionnaires := make([]*model.Questionnaire, 0, limit)

	rows, err := qr.db.Query(query, args...)
	if err != nil {
		return questionnaires, err
	}
	defer rows.Close()

	for rows.Next() {
		q := model.NewQuestionnaire()

		fields := q.GetFieldPointers()
		err := rows.Scan(fields...)
		if err != nil {
			return questionnaires[0:0], err
		}

		q.Photos, err = qr.photoRepository.GetByQuestionnaireID(q.ID)
		if err != nil {
			return questionnaires[0:0], err
		}

		questionnaires = append(questionnaires, q)
	}

	return questionnaires, nil
}

// PickUpCouples implements UserRepository interface
func (qr *questionnaireDBRepository) PickUp(questionnaire *model.Questionnaire, limit, offset int) ([]*model.Questionnaire, error) {
	questionnaires := make([]*model.Questionnaire, 0, limit)

	livingCondition, err := sqlbuilder.PostgreSQL.Interpolate(
		"ORDER BY CASE WHEN country = $1 THEN 0 ELSE 1 END, CASE WHEN city = $1 THEN 0 ELSE 1 END",
		[]interface{}{questionnaire.Country, questionnaire.City},
	)
	if err != nil {
		return questionnaires, err
	}

	sb := qr.createSelectBuilder()
	sbAssessments := qr.createSelectBuilder()

	query, args := sb.
		Select("questionnaires.*").
		From(questionnaireTableName).
		Join(userTableName, "users.id = questionnaires.user_id").
		Where("users.is_blocked = false").
		Where("questionnaires.is_active = true").
		Where(sb.Between("age", questionnaire.AgeRange.From, questionnaire.AgeRange.To)).
		Where(sb.NotEqual("questionnaires.id", questionnaire.ID)).
		Where(sb.NotIn(
			"questionnaires.id",
			sbAssessments.
				Select("to_questionnaire_id").
				From(assessmentTableName).
				Where(sbAssessments.Equal("assessments.from_questionnaire_id", questionnaire.ID)),
		)).
		SQL(livingCondition).
		Limit(limit).
		Offset(offset).
		Build()

	rows, err := qr.db.Query(query, args...)
	if err != nil {
		return questionnaires, err
	}
	defer rows.Close()

	for rows.Next() {
		q := new(model.Questionnaire)

		err := rows.Scan(q.GetFieldPointers()...)
		if err != nil {
			return questionnaires[0:0], err
		}

		q.Photos, err = qr.photoRepository.GetByQuestionnaireID(q.ID)
		if err != nil {
			return questionnaires[0:0], err
		}

		questionnaires = append(questionnaires, q)
	}

	return questionnaires, nil
}

func (qr *questionnaireDBRepository) Add(questionnaire *model.Questionnaire) error {
	if err := questionnaire.Validate(); err != nil {
		return err
	}

	ib := sqlbuilder.NewInsertBuilder()
	sql, args := ib.InsertInto(questionnaireTableName).
		Cols(
			"user_id",
			"name",
			"age",
			"gender",
			"orientation",
			"meeting_purpose",
			"age_range_from",
			"age_range_to",
			"country",
			"city",
			"about",
		).
		Values(
			questionnaire.UserID,
			questionnaire.Name,
			questionnaire.Age,
			questionnaire.Gender,
			questionnaire.Orientation,
			questionnaire.MeetingPurpose,
			questionnaire.AgeRange.From,
			questionnaire.AgeRange.To,
			questionnaire.Country,
			questionnaire.City,
			questionnaire.About,
		).
		BuildWithFlavor(app.SQLBuilderFlavor)

	_, err := qr.db.Exec(sql, args...)

	return err
}

func (qr *questionnaireDBRepository) Update(questionnaire *model.Questionnaire) error {
	if err := questionnaire.Validate(); err != nil {
		return err
	}

	questionnaire.BeforeUpdate()

	ub := sqlbuilder.NewUpdateBuilder()
	sql, args := ub.Update(questionnaireTableName).
		Set(
			ub.Assign("name", questionnaire.Name),
			ub.Assign("age", questionnaire.Age),
			ub.Assign("gender", questionnaire.Gender),
			ub.Assign("orientation", questionnaire.Orientation),
			ub.Assign("meeting_purpose", questionnaire.MeetingPurpose),
			ub.Assign("age_range_from", questionnaire.AgeRange.From),
			ub.Assign("age_range_to", questionnaire.AgeRange.To),
			ub.Assign("country", questionnaire.Country),
			ub.Assign("city", questionnaire.City),
			ub.Assign("about", questionnaire.About),
			ub.Assign("is_active", questionnaire.IsActive),
		).
		Where(ub.Equal("user_id", questionnaire.UserID)).
		BuildWithFlavor(app.SQLBuilderFlavor)

	_, err := qr.db.Exec(sql, args...)

	return err
}
