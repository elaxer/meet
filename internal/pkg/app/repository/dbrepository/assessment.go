package dbrepository

import (
	"context"
	"database/sql"
	"errors"
	"meet/internal/pkg/app/database"
	"meet/internal/pkg/app/model"
	"meet/internal/pkg/app/repository"
)

const assessmentTableName = "assessments"

type assessmentRepository struct {
	conn database.Connection
}

func NewAssessmentRepository(conn database.Connection) repository.AssessmentRepository {
	return &assessmentRepository{conn}
}

func (ar *assessmentRepository) Has(usersDirection model.Direction) (bool, error) {
	sb := newSelectBuilder()
	sql, args := sb.
		Select("1").
		From(assessmentTableName).
		Where(sb.Equal("from_user_id", usersDirection.FromID)).
		Where(sb.Equal("to_user_id", usersDirection.ToID)).
		Limit(1).
		Build()

	res, err := ar.conn.Exec(sql, args...)
	if err != nil {
		return false, err
	}

	ra, err := res.RowsAffected()

	return ra > 0, err
}

func (ar *assessmentRepository) Get(usersDirection model.Direction) (*model.Assessment, error) {
	sb := newSelectBuilder()
	query, args := sb.
		Select("*").
		From(assessmentTableName).
		Where(sb.Equal("from_user_id", usersDirection.FromID)).
		Where(sb.Equal("to_user_id", usersDirection.ToID)).
		Limit(1).
		Build()

	assessment := new(model.Assessment)

	row := ar.conn.QueryRow(query, args...)
	if err := row.Scan(assessment.GetFieldPointers()...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}

		return nil, err
	}

	return assessment, nil
}

func (ar *assessmentRepository) Add(ctx context.Context, assessment *model.Assessment) error {
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

	conn := database.TxOrDB(ctx, ar.conn)
	row := conn.QueryRow(sql, args...)
	if err := row.Scan(&id); err != nil {
		return err
	}

	assessment.ID = id

	return nil
}

func (ar *assessmentRepository) Remove(ctx context.Context, assessment *model.Assessment) error {
	if err := assessment.Validate(); err != nil {
		return err
	}

	db := newDeleteBuilder()
	query, args := db.DeleteFrom(assessmentTableName).Where(db.Equal("id", assessment.ID)).Build()

	_, err := ar.conn.Exec(query, args...)

	return err
}
