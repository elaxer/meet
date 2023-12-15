package model

import (
	"github.com/guregu/null"
)

type Decision string

const (
	DecisionLike      Decision = "like"
	DecisionSuperlike Decision = "dislike"
)

type Assessment struct {
	BaseModel
	Direction Direction   `json:"questionnaire_direction"`
	Message   null.String `json:"message"`
	Decision  Decision    `json:"decision"`
}

// GetFieldPointers реализует интерфейс Model
func (a *Assessment) GetFieldPointers() []interface{} {
	fields := append(a.BaseModel.GetFieldPointers(), a.Direction.GetFieldPointers()...)

	return append(fields, &a.Message, &a.Decision)
}

func (a *Assessment) Validate() error {
	err := a.Direction.Validate()

	return err
}
