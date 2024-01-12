package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/guregu/null"
)

const (
	AgeMin = 18
	AgeMax = 65
)

var (
	errAgeMin = NewValidationError("age", "возраст должен быть не менее %d лет", AgeMin)
	errAgeMax = NewValidationError("age", "возраст должен быть не менее %d лет", AgeMax)

	errAgeRangeInvalid = NewValidationError("age_range", "минимальный возраст не может быть больше максимального возраста")
)

// todo
type BirthDate struct{ null.Time }

func NewBirthDate(year int, month time.Month, day int) BirthDate {
	return BirthDateFrom(time.Date(year, month, day, 0, 0, 0, 0, time.UTC))
}

func BirthDateFrom(date time.Time) BirthDate {
	return BirthDate{Time: null.TimeFrom(date)}
}

func (bd BirthDate) Age(currentTime time.Time) int {
	if !bd.Valid {
		return 0
	}

	birthDate := bd.Time

	age := currentTime.Year() - birthDate.Time.Year()

	if currentTime.YearDay() < birthDate.Time.YearDay() {
		age--
	}

	return age
}

func (bd BirthDate) String() string {
	return time.Time(bd.Time.Time).Format("2006-01-02")
}

func (bd BirthDate) MarshalJSON() ([]byte, error) {
	j := fmt.Sprintf("\"%s\"", bd.String())

	return []byte(j), nil
}

func (bd *BirthDate) UnmarshalJSON(data []byte) error {
	var bdStr string
	err := json.Unmarshal(data, &bdStr)
	if err != nil {
		return err
	}

	bdTime, err := time.Parse("2006-01-02", bdStr)
	if err != nil {
		return err
	}

	*bd = BirthDateFrom(bdTime)

	return nil
}

func (bd BirthDate) Validate(currentTime time.Time) error {
	errs := &ValidationErrors{}

	age := bd.Age(currentTime)

	if age < AgeMin {
		errs.Append(errAgeMin)
	}
	if age > AgeMax {
		errs.Append(errAgeMax)
	}

	if errs.Empty() {
		return nil
	}

	return errs
}

func (bd BirthDate) Value() (driver.Value, error) {
	return bd.Time, nil
}

func (b *BirthDate) Scan(value interface{}) error {
	switch v := value.(type) {
	case time.Time:
		*b = BirthDateFrom(v)
	case nil:
		b = nil
	default:
		return fmt.Errorf("неподдерживаемый тип для сканирования: %T", value)
	}

	return nil
}

type AgeRange struct {
	Min null.Int `json:"min"`
	Max null.Int `json:"max"`
}

func NewAgeRange(min, max int) AgeRange {
	return AgeRange{null.IntFrom(int64(min)), null.IntFrom(int64(max))}
}

func (ar *AgeRange) Validate() error {
	errs := &ValidationErrors{}

	if ar.Min.Valid && ar.Min.Int64 < AgeMin {
		errs.Append(errAgeMin)
	}

	if ar.Max.Valid && ar.Max.Int64 > AgeMax {
		errs.Append(errAgeMin)
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
