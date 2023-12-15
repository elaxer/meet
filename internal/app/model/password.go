package model

import "golang.org/x/crypto/bcrypt"

var (
	ErrPasswordEmpty       = NewValidationError("password", "пароль не может быть пустым")
	ErrPasswordWrongLength = NewValidationError(
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
	if p == "" {
		return ErrPasswordEmpty
	}

	if len(p) < passwordMinLength || len(p) > passwordMaxLength {
		return ErrPasswordWrongLength
	}

	return nil
}

func (p Password) GetHash() (string, error) {
	ph, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(ph), nil
}
