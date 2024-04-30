package model

import (
	"testing"
)

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
			fields{ageMin, ageMax},
			false,
		},
		{
			"Ages are equal",
			fields{ageMin, ageMin},
			false,
		},
		{
			"Too small min age",
			fields{ageMin - 1, ageMax},
			true,
		},
		{
			"Too big max age",
			fields{ageMin, ageMax + 1},
			true,
		},
		{
			"All ages are small",
			fields{ageMin - 2, ageMin - 1},
			true,
		},
		{
			"All ages are big",
			fields{ageMax + 1, ageMax + 2},
			true,
		},
		{
			"Min age are greater than max age",
			fields{ageMin + 1, ageMin},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ar := NewAgeRange(tt.fields.From, tt.fields.To)
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
			ar := NewAgeRange(tt.fields.From, tt.fields.To)
			if got := ar.InRange(tt.args.age); got != tt.want {
				t.Errorf("AgeRange.InRange() = %v, want %v", got, tt.want)
			}
		})
	}
}
