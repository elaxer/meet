package model

import (
	"strings"

	"github.com/guregu/null"
)

type Decision string

const (
	DecisionLike      Decision = "like"
	DecisionSuperlike Decision = "dislike"
)

const assessmentMessageLenMax = 2048

var (
	errAssessmentMessageTooLong = NewValidationError("message", "длина сообщения не должна превышать %d символов", assessmentMessageLenMax)
)

type Assessment struct {
	BaseModel
	UsersDirection Direction   `json:"users_direction"`
	Message        null.String `json:"message"`
	Decision       Decision    `json:"decision"`
	IsMutual       bool        `json:"is_mutual"`
}

// GetFieldPointers реализует интерфейс Model
func (a *Assessment) GetFieldPointers() []interface{} {
	fields := append(a.BaseModel.GetFieldPointers(), a.UsersDirection.GetFieldPointers()...)

	return append(fields, &a.Message, &a.Decision)
}

func (a *Assessment) BeforeAdd() {
	a.BaseModel.BeforeAdd()

	if a.Message.Valid {
		a.Message.String = strings.TrimSpace(a.Message.String)
	}
}

func (a *Assessment) Validate() error {
	errs := &ValidationErrors{}

	if err := a.UsersDirection.Validate(); err != nil {
		errs.Append(err)
	}
	if a.Message.Valid && len(strings.TrimSpace(a.Message.String)) > messageLenMax {
		errs.Append(errAssessmentMessageTooLong)
	}

	if errs.Empty() {
		return nil
	}

	return errs
}
