package repository

import (
	"context"
	"database/sql"
	"meet/internal/pkg/app/model"
	"meet/internal/pkg/app/repository/transaction"
)

const countryTableName = "countries"

type CountryRepository interface {
	GetAll() ([]*model.Country, error)
	Add(ctx context.Context, country *model.Country) error
}

type countryDBRepository struct {
	db *sql.DB
}

func NewCountryDBRepository(db *sql.DB) CountryRepository {
	return &countryDBRepository{db}
}

func (cr *countryDBRepository) GetAll() ([]*model.Country, error) {
	query, args := newSelectBuilder().Select("*").From(countryTableName).Build()

	var countries []*model.Country

	rows, err := cr.db.Query(query, args...)
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

func (cr *countryDBRepository) Add(ctx context.Context, country *model.Country) error {
	country.BeforeAdd()

	if err := country.Validate(); err != nil {
		return err
	}

	query, args := newInsertBuilder().
		InsertInto(countryTableName).
		Cols("id", "region_id", "name", "name_native", "emoji").
		Values(country.ID, country.RegionID, country.Name, country.NameNative, country.Emoji).
		Build()

	conn := transaction.TxOrDB(ctx, cr.db)
	_, err := conn.Exec(query, args...)

	return err
}
