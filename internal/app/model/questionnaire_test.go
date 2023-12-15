package model

import (
	"reflect"
	"testing"
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

func TestAgeRange_Validate(t *testing.T) {
	type fields struct {
		From int
		To   int
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			"Ages are correct",
			fields{AgeMin, AgeMax},
			false,
		},
		{
			"Ages are equal",
			fields{AgeMin, AgeMin},
			false,
		},
		{
			"Too small min age",
			fields{AgeMin - 1, AgeMax},
			true,
		},
		{
			"Too big max age",
			fields{AgeMin, AgeMax + 1},
			true,
		},
		{
			"All ages are small",
			fields{AgeMin - 2, AgeMin - 1},
			true,
		},
		{
			"All ages are big",
			fields{AgeMax + 1, AgeMax + 2},
			true,
		},
		{
			"Min age are greater than max age",
			fields{AgeMin + 1, AgeMin},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ar := &AgeRange{
				From: tt.fields.From,
				To:   tt.fields.To,
			}
			if err := ar.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("AgeRange.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
