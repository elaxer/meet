package service

import (
	"errors"
	"meet/internal/pkg/app/model"
	"meet/internal/pkg/app/repository"
	"time"
)

var (
	ErrAlreadyAssessed            = errors.New("оценка уже произведена")
	ErrQuestionnairesIncompatible = errors.New("анкеты несовместимы")
)

type AssessmentService struct {
	assessmentRepository    repository.AssessmentRepository
	questionnaireRepository repository.QuestionnaireRepository
}

func newAssessmentService(
	assessmentRepository repository.AssessmentRepository,
	questionnaireRepository repository.QuestionnaireRepository,
) *AssessmentService {
	return &AssessmentService{assessmentRepository, questionnaireRepository}
}

func (as *AssessmentService) Assess(assessment *model.Assessment, currentTime time.Time) error {
	if err := assessment.Validate(); err != nil {
		return err
	}

	hasAssessment, err := as.assessmentRepository.HasByDirection(assessment.UsersDirection)
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

	qTo, err := as.questionnaireRepository.GetByUserID(assessment.UsersDirection.FromID)
	if err != nil {
		return err
	}

	if !qFrom.CheckCompatibilities(qTo, currentTime) {
		return ErrQuestionnairesIncompatible
	}

	assessment.IsMutual, err = as.assessmentRepository.HasByDirection(assessment.UsersDirection.NewReversed())
	if err != nil {
		return err
	}

	err = as.assessmentRepository.Add(assessment)

	return err
}

func (as *AssessmentService) IsCouple(usersDirection model.Direction) (bool, error) {
	has, err := as.assessmentRepository.HasByDirection(usersDirection)
	if err != nil {
		return false, err
	}
	if !has {
		return false, nil
	}

	has, err = as.assessmentRepository.HasByDirection(usersDirection)

	return has, err
}
