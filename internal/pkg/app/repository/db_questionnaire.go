package repository

import (
	"context"
	"database/sql"
	"errors"
	"meet/internal/pkg/app/model"
	"meet/internal/pkg/app/repository/transaction"

	"github.com/huandu/go-sqlbuilder"
)

// questionnaireTableName represents the name of the questionnaires table
const questionnaireTableName = "questionnaires"

// QuestionnaireRepository is repository for model.Questionnaire model
type QuestionnaireRepository interface {
	GetByUserID(userID int) (*model.Questionnaire, error)
	HasByUserID(userID int) (bool, error)
	Couples(userID, limit, offset int) ([]*model.Questionnaire, error)
	PickUp(userID, limit, offset int) ([]*model.Questionnaire, error)
	Add(ctx context.Context, questionnaire *model.Questionnaire) error
	Update(questionnaire *model.Questionnaire) error
}

type questionnaireDBRepository struct {
	db              *sql.DB
	photoRepository PhotoRepository
}

// newQuestionnaireRepository is the constructor of the questionnaireDBRepository
func NewQuestionnaireDBRepository(db *sql.DB, photoRepository PhotoRepository) QuestionnaireRepository {
	return &questionnaireDBRepository{db, photoRepository}
}

// GetByUserID implements QuestionnaireRepository interface
func (qr *questionnaireDBRepository) GetByUserID(userID int) (*model.Questionnaire, error) {
	sb := newSelectBuilder()
	query, args := sb.Select("*").
		From(questionnaireTableName).
		Where(sb.Equal("user_id", userID)).
		Limit(1).
		Build()

	q := model.NewQuestionnaireEmpty()

	row := qr.db.QueryRow(query, args...)
	err := row.Scan(q.GetFieldPointers()...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	q.FSM.SetState(q.State())
	q.Photos, err = qr.photoRepository.GetByQuestionnaireID(q.ID)
	if err != nil {
		return nil, err
	}

	return q, err
}

func (qr *questionnaireDBRepository) HasByUserID(userID int) (bool, error) {
	sb := newSelectBuilder()
	query, args := sb.Select("1").
		From(questionnaireTableName).
		Where(sb.Equal("user_id", userID)).
		Limit(1).
		Build()

	res, err := qr.db.Exec(query, args...)
	if err != nil {
		return false, err
	}

	ra, err := res.RowsAffected()

	return ra > 0, err
}

// Couples реализует интерфейс UserRepository
// Couples возвращает список анкет, с которыми пользователь состоит в паре
func (qr *questionnaireDBRepository) Couples(userID, limit, offset int) ([]*model.Questionnaire, error) {
	sb := newSelectBuilder()
	query, args := sb.
		Select("q.*").
		From(coupleTableName).
		JoinWithOption(
			sqlbuilder.LeftJoin,
			sb.As(questionnaireTableName, "q"),
			sb.Or(
				sb.And("user_id = from_user_id", sb.Equal("to_user_id", userID)),
				sb.And("user_id = to_user_id", sb.Equal("from_user_id", userID)),
			),
		).
		JoinWithOption(sqlbuilder.LeftJoin, sb.As(userTableName, "u"), "u.id = q.user_id").
		Where(
			sb.Or(sb.Equal("from_user_id", userID), sb.Equal("to_user_id", userID)),
		).
		Where(sb.NotEqual("user_id", userID)).
		Where("u.is_blocked = false").
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
		q := new(model.Questionnaire)

		err := rows.Scan(q.GetFieldPointers()...)
		if err != nil {
			return questionnaires[0:0], err
		}

		q.FSM.SetState(q.State())
		q.Photos, err = qr.photoRepository.GetByQuestionnaireID(q.ID)
		if err != nil {
			return questionnaires[0:0], err
		}

		questionnaires = append(questionnaires, q)
	}

	return questionnaires, nil
}

