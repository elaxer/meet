package model

import (
	"strings"
	"time"

	"github.com/guregu/null"
	"github.com/looplab/fsm"
)

type Gender struct{ null.Bool }

var (
	GenderMale   = Gender{Bool: null.BoolFrom(false)}
	GenderFemale = Gender{Bool: null.BoolFrom(true)}
)

type Orientation struct{ null.String }

var (
	OrientationHetero = Orientation{null.StringFrom("hetero")}
	OrientationHomo   = Orientation{null.StringFrom("homo")}
	OrientationBi     = Orientation{null.StringFrom("bi")}
)

type MeetingPurpose struct{ null.String }

var (
	MeetingPurposeFriendship   = MeetingPurpose{null.StringFrom("friendship")}
	MeetingPurposeRelationship = MeetingPurpose{null.StringFrom("relationship")}
	MeetingPurposeSex          = MeetingPurpose{null.StringFrom("sex")}
)

const (
	ageMin = 18
	ageMax = 65

	questionnaireAboutLengthMax = 2048
	questionnaireNameLengthMin  = 2
	questionnaireNameLengthMax  = 32
)

const (
	EventQuestionnaireFillName        = "fill_name"
	EventQuestionnaireFillBirthDate   = "fill_birth_date"
	EventQuestionnaireFillGender      = "fill_gender"
	EventQuestionnaireFillOrientation = "fill_orientation"
	EventQuestionnaireFillAgeRangeMin = "fill_age_range_min"
	EventQuestionnaireFillAgeRangeMax = "fill_age_range_max"
	EventQuestionnaireFillCity        = "fill_city"
	EventQuestionnaireFillAbout       = "fill_about"

	EventQuestionnaireComplete = "complete"
)

const (
	StateQuestionnaireCreated   = "created"
	StateQuestionnaireCompleted = "completed"

	StateQuestionnaireFillingName        = "filling_name"
	StateQuestionnaireFillingBirthDate   = "filling_birth_date"
	StateQuestionnaireFillingGender      = "filling_gender"
	StateQuestionnaireFillingOrientation = "filling_orientation"
	StateQuestionnaireFillingAgeRangeMin = "filling_age_range_min"
	StateQuestionnaireFillingAgeRangeMax = "filling_age_range_max"
	StateQuestionnaireFillingCity        = "filling_city"
	StateQuestionnaireFillingAbout       = "filling_about"
)

var questionnaireEventDes = []fsm.EventDesc{
	{Name: EventQuestionnaireFillName, Src: []string{StateQuestionnaireCreated}, Dst: StateQuestionnaireFillingName},
	{Name: EventQuestionnaireFillBirthDate, Src: []string{StateQuestionnaireFillingName, StateQuestionnaireCompleted}, Dst: StateQuestionnaireFillingBirthDate},
	{Name: EventQuestionnaireFillGender, Src: []string{StateQuestionnaireFillingBirthDate, StateQuestionnaireCompleted}, Dst: StateQuestionnaireFillingGender},
	{Name: EventQuestionnaireFillOrientation, Src: []string{StateQuestionnaireFillingGender, StateQuestionnaireCompleted}, Dst: StateQuestionnaireFillingOrientation},
	{Name: EventQuestionnaireFillAgeRangeMin, Src: []string{StateQuestionnaireFillingOrientation, StateQuestionnaireCompleted}, Dst: StateQuestionnaireFillingAgeRangeMin},
	{Name: EventQuestionnaireFillAgeRangeMax, Src: []string{StateQuestionnaireFillingAgeRangeMin, StateQuestionnaireCompleted}, Dst: StateQuestionnaireFillingAgeRangeMax},
	{Name: EventQuestionnaireFillCity, Src: []string{StateQuestionnaireFillingAgeRangeMax, StateQuestionnaireCompleted}, Dst: StateQuestionnaireFillingCity},
	{Name: EventQuestionnaireFillAbout, Src: []string{StateQuestionnaireFillingCity, StateQuestionnaireCompleted}, Dst: StateQuestionnaireFillingAbout},
	{Name: EventQuestionnaireComplete, Src: []string{StateQuestionnaireFillingAbout}, Dst: StateQuestionnaireCompleted},
}

var (
	errQuestionnaireNameTooShort = NewValidationError("name", "длина имени должна быть не менее %d символов", questionnaireNameLengthMin)
	errQuestionnaireNameTooLong  = NewValidationError("name", "длина имени должна быть не более %d символов", questionnaireNameLengthMax)

	errAgeMin = NewValidationError("age", "возраст должен быть не менее %d лет", ageMin)
	errAgeMax = NewValidationError("age", "возраст должен быть не менее %d лет", ageMax)

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
	CityID         null.Int       `json:"city_id"`
	About          string         `json:"about"`
	Photos         []*Photo       `json:"photos"`
	IsActive       bool           `json:"is_active"`
	state          string
	*fsm.FSM
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
		&q.AgeRange.Min,
		&q.AgeRange.Max,
		&q.CityID,
		&q.About,
		&q.IsActive,
		&q.state,
	)
}

func NewQuestionnaireEmpty() *Questionnaire {
	q := new(Questionnaire)
	q.FSM = fsm.NewFSM(StateQuestionnaireCreated, questionnaireEventDes, fsm.Callbacks{})

	return q
}

func NewQuestionnaire(userID int, name string) *Questionnaire {
	q := NewQuestionnaireEmpty()
	q.UserID = userID
	q.Name = name

	return q
}

func (q *Questionnaire) State() string {
	return q.state
}

func (q *Questionnaire) BeforeAdd() {
	q.BaseModel.BeforeAdd()

	q.About = strings.TrimSpace(q.About)
}

func (q *Questionnaire) BeforeUpdate() {
	q.BaseModel.BeforeUpdate()

	q.About = strings.TrimSpace(q.About)
}

func (q *Questionnaire) Validate() error {
	currentTime := time.Now()

	errs := &ValidationErrors{}

	if err := q.AgeRange.Validate(); err != nil {
		errs.Append(err)
	}

	if q.BirthDate.Valid {
		if err := q.BirthDate.Validate(currentTime); err != nil {
			errs.Append(err)
		}
	}

	if len(strings.TrimSpace(q.About)) > questionnaireAboutLengthMax {
		errs.Append(errQuestionnaireAboutTooLong)
	}

	name := strings.TrimSpace(q.Name)
	if len(name) < questionnaireNameLengthMin {
		errs.Append(errQuestionnaireNameTooShort)
	}
	if len(name) > questionnaireNameLengthMax {
		errs.Append(errQuestionnaireNameTooLong)
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
		if q.Gender == GenderMale {
			return []Gender{GenderFemale}
		}
		return []Gender{GenderMale}
	case OrientationHomo:
		return []Gender{q.Gender}
	default:
		return []Gender{}
	}
}

func (q *Questionnaire) checkCompatibility(questionnaire *Questionnaire, currentTime time.Time) bool {
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

func (q *Questionnaire) IsCompleted() bool {
	checks := []bool{
		q.BirthDate.Valid,
		q.Gender.Valid,
		q.Orientation.Valid,
		q.MeetingPurpose.Valid,
		q.AgeRange.Min.Valid,
		q.AgeRange.Max.Valid,
		q.Gender.Valid,
		q.CityID.Valid,
	}

	for _, check := range checks {
		if !check {
			return false
		}
	}

	return true
}

func (q *Questionnaire) IsReady() bool {
	return q.FSM.Is(StateQuestionnaireCreated) || q.FSM.Is(StateQuestionnaireCompleted)
}
