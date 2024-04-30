package dbrepository

import (
	"context"
	"meet/internal/pkg/app/database"
	"meet/internal/pkg/app/model"
	"meet/internal/pkg/app/repository"
)

const countryTableName = "countries"

type countryRepository struct {
	conn database.Connection
}

func NewCountryRepository(conn database.Connection) repository.CountryRepository {
	return &countryRepository{conn}
}

func (cr *countryRepository) GetAll() ([]*model.Country, error) {
	query, args := newSelectBuilder().Select("*").From(countryTableName).Build()

	var countries []*model.Country

	rows, err := cr.conn.Query(query, args...)
	if err != nil {
		return countries, err
	}

	for rows.Next() {
		country := new(model.Country)

		if err := rows.Scan(country.GetFieldPointers()...); err != nil {
			return countries[0:0], err
		}

		countries = append(countries, country)
	}

	return countries, nil
}

func (cr *countryRepository) Add(ctx context.Context, country *model.Country) error {
	country.BeforeAdd()

	if err := country.Validate(); err != nil {
		return err
	}

	query, args := newInsertBuilder().
		InsertInto(countryTableName).
		Cols("id", "region_id", "name", "name_native", "emoji").
		Values(country.ID, country.RegionID, country.Name, country.NameNative, country.Emoji).
		Build()

	conn := database.TxOrDB(ctx, cr.conn)
	_, err := conn.Exec(query, args...)

	return err
}
