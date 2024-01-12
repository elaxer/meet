package repository

import (
	"context"
	"database/sql"
	"meet/internal/pkg/app/model"
	"meet/internal/pkg/app/repository/transaction"
)

const regionTableName = "regions"

type RegionRepository interface {
	Add(ctx context.Context, region *model.Region) error
}

type regionDBRepository struct {
	db *sql.DB
}

func NewRegionDBRepository(db *sql.DB) RegionRepository {
	return &regionDBRepository{db}
}

func (rr *regionDBRepository) Add(ctx context.Context, region *model.Region) error {
	region.BeforeAdd()

	if err := region.Validate(); err != nil {
		return err
	}

	query, args := newInsertBuilder().
		InsertInto(regionTableName).
		Cols("id", "name").
		Values(region.ID, region.Name).
		Build()

	conn := transaction.TxOrDB(ctx, rr.db)
	_, err := conn.Exec(query, args...)

	return err
}
