package model

import (
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

const (
	loginLengthMin = 3
	loginLengthMax = 25
)

var (
	ErrUserLoginEmpty             = NewValidationError("login", "логин пользователя не может быть пустым")
	ErrUserLoginTooShort          = NewValidationError("login", "длина логина пользователя должна быть не менее %d символов", loginLengthMin)
	ErrUserLoginTooLong           = NewValidationError("login", "длина логина пользователя должна быть не более %d символов", loginLengthMax)
	ErrUserLoginInvalidCharacters = NewValidationError("login", "логин пользователя должен содержать в себе только буквы и цифры")

	ErrUserPasswordHashEmpty = NewValidationError("password_hash", "хэш пароля пользователя не может быть пустым")
)

type User struct {
	BaseModel
	Login        string `json:"login"`
	PasswordHash string `json:"-"`
	IsBlocked    bool   `json:"-"`
}

// GetFieldPointers реализует интерфейс Model
func (u *User) GetFieldPointers() []interface{} {
	return append(
		u.BaseModel.GetFieldPointers(),
		&u.Login,
		&u.PasswordHash,
		&u.IsBlocked,
	)
}

func (u *User) BeforeAdd() {
	u.BaseModel.BeforeAdd()

	u.Login = strings.TrimSpace(u.Login)
}

func (u *User) BeforeUpdate() {
	u.BaseModel.BeforeUpdate()

	u.Login = strings.TrimSpace(u.Login)
}

func (u *User) Validate() error {
	errs := &ValidationErrors{}

	login := strings.TrimSpace(u.Login)
	if login == "" {
		errs.Append(ErrUserLoginEmpty)
	}
	if len(login) < loginLengthMin {
		errs.Append(ErrUserLoginTooShort)
	}
	if len(login) > loginLengthMax {
		errs.Append(ErrUserLoginTooLong)
	}
	if !regexp.MustCompile("^[a-zA-Z0-9]+$").MatchString(login) {
		errs.Append(ErrUserPasswordHashEmpty)
	}

	if strings.TrimSpace(u.PasswordHash) == "" {
		errs.Append(ErrUserPasswordHashEmpty)
	}

	if errs.Empty() {
		return nil
	}

	return errs
}

func (u *User) ComparePassword(password Password) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)) == nil
}

func (u *User) Ban() {
	u.IsBlocked = true
}

func (u *User) Unban() {
	u.IsBlocked = false
}
