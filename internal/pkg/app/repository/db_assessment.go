package repository

import (
	"context"
	"database/sql"
	"errors"
	"meet/internal/pkg/app/model"
	"meet/internal/pkg/app/repository/transaction"
)

const assessmentTableName = "assessments"

type AssessmentRepository interface {
	Has(usersDirection model.Direction) (bool, error)
	Get(usersDirection model.Direction) (*model.Assessment, error)
	Add(ctx context.Context, assessment *model.Assessment) error
	Remove(ctx context.Context, assessment *model.Assessment) error
}

type assessmentDBRepository struct {
	db *sql.DB
}

func NewAssessmentDBRepository(db *sql.DB) AssessmentRepository {
	return &assessmentDBRepository{db}
}

func (ar *assessmentDBRepository) Has(usersDirection model.Direction) (bool, error) {
	sb := newSelectBuilder()
	sql, args := sb.
		Select("1").
		From(assessmentTableName).
		Where(sb.Equal("from_user_id", usersDirection.FromID)).
		Where(sb.Equal("to_user_id", usersDirection.ToID)).
		Limit(1).
		Build()

	res, err := ar.db.Exec(sql, args...)
	if err != nil {
		return false, err
	}

	ra, err := res.RowsAffected()

	return ra > 0, err
}

func (ar *assessmentDBRepository) Get(usersDirection model.Direction) (*model.Assessment, error) {
	sb := newSelectBuilder()
	query, args := sb.
		Select("*").
		From(assessmentTableName).
		Where(sb.Equal("from_user_id", usersDirection.FromID)).
		Where(sb.Equal("to_user_id", usersDirection.ToID)).
		Limit(1).
		Build()

	assessment := new(model.Assessment)

	row := ar.db.QueryRow(query, args...)
	if err := row.Scan(assessment.GetFieldPointers()...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return assessment, nil
}

func (ar *assessmentDBRepository) Add(ctx context.Context, assessment *model.Assessment) error {
	assessment.BeforeAdd()

	if err := assessment.Validate(); err != nil {
		return err
	}

	ib := newInsertBuilder()
	sql, args := ib.
		InsertInto(assessmentTableName).
		Cols("created_at", "from_user_id", "to_user_id", "message", "decision").
		Values(assessment.CreatedAt, assessment.UsersDirection.FromID, assessment.UsersDirection.ToID, assessment.Message, assessment.Decision).
		SQL("RETURNING id").
		Build()

	var id int

	conn := transaction.TxOrDB(ctx, ar.db)
	row := conn.QueryRow(sql, args...)
	if err := row.Scan(&id); err != nil {
		return err
	}

	assessment.ID = id

	return nil
}

func (ar *assessmentDBRepository) Remove(ctx context.Context, assessment *model.Assessment) error {
	if err := assessment.Validate(); err != nil {
		return err
	}

	db := newDeleteBuilder()
	query, args := db.DeleteFrom(assessmentTableName).Where(db.Equal("id", assessment.ID)).Build()

	_, err := ar.db.Exec(query, args...)

	return err
}
