package model

import (
	"testing"
)

func TestValidationError_Error(t *testing.T) {
	type fields struct {
		Subject      string
		Message      string
		placeholders []any
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"Error text without any placeholders",
			fields{"login", "логин не может быть пустым", []any{}},
			"login: логин не может быть пустым",
		},
		{
			"Erros text with placeholders",
			fields{"text", "длина текста должна быть не более %d символов", []any{512}},
			"text: длина текста должна быть не более 512 символов",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ve := &ValidationError{
				Subject:      tt.fields.Subject,
				Message:      tt.fields.Message,
				placeholders: tt.fields.placeholders,
			}
			if got := ve.Error(); got != tt.want {
				t.Errorf("ValidationError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidationErrors_Error(t *testing.T) {
	type fields struct {
		Errors []error
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"One error",
			fields{
				[]error{
					NewValidationError("login", "логин не может быть пустым"),
				},
			},
			"login: логин не может быть пустым",
		},
		{
			"Several errors",
			fields{
				[]error{
					NewValidationError("login", "логин не может быть пустым"),
					NewValidationError("text", "длина текста должна быть не более %d символов", 512),
				},
			},
			"login: логин не может быть пустым\ntext: длина текста должна быть не более 512 символов",
		},
		{
			"Embedded validation errors",
			fields{
				[]error{
					NewValidationError("login", "логин не может быть пустым"),
					&ValidationErrors{
						Errors: []error{
							NewValidationError("age", "возраст должен быть заполнен"),
							NewValidationError("age", "возраст должен быть не менее %d лет", 18),
						},
					},
					NewValidationError("text", "длина текста должна быть не более %d символов", 512),
				},
			},
			"login: логин не может быть пустым\nage: возраст должен быть заполнен\nage: возраст должен быть не менее 18 лет\ntext: длина текста должна быть не более 512 символов",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ve := &ValidationErrors{
				Errors: tt.fields.Errors,
			}
			if got := ve.Error(); got != tt.want {
				t.Errorf("ValidationErrors.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}
