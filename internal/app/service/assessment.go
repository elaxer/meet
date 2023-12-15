package service

import (
	"errors"
	"meet/internal/app/model"
	"meet/internal/app/repository"

	"github.com/guregu/null"
)

type AssessmentService struct {
	assessmentRepository repository.AssessmentRepository
}

func newAssessmentService(
	assessmentRepository repository.AssessmentRepository,
) *AssessmentService {
	return &AssessmentService{assessmentRepository}
}

func (as *AssessmentService) Assess(usersDirection model.Direction, decision model.Decision, message null.String) (bool, error) {
	a := &model.Assessment{
		UsersDirection: usersDirection,
		Decision:       decision,
		Message:        message,
	}
	if err := a.Validate(); err != nil {
		return false, err
	}

	hasLike, err := as.assessmentRepository.HasByDirection(usersDirection)
	if err != nil {
		return false, err
	}
	if hasLike {
		return false, errors.New("оценка уже произведена")
	}

	if err = as.assessmentRepository.Add(a); err != nil {
		return false, err
	}

	isMutual, err := as.assessmentRepository.HasByDirection(usersDirection.NewReversed())
	if err != nil {
		return false, nil
	}

	return isMutual, nil
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
