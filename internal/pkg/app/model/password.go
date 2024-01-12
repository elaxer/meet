package model

import (
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var (
	errPasswordEmpty       = NewValidationError("password", "пароль не может быть пустым")
	errPasswordWrongLength = NewValidationError(
		"password",
		"длина пароля должна быть не менее %d и не более %d символов",
		passwordMinLength,
		passwordMaxLength,
	)
)

const (
	passwordMinLength = 5
	passwordMaxLength = 64
)

type Password string

func (p Password) Validate() error {
	errs := &ValidationErrors{}

	if strings.TrimSpace(string(p)) == "" {
		errs.Append(errPasswordEmpty)
	}
	if len(p) < passwordMinLength || len(p) > passwordMaxLength {
		errs.Append(errPasswordWrongLength)
	}

	if errs.Empty() {
		return nil
	}

	return errs
}

func (p Password) Hash() (string, error) {
	ph, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(ph), nil
}
