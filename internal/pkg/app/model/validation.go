package model

import (
	"fmt"
	"strings"
)

type Validatable interface {
	Validate() []error
}

type ValidationError struct {
	Subject      string `json:"subject"`
	Message      string `json:"message"`
	placeholders []any  `json:"-"`
}

func NewValidationError(subject, message string, placeholders ...any) *ValidationError {
	return &ValidationError{
		Subject:      subject,
		Message:      message,
		placeholders: placeholders,
	}
}

func (ve *ValidationError) Error() string {
	return ve.Subject + ": " + fmt.Sprintf(ve.Message, ve.placeholders...)
}

type ValidationErrors struct {
	Errors []error `json:"validation_errors"`
}

func (ve *ValidationErrors) Error() string {
	var errStrings []string
	for _, err := range ve.Errors {
		errStrings = append(errStrings, err.Error())
	}

	return strings.Join(errStrings, "\n")
}

func (ve *ValidationErrors) Append(err ...error) {
	ve.Errors = append(ve.Errors, err...)
}

func (ve *ValidationErrors) Empty() bool {
	return len(ve.Errors) == 0
}
