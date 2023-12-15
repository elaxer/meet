package model

import (
	"fmt"
	"time"

	"github.com/guregu/null"
)

type Model interface {
	BeforeUpdate()
	GetFieldPointers() []interface{}
}

type Validatable interface {
	Validate() error
}

type ValidationError struct {
	Subject      string `json:"subject"`
	Message      string `json:"message"`
	placeholders []any  `json:"-"`
}

func NewValidationError(subject, message string, placeholders ...any) *ValidationError {
	return &ValidationError{
		Subject:      subject,
		Message:      message,
		placeholders: placeholders,
	}
}

func (ve *ValidationError) Error() string {
	return fmt.Sprintf(ve.Message, ve.placeholders...)
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
