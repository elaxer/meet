package model

import (
	"regexp"
	"strings"

	"github.com/guregu/null"
	"golang.org/x/crypto/bcrypt"
)

const (
	loginLengthMin = 3
	loginLengthMax = 25
)

var (
	errUserLoginEmpty             = NewValidationError("login", "логин пользователя не может быть пустым")
	errUserLoginTooShort          = NewValidationError("login", "длина логина пользователя должна быть не менее %d символов", loginLengthMin)
	errUserLoginTooLong           = NewValidationError("login", "длина логина пользователя должна быть не более %d символов", loginLengthMax)
	errUserLoginInvalidCharacters = NewValidationError("login", "логин пользователя должен содержать в себе только буквы и цифры")

	errUserTgIDInvalid = NewValidationError("tg_id", "неверное значение идентификатора пользователя Telegram")
)

type User struct {
	BaseModel
	Login        string      `json:"login"`
	PasswordHash null.String `json:"-"`
	TgID         null.Int    `json:"tg_id"`
	IsBlocked    bool        `json:"-"`
}

// GetFieldPointers реализует интерфейс Model
func (u *User) GetFieldPointers() []interface{} {
	return append(
		u.BaseModel.GetFieldPointers(),
		&u.Login,
		&u.PasswordHash,
		&u.IsBlocked,
		&u.TgID,
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
		errs.Append(errUserLoginEmpty)
	}
	if len(login) < loginLengthMin {
		errs.Append(errUserLoginTooShort)
	}
	if len(login) > loginLengthMax {
		errs.Append(errUserLoginTooLong)
	}
	if !regexp.MustCompile("^[a-zA-Z0-9]+$").MatchString(login) {
		errs.Append(errUserLoginInvalidCharacters)
	}

	if u.TgID.Valid && u.TgID.IsZero() {
		errs.Append(errUserTgIDInvalid)
	}

	if errs.Empty() {
		return nil
	}

	return errs
}

func (u *User) ComparePassword(password Password) bool {
	return u.PasswordHash.Valid && bcrypt.CompareHashAndPassword([]byte(u.PasswordHash.String), []byte(password)) == nil
}

func (u *User) Ban() {
	u.IsBlocked = true
}

func (u *User) Unban() {
	u.IsBlocked = false
}
