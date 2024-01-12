package model

import "strings"

var (
	errRegionIDEmpty   = NewValidationError("id", "идентификатор региона не может быть пустым")
	errRegionNameEmpty = NewValidationError("name", "название региона не может быть пустым")
)

type Region struct {
	ID   int    `json:"id,string"`
	Name string `json:"name"`
}

func (r *Region) BeforeAdd() {
	r.Name = strings.TrimSpace(r.Name)
}

func (r *Region) Validate() error {
	errors := &ValidationErrors{}

	if r.ID == 0 {
		errors.Append(errRegionIDEmpty)
	}
	if strings.TrimSpace(r.Name) == "" {
		errors.Append(errRegionNameEmpty)
	}

	if errors.Empty() {
		return nil
	}

	return errors
}
