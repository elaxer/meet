package model

import (
	"reflect"
	"testing"
	"time"
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
			fields{BirthDate(currentDate.AddDate(-20, 0, 0)), GenderMale, OrientationHetero, AgeRange{18, 30}},
			args{
				&Questionnaire{
					BirthDate:   BirthDate(currentDate.AddDate(-23, 0, 0)),
					Gender:      GenderFemale,
					Orientation: OrientationHetero,
					AgeRange:    AgeRange{20, 40},
				},
			},
			true,
		},
		{
			"The questionnaires are not compatible because of age",
			fields{BirthDate(currentDate.AddDate(-20, 0, 0)), GenderMale, OrientationHetero, AgeRange{18, 30}},
			args{
				&Questionnaire{
					BirthDate:   BirthDate(currentDate.AddDate(-31, 0, 0)),
					Gender:      GenderFemale,
					Orientation: OrientationHetero,
					AgeRange:    AgeRange{20, 40},
				},
			},
			false,
		},
		{
			"The questionnaires are not compatible because of orientation",
			fields{BirthDate(currentDate.AddDate(-20, 0, 0)), GenderMale, OrientationHetero, AgeRange{18, 30}},
			args{
				&Questionnaire{
					BirthDate:   BirthDate(currentDate.AddDate(-23, 0, 0)),
					Gender:      GenderFemale,
					Orientation: OrientationHomo,
					AgeRange:    AgeRange{20, 40},
				},
			},
			false,
		},
		{
			"The questionnaires are not compatible because of orientation (2)",
			fields{BirthDate(currentDate.AddDate(-20, 0, 0)), GenderMale, OrientationHomo, AgeRange{18, 30}},
			args{
				&Questionnaire{
					BirthDate:   BirthDate(currentDate.AddDate(-23, 0, 0)),
					Gender:      GenderFemale,
					Orientation: OrientationHomo,
					AgeRange:    AgeRange{20, 40},
				},
			},
			false,
		},
		{
			"The questionnaires are not compatible because of orientation (3)",
			fields{BirthDate(currentDate.AddDate(-20, 0, 0)), GenderFemale, OrientationHetero, AgeRange{18, 30}},
			args{
				&Questionnaire{
					BirthDate:   BirthDate(currentDate.AddDate(-20, 0, 0)),
					Gender:      GenderFemale,
					Orientation: OrientationBi,
					AgeRange:    AgeRange{18, 30},
				},
			},
			false,
		},
		{
			"The questionnaires are not compatible because of age and orientation",
			fields{BirthDate(currentDate.AddDate(-20, 0, 0)), GenderMale, OrientationHetero, AgeRange{18, 30}},
			args{
				&Questionnaire{
					BirthDate:   BirthDate(currentDate.AddDate(-30, 0, 0)),
					Gender:      GenderFemale,
					Orientation: OrientationHomo,
					AgeRange:    AgeRange{20, 40},
				},
			},
			false,
		},
		{
			"The questionnaires are not compatible because of age and orientation (reversed)",
			fields{BirthDate(currentDate.AddDate(-30, 0, 0)), GenderFemale, OrientationHomo, AgeRange{20, 40}},
			args{
				&Questionnaire{
					BirthDate:   BirthDate(currentDate.AddDate(-20, 0, 0)),
					Gender:      GenderMale,
					Orientation: OrientationHetero,
					AgeRange:    AgeRange{18, 30},
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Questionnaire{
				BirthDate:   tt.fields.BirthDate,
				Gender:      tt.fields.Gender,
				Orientation: tt.fields.Orientation,
				AgeRange:    tt.fields.AgeRange,
			}

			if got := q.CheckCompatibilities(tt.args.questionnaire, currentDate); got != tt.want {
				t.Errorf("Questionnaire.CheckCompatibility() = %v, want %v", got, tt.want)
			}
		})
	}
}
