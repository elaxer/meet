package model

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	BaseModel
	Login        string `json:"login"`
	PasswordHash string `json:"-"`
	IsBlocked    bool   `json:"-"`
}

func NewUser() *User {
	return new(User)
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

func (u *User) Validate() error {
	if u.Login == "" {
		// TODO
		return errors.New("логин пользователя не может быть пустым")
	}
	if u.PasswordHash == "" {
		return errors.New("хэш пароля пользователя не может быть пустым")
	}

	return nil
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
