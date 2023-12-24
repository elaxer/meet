package repository

import (
	"database/sql"
	"errors"
	"meet/internal/pkg/app"

	"github.com/huandu/go-sqlbuilder"
)

var (
	ErrNotFound  = errors.New("не удалось найти модель в репозитории")
	ErrDuplicate = errors.New("в репозитории уже сущетсвует модель с таким идентификатором")
)

type RepositoryContainer struct {
	assessmentRepository    AssessmentRepository
	messageRepository       MessageRepository
	photoRepository         PhotoRepository
	questionnaireRepository QuestionnaireRepository
	userRepository          UserRepository
}

func NewRepositoryContainer(db *sql.DB) *RepositoryContainer {
	pr := newPhotoRepository(db)
	return &RepositoryContainer{
		assessmentRepository:    newAssessmentRepository(db),
		messageRepository:       newMessageRepository(db),
		photoRepository:         pr,
		userRepository:          newUserRepository(db),
		questionnaireRepository: newQuestionnaireRepository(db, pr),
	}
}

func (rc *RepositoryContainer) Assessment() AssessmentRepository {
	return rc.assessmentRepository
}

func (rc *RepositoryContainer) Message() MessageRepository {
	return rc.messageRepository
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
