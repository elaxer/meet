package model

import (
	"testing"
	"time"
)

func TestBirthDate_Age(t *testing.T) {
	type args struct {
		currentTime time.Time
	}
	tests := []struct {
		name string
		bd   BirthDate
		args args
		want int
	}{
		{
			"Birthday is happened",
			BirthDateFrom(time.Date(2000, time.April, 1, 0, 0, 0, 0, time.UTC)),
			args{time.Date(2020, time.May, 1, 0, 0, 0, 0, time.UTC)},
			20,
		},
		{
			"Birthday is not happened",
			BirthDateFrom(time.Date(2000, time.April, 1, 0, 0, 0, 0, time.UTC)),
			args{time.Date(2020, time.March, 1, 0, 0, 0, 0, time.UTC)},
			19,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.bd.Age(tt.args.currentTime); got != tt.want {
				t.Errorf("BirthDate.Age() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBirthDate_String(t *testing.T) {
	tests := []struct {
		name string
		bd   BirthDate
		want string
	}{
		{
			"Birth date string",
			BirthDateFrom(time.Date(2023, time.February, 11, 0, 0, 0, 0, time.UTC)),
			"2023-02-11",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.bd.String(); got != tt.want {
				t.Errorf("BirthDate.String() = %v, want %v", got, tt.want)
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
