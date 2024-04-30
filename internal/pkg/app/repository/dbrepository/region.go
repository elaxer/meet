package dbrepository

import (
	"context"
	"meet/internal/pkg/app/database"
	"meet/internal/pkg/app/model"
	"meet/internal/pkg/app/repository"
)

const regionTableName = "regions"

type regionRepository struct {
	conn database.Connection
}

func NewRegionRepository(conn database.Connection) repository.RegionRepository {
	return &regionRepository{conn}
}

func (rr *regionRepository) Add(ctx context.Context, region *model.Region) error {
	region.BeforeAdd()

	if err := region.Validate(); err != nil {
		return err
	}

	query, args := newInsertBuilder().
		InsertInto(regionTableName).
		Cols("id", "name").
		Values(region.ID, region.Name).
		Build()

	conn := database.TxOrDB(ctx, rr.conn)
	_, err := conn.Exec(query, args...)

	return err
}
