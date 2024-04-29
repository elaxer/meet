package service

import (
	"context"
	"meet/internal/pkg/app/model"
	"meet/internal/pkg/app/repository"
	"testing"

	"github.com/guregu/null"
)

func Test_userService_Register(t *testing.T) {
	const login = "elaxer"

	ur := repository.NewUserRepository()
	us := NewUserService(ur, repository.NewUserRepository())

	u, err := us.Register(login, "123456")
	if err != nil {
		t.Errorf("userService.Register() = %s", err)
	}
	if u == nil {
		t.Errorf("user is nil")
	}

	has, err := ur.HasByLogin(login)
	if err != nil {
		t.Errorf("userRepository.HasByLogin() = %s", err)
	}
	if !has {
		t.Errorf("the user did not added to the repository")
	}
}

func Test_userService_ChangePassword(t *testing.T) {
	const passwordHash = "!@#$%^$@!"
	u := &model.User{Login: "elaxer", PasswordHash: null.StringFrom(passwordHash)}

	ur := repository.NewUserRepository()
	if err := ur.Add(context.Background(), u); err != nil {
		t.Errorf("userRepository.Add() = %s", err)
	}

	us := NewUserService(ur, repository.NewUserRepository())
	if err := us.ChangePassword(u, "hello world"); err != nil {
		t.Errorf("userService.ChangePassword() = %s", err)
	}

	if u.PasswordHash.String == passwordHash {
		t.Errorf("the user's password has was not changed")
	}
}
