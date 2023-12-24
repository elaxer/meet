package model

import (
	"fmt"
	"time"
)

const (
	AgeMin = 18
	AgeMax = 65
)

var (
	ErrAgeMin = NewValidationError("age", "возраст должен быть не менее %d лет", AgeMin)
	ErrAgeMax = NewValidationError("age", "возраст должен быть не менее %d лет", AgeMax)

	ErrAgeRangeInvalid = NewValidationError("age_range", "минимальный возраст не может быть больше максимального возраста")
)

type BirthDate time.Time

func NewBirthDate(year int, month time.Month, day int) BirthDate {
	return BirthDate(time.Date(year, month, day, 0, 0, 0, 0, time.UTC))
}

func (bd BirthDate) Age(currentTime time.Time) int {
	birthDate := time.Time(bd)

	age := currentTime.Year() - birthDate.Year()

	if currentTime.YearDay() < birthDate.YearDay() {
		age--
	}

	return age
}

func (bd BirthDate) String() string {
	return time.Time(bd).Format("2006-01-02")
}

func (bd BirthDate) MarshalJSON() ([]byte, error) {
	j := fmt.Sprintf("\"%s\"", bd.String())

	return []byte(j), nil
}

func (bd BirthDate) Validate(currentTime time.Time) error {
	errs := &ValidationErrors{}

	if bd.Age(currentTime) < AgeMin {
		errs.Append(ErrAgeMin)
	}
	if bd.Age(currentTime) > AgeMax {
		errs.Append(ErrAgeMax)
	}

	if errs.Empty() {
		return nil
	}

	return errs
}

type AgeRange struct {
	From int `json:"from"`
	To   int `json:"to"`
}

func (ar *AgeRange) Validate() error {
	errs := &ValidationErrors{}

	if ar.From < AgeMin {
		errs.Append(ErrAgeMin)
	}
	if ar.To > AgeMax {
		errs.Append(ErrAgeMax)
	}
	if ar.From > ar.To {
		errs.Append(ErrAgeRangeInvalid)
	}

	if errs.Empty() {
		return nil
	}

	return errs
}

func (ar *AgeRange) InRange(age int) bool {
	return age >= ar.From && age <= ar.To
}
