package model

import (
	"reflect"
	"testing"
	"time"

	"github.com/guregu/null"
)

func TestQuestionnaire_GetPreferredGenders(t *testing.T) {
	type fields struct {
		Gender      Gender
		Orientation Orientation
	}
	tests := []struct {
		name   string
		fields fields
		want   []Gender
	}{
		{
			"Hetero Male",
			fields{
				Gender:      GenderMale,
				Orientation: OrientationHetero,
			},
			[]Gender{GenderFemale},
		},
		{
			"Homo Male",
			fields{
				Gender:      GenderMale,
				Orientation: OrientationHomo,
			},
			[]Gender{GenderMale},
		},
		{
			"Bi Male",
			fields{
				Gender:      GenderMale,
				Orientation: OrientationBi,
			},
			[]Gender{GenderMale, GenderFemale},
		},
		{
			"Hetero Female",
			fields{
				Gender:      GenderFemale,
				Orientation: OrientationHetero,
			},
			[]Gender{GenderMale},
		},
		{
			"Homo Female",
			fields{
				Gender:      GenderFemale,
				Orientation: OrientationHomo,
			},
			[]Gender{GenderFemale},
		},
		{
			"Bi Female",
			fields{
				Gender:      GenderMale,
				Orientation: OrientationBi,
			},
			[]Gender{GenderMale, GenderFemale},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Questionnaire{
				Gender:      tt.fields.Gender,
				Orientation: tt.fields.Orientation,
			}
			if got := q.PreferredGenders(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Questionnaire.PreferredGenders() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuestionnaire_CheckCompatibilities(t *testing.T) {
	type fields struct {
		BirthDate   BirthDate
		Gender      Gender
		Orientation Orientation
		AgeRange    AgeRange
		IsActive    bool
	}
	type args struct {
		questionnaire *Questionnaire
	}

	currentDate := time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			"The questionnaires are compatible",
			// The "-20" means 20 years old
			fields{BirthDateFrom(currentDate.AddDate(-20, 0, 0)), GenderMale, OrientationHetero, NewAgeRange(18, 30), true},
			args{
				&Questionnaire{
					BirthDate:   BirthDateFrom(currentDate.AddDate(-23, 0, 0)),
					Gender:      GenderFemale,
					Orientation: OrientationHetero,
					AgeRange:    NewAgeRange(20, 40),
					IsActive:    true,
				},
			},
			true,
		},
		{
			"The questionnaires are not compatible because of age",
			fields{BirthDateFrom(currentDate.AddDate(-20, 0, 0)), GenderMale, OrientationHetero, NewAgeRange(18, 30), true},
			args{
				&Questionnaire{
					BirthDate:   BirthDateFrom(currentDate.AddDate(-31, 0, 0)),
					Gender:      GenderFemale,
					Orientation: OrientationHetero,
					AgeRange:    NewAgeRange(20, 40),
					IsActive:    true,
				},
			},
			false,
		},
		{
			"The questionnaires are not compatible because of orientation",
			fields{BirthDateFrom(currentDate.AddDate(-20, 0, 0)), GenderMale, OrientationHetero, NewAgeRange(18, 30), true},
			args{
				&Questionnaire{
					BirthDate:   BirthDateFrom(currentDate.AddDate(-23, 0, 0)),
					Gender:      GenderFemale,
					Orientation: OrientationHomo,
					AgeRange:    NewAgeRange(20, 40),
					IsActive:    true,
				},
			},
			false,
		},
		{
			"The questionnaires are not compatible because of orientation (2)",
			fields{BirthDateFrom(currentDate.AddDate(-20, 0, 0)), GenderMale, OrientationHomo, NewAgeRange(18, 30), true},
			args{
				&Questionnaire{
					BirthDate:   BirthDateFrom(currentDate.AddDate(-23, 0, 0)),
					Gender:      GenderFemale,
					Orientation: OrientationHomo,
					AgeRange:    NewAgeRange(20, 40),
					IsActive:    true,
				},
			},
			false,
		},
		{
			"The questionnaires are not compatible because of orientation (3)",
			fields{BirthDateFrom(currentDate.AddDate(-20, 0, 0)), GenderFemale, OrientationHetero, NewAgeRange(18, 30), true},
			args{
				&Questionnaire{
					BirthDate:   BirthDateFrom(currentDate.AddDate(-20, 0, 0)),
					Gender:      GenderFemale,
					Orientation: OrientationBi,
					AgeRange:    NewAgeRange(18, 30),
					IsActive:    true,
				},
			},
			false,
		},
		{
			"The questionnaires are not compatible because of age and orientation",
			fields{BirthDateFrom(currentDate.AddDate(-20, 0, 0)), GenderMale, OrientationHetero, NewAgeRange(18, 30), true},
			args{
				&Questionnaire{
					BirthDate:   BirthDateFrom(currentDate.AddDate(-30, 0, 0)),
					Gender:      GenderFemale,
					Orientation: OrientationHomo,
					AgeRange:    NewAgeRange(20, 40),
					IsActive:    true,
				},
			},
			false,
		},
		{
			"The questionnaires are not compatible because of age and orientation (reversed)",
			fields{BirthDateFrom(currentDate.AddDate(-30, 0, 0)), GenderFemale, OrientationHomo, NewAgeRange(20, 40), true},
			args{
				&Questionnaire{
					BirthDate:   BirthDateFrom(currentDate.AddDate(-20, 0, 0)),
					Gender:      GenderMale,
					Orientation: OrientationHetero,
					AgeRange:    NewAgeRange(18, 30),
					IsActive:    true,
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Questionnaire{
				Name:        "test",
				BirthDate:   tt.fields.BirthDate,
				Gender:      tt.fields.Gender,
				Orientation: tt.fields.Orientation,
				AgeRange:    tt.fields.AgeRange,
				IsActive:    tt.fields.IsActive,
			}

			tt.args.questionnaire.Name = "test"

			if got := q.CheckCompatibilities(tt.args.questionnaire, currentDate); got != tt.want {
				t.Errorf("Questionnaire.CheckCompatibility() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuestionnaire_IsCompleted(t *testing.T) {
	q := NewQuestionnaire(1, "test_user")
	if err := q.Validate(); err != nil {
		t.Error(err)
	}
	if q.IsCompleted() {
		t.Errorf("the questionnaire couldn't be completed at this moment")
	}

	q.BirthDate = NewBirthDate(2000, 1, 1)
	if err := q.Validate(); err != nil {
		t.Error(err)
	}
	if q.IsCompleted() {
		t.Errorf("the questionnaire couldn't be completed at this moment")
	}

	q.Gender = GenderMale
	if err := q.Validate(); err != nil {
		t.Error(err)
	}
	if q.IsCompleted() {
		t.Errorf("the questionnaire couldn't be completed at this moment")
	}

	q.Orientation = OrientationHetero
	if err := q.Validate(); err != nil {
		t.Error(err)
	}
	if q.IsCompleted() {
		t.Errorf("the questionnaire couldn't be completed at this moment")
	}

	q.MeetingPurpose = MeetingPurposeFriendship
	if err := q.Validate(); err != nil {
		t.Error(err)
	}
	if q.IsCompleted() {
		t.Errorf("the questionnaire couldn't be completed at this moment")
	}

	q.AgeRange = AgeRange{Min: null.IntFrom(18), Max: null.IntFrom(18)}
	if err := q.Validate(); err != nil {
		t.Error(err)
	}
	if q.IsCompleted() {
		t.Errorf("the questionnaire couldn't be completed at this moment")
	}

	q.CityID = null.IntFrom(1)
	if err := q.Validate(); err != nil {
		t.Error(err)
	}

	if !q.IsCompleted() {
		t.Errorf("the questionnaire had to be completed at this moment")
	}
}
