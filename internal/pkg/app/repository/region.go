package repository

import (
	"context"
	"meet/internal/pkg/app/model"
)

type RegionRepository interface {
	Add(ctx context.Context, region *model.Region) error
}
