package repository

import (
	"database/sql"
	"meet/internal/app"
	"meet/internal/app/model"

	"github.com/huandu/go-sqlbuilder"
)

const assessmentTableName = "assessments"

type AssessmentRepository interface {
	HasByDirection(direction model.Direction) (bool, error)
	Add(assessment *model.Assessment) error
}

type assessmentDBRepository struct {
	dbRepository
}

func newAssessmentRepository(db *sql.DB) AssessmentRepository {
	ar := new(assessmentDBRepository)
	ar.db = db

	return ar
}

func (ar *assessmentDBRepository) HasByDirection(direction model.Direction) (bool, error) {
	sb := ar.createSelectBuilder()
	sql, args := sb.
		Select("1").
		From(assessmentTableName).
		Where(sb.Equal("from_questionnaire_id", direction.FromID)).
		Where(sb.Equal("to_questionnaire_id", direction.ToID)).
		Limit(1).
		Build()

	res, err := ar.db.Exec(sql, args...)
	if err != nil {
		return false, err
	}

	ra, err := res.RowsAffected()

	return ra > 0, err
}

func (ar *assessmentDBRepository) Add(assessment *model.Assessment) error {
	if err := assessment.Validate(); err != nil {
		return err
	}

	ib := sqlbuilder.NewInsertBuilder()
	sql, args := ib.
		InsertInto(assessmentTableName).
		Cols("from_questionnaire_id", "to_questionnaire_id", "message", "decision").
		Values(assessment.Direction.FromID, assessment.Direction.ToID, assessment.Message, assessment.Decision).
		BuildWithFlavor(app.SQLBuilderFlavor)

	_, err := ar.db.Exec(sql, args...)

	return err
}
