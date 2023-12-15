package model

import (
	"errors"
	"fmt"
)

type Gender bool
type Orientation string
type MeetingPurpose string

const (
	GenderMale   Gender = false
	GenderFemale Gender = true
)

const (
	OrientationHetero Orientation = "hetero"
	OrientationHomo   Orientation = "homo"
	OrientationBi     Orientation = "bi"
)

const (
	MeetingPurposeFriendship   MeetingPurpose = "friendship"
	MeetingPurposeRelationship MeetingPurpose = "relationship"
	MeetingPurposeSex          MeetingPurpose = "sex"
)

const (
	AgeMin = 18
	AgeMax = 65
)

type AgeRange struct {
	From int `json:"from"`
	To   int `json:"to"`
}

func (ar *AgeRange) Validate() error {
	if ar.From < AgeMin {
		return fmt.Errorf("начальный возраст должен быть не менее %d лет", AgeMin)
	}
	if ar.To > AgeMax {
		return fmt.Errorf("максимальный возраст должен быть не более %d лет", AgeMax)
	}

	if ar.From > ar.To {
		return errors.New("минимальный возраст не может быть больше максимального возраста")
	}

	return nil
}

func (ar *AgeRange) InRange(age int) bool {
	return age >= ar.From && age <= ar.To
}

type Questionnaire struct {
	BaseModel
	UserID         int            `json:"user_id"`
	Name           string         `json:"name"`
	Age            int            `json:"age"`
	Gender         Gender         `json:"gender"`
	Orientation    Orientation    `json:"orientation"`
	MeetingPurpose MeetingPurpose `json:"meeting_purpose"`
	AgeRange       AgeRange       `json:"age_range"`
	Country        string         `json:"country"`
	City           string         `json:"city"`
	About          string         `json:"about"`
	Photos         []*Photo       `json:"photos"`
	IsActive       bool           `json:"is_active"`
}

func NewQuestionnaire() *Questionnaire {
	return new(Questionnaire)
}

// GetFieldPointers implements Model interface
func (q *Questionnaire) GetFieldPointers() []interface{} {
	return append(
		q.BaseModel.GetFieldPointers(),
		&q.UserID,
		&q.Name,
		&q.Age,
		&q.Gender,
		&q.Orientation,
		&q.MeetingPurpose,
		&q.AgeRange.From,
		&q.AgeRange.To,
		&q.Country,
		&q.City,
		&q.About,
		&q.IsActive,
	)
}

func (q *Questionnaire) Validate() error {
	if err := q.AgeRange.Validate(); err != nil {
		return err
	}

	if q.Age < AgeMin || q.Age > AgeMax {
		return fmt.Errorf("возраст должен быть не менее %d лет и не более %d лет", AgeMin, AgeMax)
	}

	return nil
}

func (q *Questionnaire) Activate() {
	q.IsActive = true
}

func (q *Questionnaire) Deactivate() {
	q.IsActive = false
}

func (q *Questionnaire) PreferredGenders() []Gender {
	switch q.Orientation {
	case OrientationBi:
		return []Gender{GenderMale, GenderFemale}
	case OrientationHetero:
		return []Gender{!q.Gender}
	case OrientationHomo:
		return []Gender{q.Gender}
	default:
		return []Gender{}
	}
}

func (q *Questionnaire) CheckCompatibility(questionnaire *Questionnaire) bool {
	if !q.AgeRange.InRange(questionnaire.Age) {
		return false
	}

	for _, v := range q.PreferredGenders() {
		if v == questionnaire.Gender {
			return true
		}
	}

	return false
}