func (qr *questionnaireDBRepository) PickUp(userID, limit, offset int) ([]*model.Questionnaire, error) {
	sb := newSelectBuilder()
	query, args := sb.
		Select("q.*").
		From(sb.As(questionnaireTableName, "q")).
		JoinWithOption(sqlbuilder.LeftJoin, sb.As(coupleTableName, "c1"), "q.user_id = c1.from_user_id", sb.Equal("c1.to_user_id", userID)).
		JoinWithOption(sqlbuilder.LeftJoin, sb.As(coupleTableName, "c2"), "q.user_id = c2.to_user_id", sb.Equal("c1.from_user_id", userID)).
		Join(sb.As(userTableName, "u"), "u.id = q.user_id").
		Where(sb.NotEqual("q.user_id", userID)).
		Where(sb.IsNull("c1.id")).
		Where(sb.IsNull("c2.id")).
		Where("u.is_blocked = false").
		Build()

	questionnaires := make([]*model.Questionnaire, 0, limit)

	rows, err := qr.db.Query(query, args...)
	if err != nil {
		return questionnaires, err
	}
	defer rows.Close()

	for rows.Next() {
		q := new(model.Questionnaire)

		err := rows.Scan(q.GetFieldPointers())
		if err != nil {
			return questionnaires[0:0], err
		}

		q.FSM.SetState(q.State())
		q.Photos, err = qr.photoRepository.GetByQuestionnaireID(q.ID)
		if err != nil {
			return questionnaires[0:0], err
		}

		questionnaires = append(questionnaires, q)
	}

	return questionnaires, nil
}

func (qr *questionnaireDBRepository) Add(ctx context.Context, questionnaire *model.Questionnaire) error {
	questionnaire.BeforeAdd()

	if err := questionnaire.Validate(); err != nil {
		return err
	}

	query, args := newInsertBuilder().InsertInto(questionnaireTableName).
		Cols(
			"created_at",
			"user_id",
			"name",
			"birth_date",
			"gender",
			"orientation",
			"meeting_purpose",
			"age_range_min",
			"age_range_max",
			"city_id",
			"about",
			"is_active",
			"state",
		).
		Values(
			questionnaire.CreatedAt,
			questionnaire.UserID,
			questionnaire.Name,
			questionnaire.BirthDate.Time,
			questionnaire.Gender,
			questionnaire.Orientation,
			questionnaire.MeetingPurpose,
			questionnaire.AgeRange.Min,
			questionnaire.AgeRange.Max,
			questionnaire.CityID,
			questionnaire.About,
			questionnaire.IsActive,
			questionnaire.FSM.Current(),
		).
		SQL("RETURNING id").
		Build()

	var id int

	conn := transaction.TxOrDB(ctx, qr.db)
	row := conn.QueryRow(query, args...)
	if err := row.Scan(&id); err != nil {
		return err
	}

	questionnaire.ID = id

	return nil
}

func (qr *questionnaireDBRepository) Update(questionnaire *model.Questionnaire) error {
	questionnaire.BeforeUpdate()

	if err := questionnaire.Validate(); err != nil {
		return err
	}

	ub := newUpdateBuilder()
	query, args := ub.Update(questionnaireTableName).
		Set(
			ub.Assign("updated_at", questionnaire.UpdatedAt),
			ub.Assign("name", questionnaire.Name),
			ub.Assign("birth_date", questionnaire.BirthDate.Time),
			ub.Assign("gender", questionnaire.Gender),
			ub.Assign("orientation", questionnaire.Orientation),
			ub.Assign("meeting_purpose", questionnaire.MeetingPurpose),
			ub.Assign("age_range_min", questionnaire.AgeRange.Min),
			ub.Assign("age_range_max", questionnaire.AgeRange.Max),
			ub.Assign("city_id", questionnaire.CityID),
			ub.Assign("about", questionnaire.About),
			ub.Assign("is_active", questionnaire.IsActive),
			ub.Assign("state", questionnaire.FSM.Current()),
		).
		Where(ub.Equal("user_id", questionnaire.UserID)).
		Build()

	_, err := qr.db.Exec(query, args...)

	return err
}
