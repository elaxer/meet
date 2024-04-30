package repository

import (
	"context"
	"meet/internal/pkg/app/model"
)

type CoupleRepository interface {
	Has(usersDirection model.Direction) (bool, error)
	Add(ctx context.Context, couple *model.Couple) error
}
