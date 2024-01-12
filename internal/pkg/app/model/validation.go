package model

import (
	"encoding/json"
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

func (ve *ValidationError) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Subject string `json:"subject"`
		Message string `json:"message"`
	}{ve.Subject, ve.String()})
}

func NewValidationError(subject, message string, placeholders ...any) *ValidationError {
	return &ValidationError{
		Subject:      subject,
		Message:      message,
		placeholders: placeholders,
	}
}

func (ve *ValidationError) Error() string {
	return ve.Subject + ": " + ve.String()
}

func (ve *ValidationError) String() string {
	return fmt.Sprintf(ve.Message, ve.placeholders...)
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

func (ve *ValidationErrors) First() error {
	if ve.Empty() {
		return nil
	}

	return ve.Errors[0]
}
