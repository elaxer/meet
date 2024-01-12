package service

import (
	"context"
	"database/sql"
	"errors"
	"meet/internal/pkg/app/model"
	"meet/internal/pkg/app/repository"
	"meet/internal/pkg/app/repository/transaction"
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

	qTo, err := as.questionnaireRepository.GetByUserID(assessment.UsersDirection.FromID)
	if err != nil {
		return err
	}

	if !qFrom.CheckCompatibilities(qTo, time.Now()) {
		return ErrQuestionnairesIncompatible
	}

	reversedAssessment, err := as.assessmentRepository.Get(assessment.UsersDirection.NewReversed())
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return err
	}

	if reversedAssessment == nil {
		err := as.assessmentRepository.Add(context.Background(), assessment)

		return err
	}

	ctx, tx, err := transaction.BeginTx(context.Background(), as.db)
	if err != nil {
		return err
	}

	assessment.IsMutual = true
	if err := as.assessmentRepository.Remove(ctx, reversedAssessment); err != nil {
		tx.Rollback()

		return err
	}

	couple := new(model.Couple)
	couple.UsersDirection.FromID = assessment.UsersDirection.FromID
	couple.UsersDirection.ToID = assessment.UsersDirection.ToID

	if err := as.coupleRepository.Add(ctx, couple); err != nil {
		tx.Rollback()

		return err
	}

	tx.Commit()

	return nil
}
