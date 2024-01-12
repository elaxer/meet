package service

import (
	"context"
	"meet/internal/pkg/app/model"
	"meet/internal/pkg/app/repository"

	"github.com/guregu/null"
)

type UserService interface {
	Register(login string, password model.Password) (*model.User, error)
	Create(ctx context.Context, tgID int64, login string) (*model.User, error)
	ChangePassword(user *model.User, password model.Password) error
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{userRepository}
}

func (us *userService) Register(login string, password model.Password) (*model.User, error) {
	hasUser, err := us.userRepository.HasByLogin(login)
	if err != nil {
		return nil, err
	}
	if hasUser {
		return nil, repository.ErrDuplicate
	}

	if err := password.Validate(); err != nil {
		return nil, err
	}

	hash, err := password.Hash()
	if err != nil {
		return nil, err
	}

	u := new(model.User)
	u.Login = login
	u.PasswordHash = null.StringFrom(hash)

	if err := us.userRepository.Add(context.Background(), u); err != nil {
		return nil, err
	}

	return u, nil
}

func (us *userService) Create(ctx context.Context, tgID int64, login string) (*model.User, error) {
	hasUser, err := us.userRepository.HasByTgID(tgID)
	if err != nil {
		return nil, err
	}
	if hasUser {
		return nil, repository.ErrDuplicate
	}

	u := new(model.User)
	u.Login = login
	u.TgID = null.IntFrom(tgID)

	err = us.userRepository.Add(ctx, u)

	return u, err
}

func (us *userService) ChangePassword(user *model.User, password model.Password) error {
	err := password.Validate()
	if err != nil {
		return err
	}

	hash, err := password.Hash()
	if err != nil {
		return err
	}

	user.PasswordHash = null.StringFrom(hash)
	if err != nil {
		return err
	}

	err = us.userRepository.Update(user)

	return err
}
