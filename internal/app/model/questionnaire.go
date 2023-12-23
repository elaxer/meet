package model

import (
	"strings"
	"time"
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

const questionnaireAboutLengthMax = 2048

var (
	errQuestionnaireEmptyCountry = NewValidationError("country", "необходимо указать страну")
	errQuestionnaireEmptyCity    = NewValidationError("city", "необходимо указать город")
	errQuestionnaireAboutTooLong = NewValidationError("about", "текст описания анкеты не должен превышать %d символов", questionnaireAboutLengthMax)
)

type Questionnaire struct {
	BaseModel
	UserID         int            `json:"user_id"`
	Name           string         `json:"name"`
	BirthDate      BirthDate      `json:"birth_date"`
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

// GetFieldPointers implements Model interface
func (q *Questionnaire) GetFieldPointers() []interface{} {
	return append(
		q.BaseModel.GetFieldPointers(),
		&q.UserID,
		&q.Name,
		&q.BirthDate,
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

func (q *Questionnaire) BeforeAdd() {
	q.BaseModel.BeforeAdd()

	q.Country = strings.TrimSpace(q.Country)
	q.City = strings.TrimSpace(q.City)
	q.About = strings.TrimSpace(q.About)
}

func (q *Questionnaire) BeforeUpdate() {
	q.BaseModel.BeforeUpdate()

	q.Country = strings.TrimSpace(q.Country)
	q.City = strings.TrimSpace(q.City)
	q.About = strings.TrimSpace(q.About)
}

func (q *Questionnaire) Validate(currentTime time.Time) error {
	errs := &ValidationErrors{}

	if err := q.AgeRange.Validate(); err != nil {
		errs.Append(err)
	}
	if err := q.BirthDate.Validate(currentTime); err != nil {
		errs.Append(err)
	}
	if strings.TrimSpace(q.Country) == "" {
		errs.Append(errQuestionnaireEmptyCountry)
	}
	if strings.TrimSpace(q.City) == "" {
		errs.Append(errQuestionnaireEmptyCity)
	}
	if len(strings.TrimSpace(q.About)) > questionnaireAboutLengthMax {
		errs.Append(errQuestionnaireAboutTooLong)
	}

	if errs.Empty() {
		return nil
	}

	return errs
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

func (q *Questionnaire) checkCompatibility(questionnaire *Questionnaire, currentTime time.Time) bool {
	// todo add tests
	if !q.IsActive {
		return false
	}

	if !q.AgeRange.InRange(questionnaire.BirthDate.Age(currentTime)) {
		return false
	}

	for _, v := range q.PreferredGenders() {
		if v == questionnaire.Gender {
			return true
		}
	}

	return false
}

func (q *Questionnaire) CheckCompatibilities(questionnaire *Questionnaire, currentTime time.Time) bool {
	return q.checkCompatibility(questionnaire, currentTime) && questionnaire.checkCompatibility(q, currentTime)
}
