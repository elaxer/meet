package fixture

import (
	"meet/internal/pkg/app/model"
	"time"

	"github.com/guregu/null"
)

func baseModelSeq() func() model.BaseModel {
	var id int

	return func() model.BaseModel {
		id++

		return model.BaseModel{
			ID:        id,
			CreatedAt: time.Date(2024, 1, 1, 12, 15, 55, 0, time.UTC),
			UpdatedAt: null.Time{},
		}
	}
}
