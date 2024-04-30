package repository

import (
	"context"
	"meet/internal/pkg/app/model"
)

type AssessmentRepository interface {
	Has(usersDirection model.Direction) (bool, error)
	Get(usersDirection model.Direction) (*model.Assessment, error)
	Add(ctx context.Context, assessment *model.Assessment) error
	Remove(ctx context.Context, assessment *model.Assessment) error
}
