package service

import (
	"errors"
	"meet/internal/app/model"
	"meet/internal/app/repository"

	"github.com/guregu/null"
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

func (as *AssessmentService) Assess(userID, questionnaireID int, decision model.Decision, message null.String) (bool, error) {
	q, err := as.questionnaireRepository.GetByUserID(userID)
	if err != nil {
		return false, err
	}

	a := &model.Assessment{
		Direction: model.Direction{FromID: q.ID, ToID: questionnaireID},
		Decision:  decision,
		Message:   message,
	}
	if err := a.Validate(); err != nil {
		return false, err
	}

	hasLike, err := as.assessmentRepository.HasByDirection(a.Direction)
	if err != nil {
		return false, err
	}
	if hasLike {
		return false, errors.New("оценка уже произведена")
	}

	if err = as.assessmentRepository.Add(a); err != nil {
		return false, err
	}

	isMutual, err := as.assessmentRepository.HasByDirection(a.Direction.NewReversed())
	if err != nil {
		return false, nil
	}

	return isMutual, nil
}
