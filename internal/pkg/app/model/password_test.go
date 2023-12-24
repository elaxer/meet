package model

import "testing"

func TestPassword_Validate(t *testing.T) {
	tests := []struct {
		name    string
		p       Password
		wantErr bool
	}{
		{
			"Password is valid",
			Password("123456"),
			false,
		},
		{
			"Password is too short",
			Password("1234"),
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.p.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Password.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
