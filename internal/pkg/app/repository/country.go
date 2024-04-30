package repository

import (
	"context"
	"meet/internal/pkg/app/model"
)

type CountryRepository interface {
	GetAll() ([]*model.Country, error)
	Add(ctx context.Context, country *model.Country) error
}
