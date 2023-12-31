package repository

import (
	"database/sql"
	"errors"
	"meet/internal/app"

	"github.com/huandu/go-sqlbuilder"
)

var (
	ErrNotFound  = errors.New("не удалось найти модель в репозитории")
	ErrDuplicate = errors.New("в репозитории уже сущетсвует модель с таким идентификатором")
)

type RepositoryContainer struct {
	assessmentRepository    AssessmentRepository
	photoRepository         PhotoRepository
	questionnaireRepository QuestionnaireRepository
	userRepository          UserRepository
}

func NewRepositoryContainer(db *sql.DB) *RepositoryContainer {
	rc := &RepositoryContainer{
		assessmentRepository: newAssessmentRepository(db),
		photoRepository:      newPhotoRepository(db),
		userRepository:       newUserRepository(db),
	}

	rc.questionnaireRepository = newQuestionnaireRepository(db, rc.photoRepository)

	return rc
}

func (rc *RepositoryContainer) Assessment() AssessmentRepository {
	return rc.assessmentRepository
}

func (rc *RepositoryContainer) Photo() PhotoRepository {
	return rc.photoRepository
}

func (rc *RepositoryContainer) Questionnaire() QuestionnaireRepository {
	return rc.questionnaireRepository
}

func (rc *RepositoryContainer) User() UserRepository {
	return rc.userRepository
}

// Реализует интерфейс Repository
type dbRepository struct {
	db *sql.DB
}

func (r *dbRepository) createSelectBuilder() *sqlbuilder.SelectBuilder {
	sb := sqlbuilder.NewSelectBuilder()
	sb.SetFlavor(app.SQLBuilderFlavor)

	return sb
}

// TODO more createBuilder
// TODO rename createSelectBuilder
// TODO validate model after select
