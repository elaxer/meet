package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/guregu/null"
)

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

	if age < ageMin {
		errs.Append(errAgeMin)
	}
	if age > ageMax {
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
