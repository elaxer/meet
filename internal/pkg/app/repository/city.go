package repository

import (
	"context"
	"meet/internal/pkg/app/model"
)

type CityRepository interface {
	GetByCountryID(countryID, limit, offset int) ([]*model.City, error)
	Add(ctx context.Context, city *model.City) error
}
