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
