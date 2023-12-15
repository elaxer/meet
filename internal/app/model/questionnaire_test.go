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

func TestAgeRange_Match(t *testing.T) {
	type fields struct {
		From int
		To   int
	}
	type args struct {
		age int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			"An age is in the range",
			fields{20, 30},
			args{25},
			true,
		},
		{
			"An age is in the min threshold",
			fields{20, 30},
			args{20},
			true,
		},
		{
			"An age is in the max threshold",
			fields{20, 30},
			args{30},
			true,
		},
		{
			"An age is below the min threshold",
			fields{20, 30},
			args{19},
			false,
		},
		{
			"An age is above the max threshold",
			fields{20, 30},
			args{31},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ar := &AgeRange{
				From: tt.fields.From,
				To:   tt.fields.To,
			}
			if got := ar.InRange(tt.args.age); got != tt.want {
				t.Errorf("AgeRange.InRange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuestionnaire_CheckCompatibility(t *testing.T) {
	type fields struct {
		Age         int
		Gender      Gender
		Orientation Orientation
		AgeRange    AgeRange
	}
	type args struct {
		questionnaire *Questionnaire
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			"The questionnaires are compatible",
			fields{20, GenderMale, OrientationHetero, AgeRange{18, 30}},
			args{&Questionnaire{Age: 23, Gender: GenderFemale, Orientation: OrientationHetero, AgeRange: AgeRange{20, 40}}},
			true,
		},
		{
			"The questionnaires are not compatible because of age",
			fields{20, GenderMale, OrientationHetero, AgeRange{18, 30}},
			args{&Questionnaire{Age: 31, Gender: GenderFemale, Orientation: OrientationHetero, AgeRange: AgeRange{20, 40}}},
			false,
		},
		{
			"The questionnaires are not compatible because of orientation",
			fields{20, GenderMale, OrientationHetero, AgeRange{18, 30}},
			args{&Questionnaire{Age: 23, Gender: GenderFemale, Orientation: OrientationHomo, AgeRange: AgeRange{20, 40}}},
			false,
		},
		{
			"The questionnaires are not compatible because of orientation (2)",
			fields{20, GenderMale, OrientationHomo, AgeRange{18, 30}},
			args{&Questionnaire{Age: 23, Gender: GenderFemale, Orientation: OrientationHomo, AgeRange: AgeRange{20, 40}}},
			false,
		},
		{
			"The questionnaires are not compatible because of orientation (3)",
			fields{20, GenderFemale, OrientationHetero, AgeRange{18, 30}},
			args{&Questionnaire{Age: 20, Gender: GenderFemale, Orientation: OrientationBi, AgeRange: AgeRange{18, 30}}},
			false,
		},
		{
			"The questionnaires are not compatible because of age and orientation",
			fields{20, GenderMale, OrientationHetero, AgeRange{18, 30}},
			args{&Questionnaire{Age: 30, Gender: GenderFemale, Orientation: OrientationHomo, AgeRange: AgeRange{20, 40}}},
			false,
		},
		{
			"The questionnaires are not compatible because of age and orientation (reversed)",
			fields{30, GenderFemale, OrientationHomo, AgeRange{20, 40}},
			args{&Questionnaire{Age: 20, Gender: GenderMale, Orientation: OrientationHetero, AgeRange: AgeRange{18, 30}}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Questionnaire{
				Age:         tt.fields.Age,
				Gender:      tt.fields.Gender,
				Orientation: tt.fields.Orientation,
				AgeRange:    tt.fields.AgeRange,
			}
			q1Compatibility := q.CheckCompatibility(tt.args.questionnaire)
			q2Compatibility := tt.args.questionnaire.CheckCompatibility(q)
			got := q1Compatibility && q2Compatibility

			if got != tt.want {
				t.Errorf("Questionnaire.CheckCompatibility() = %v, want %v", got, tt.want)
			}
		})
	}
}
