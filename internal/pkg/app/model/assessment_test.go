package model

import (
	"reflect"
	"testing"
)

func TestUserDirection_Validate(t *testing.T) {
	type fields struct {
		FromUserID int
		ToUserID   int
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			"IDs are equals",
			fields{FromUserID: 1252364, ToUserID: 1252364},
			true,
		},
		{
			"IDs are not equals",
			fields{FromUserID: 1252364, ToUserID: 1252363},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ad := &Direction{
				FromID: tt.fields.FromUserID,
				ToID:   tt.fields.ToUserID,
			}
			if err := ad.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("UserDirection.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserDirection_NewReversed(t *testing.T) {
	type fields struct {
		FromUserID int
		ToUserID   int
	}
	tests := []struct {
		name   string
		fields fields
		want   Direction
	}{
		{
			"New struct are correctly reversed",
			fields{123, 321},
			Direction{321, 123},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ad := &Direction{
				FromID: tt.fields.FromUserID,
				ToID:   tt.fields.ToUserID,
			}
			if got := ad.NewReversed(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserDirection.NewReversed() = %v, want %v", got, tt.want)
			}
		})
	}
}
