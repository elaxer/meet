package dbrepository

import (
	"context"
	"meet/internal/pkg/app/database"
	"meet/internal/pkg/app/model"
	"meet/internal/pkg/app/repository"
)

const cityTableName = "cities"

type cityRepository struct {
	conn database.Connection
}

func NewCityRepository(conn database.Connection) repository.CityRepository {
	return &cityRepository{conn}
}

func (cr *cityRepository) GetByCountryID(countryID, limit, offset int) ([]*model.City, error) {
	sb := newSelectBuilder()
	query, args := sb.
		Select("*").
		From(cityTableName).
		Where(sb.Equal("country_id", countryID)).
		Limit(limit).
		Offset(offset).
		Build()

	var cities []*model.City

	rows, err := cr.conn.Query(query, args...)
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

func (cr *cityRepository) Add(ctx context.Context, city *model.City) error {
	city.BeforeAdd()

	if err := city.Validate(); err != nil {
		return err
	}

	query, args := newInsertBuilder().
		InsertInto(cityTableName).
		Cols("id", "country_id", "name", "latitude", "longitude").
		Values(city.ID, city.CountryID, city.Name, city.Latitude, city.Longitude).
		Build()

	conn := database.TxOrDB(ctx, cr.conn)
	_, err := conn.Exec(query, args...)

	return err
}
