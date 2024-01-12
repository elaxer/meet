package repository

import (
	"context"
	"database/sql"
	"meet/internal/pkg/app/model"
	"meet/internal/pkg/app/repository/transaction"
)

const cityTableName = "cities"

type CityRepository interface {
	GetByCountryID(countryID, limit, offset int) ([]*model.City, error)
	Add(ctx context.Context, city *model.City) error
}

type cityDBRepository struct {
	db *sql.DB
}

func NewCityDBRepository(db *sql.DB) CityRepository {
	return &cityDBRepository{db}
}

func (cr *cityDBRepository) GetByCountryID(countryID, limit, offset int) ([]*model.City, error) {
	sb := newSelectBuilder()
	query, args := sb.
		Select("*").
		From(cityTableName).
		Where(sb.Equal("country_id", countryID)).
		Limit(limit).
		Offset(offset).
		Build()

	var cities []*model.City

	rows, err := cr.db.Query(query, args...)
	if err != nil {
		return cities, err
	}

	for rows.Next() {
		city := new(model.City)

		if err := rows.Scan(city.GetFieldPointers()...); err != nil {
			return cities[0:0], err
		}

		cities = append(cities, city)
	}

	return cities, nil
}

func (cr *cityDBRepository) Add(ctx context.Context, city *model.City) error {
	city.BeforeAdd()

	if err := city.Validate(); err != nil {
		return err
	}

	query, args := newInsertBuilder().
		InsertInto(cityTableName).
		Cols("id", "country_id", "name", "latitude", "longitude").
		Values(city.ID, city.CountryID, city.Name, city.Latitude, city.Longitude).
		Build()

	conn := transaction.TxOrDB(ctx, cr.db)
	_, err := conn.Exec(query, args...)

	return err
}
