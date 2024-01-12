package model

import (
	"time"

	"github.com/guregu/null"
)

type Model interface {
	comparable
	BeforeAdd()
	BeforeUpdate()
	Validate() error
}

// BaseModel это базовая модель, содержащая повторяющиеся поля всех моделей
type BaseModel struct {
	ID        int       `json:"id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt null.Time `json:"-"`
}

func (bm *BaseModel) BeforeAdd() {
	bm.CreatedAt = time.Now()
}

func (bm *BaseModel) BeforeUpdate() {
	bm.UpdatedAt = null.TimeFrom(time.Now())
}

// GetFieldPointers реализует интерфейс Model
func (bm *BaseModel) GetFieldPointers() []interface{} {
	return []interface{}{
		&bm.ID,
		&bm.CreatedAt,
		&bm.UpdatedAt,
	}
}

type SMEvent string
type SMState string
