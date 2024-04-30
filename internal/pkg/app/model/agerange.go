package model

import (
	"github.com/guregu/null"
)

var (
	errAgeRangeInvalid = NewValidationError("age_range", "минимальный возраст не может быть больше максимального возраста")
)

type AgeRange struct {
	Min null.Int `json:"min"`
	Max null.Int `json:"max"`
}

func NewAgeRange(min, max int) AgeRange {
	return AgeRange{null.IntFrom(int64(min)), null.IntFrom(int64(max))}
}

func (ar *AgeRange) Validate() error {
	errs := &ValidationErrors{}

	if ar.Min.Valid && ar.Min.Int64 < ageMin {
		errs.Append(errAgeMin)
	}

	if ar.Max.Valid && ar.Max.Int64 > ageMax {
		errs.Append(errAgeMax)
	}

	if ar.Min.Valid && ar.Max.Valid && ar.Min.Int64 > ar.Max.Int64 {
		errs.Append(errAgeRangeInvalid)
	}

	if errs.Empty() {
		return nil
	}

	return errs
}

func (ar *AgeRange) InRange(age int) bool {
	return age >= int(ar.Min.Int64) && age <= int(ar.Max.Int64)
}
