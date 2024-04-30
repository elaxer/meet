package service

import (
	"context"
	"database/sql"
	"errors"
	"meet/internal/pkg/app/database"
	"meet/internal/pkg/app/model"
	"meet/internal/pkg/app/repository"
	"time"
)

var (
	ErrAlreadyAssessed            = errors.New("оценка уже произведена")
	ErrQuestionnairesIncompatible = errors.New("анкеты несовместимы")
)

type AssessmentService interface {
	Assess(assessment *model.Assessment) error
}

type assessmentService struct {
	db                      *sql.DB
	assessmentRepository    repository.AssessmentRepository
	coupleRepository        repository.CoupleRepository
	questionnaireRepository repository.QuestionnaireRepository
}

func NewAssessmentService(
	db *sql.DB,
	assessmentRepository repository.AssessmentRepository,
	coupleRepository repository.CoupleRepository,
	questionnaireRepository repository.QuestionnaireRepository,
) AssessmentService {
	return &assessmentService{db, assessmentRepository, coupleRepository, questionnaireRepository}
}

func (as *assessmentService) Assess(assessment *model.Assessment) error {
	if err := assessment.Validate(); err != nil {
		return err
	}

	hasAssessment, err := as.assessmentRepository.Has(assessment.UsersDirection)
	if err != nil {
		return err
	}
	if hasAssessment {
		return ErrAlreadyAssessed
	}

	qFrom, err := as.questionnaireRepository.GetByUserID(assessment.UsersDirection.FromID)
	if err != nil {
		return err
	}
	if !qFrom.IsReady() {
		return ErrQuestionnaireState
	}

	qTo, err := as.questionnaireRepository.GetByUserID(assessment.UsersDirection.ToID)
	if err != nil {
		return err
	}

	if !qFrom.CheckCompatibilities(qTo, time.Now()) {
		return ErrQuestionnairesIncompatible
	}

	backAssessment, err := as.assessmentRepository.Get(assessment.UsersDirection.NewBack())
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return err
	}

	if backAssessment == nil {
		return as.assessmentRepository.Add(context.Background(), assessment)
	}

	ctx, tx, err := database.BeginTx(context.Background(), as.db)
	if err != nil {
		return err
	}

	assessment.IsMutual = true
	if err := as.assessmentRepository.Remove(ctx, backAssessment); err != nil {
		tx.Rollback()

		return err
	}

	couple := model.NewCouple(assessment.UsersDirection.FromID, assessment.UsersDirection.ToID)

	if err := as.coupleRepository.Add(ctx, couple); err != nil {
		tx.Rollback()

		return err
	}

	return tx.Commit()
}
