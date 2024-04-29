package service

import (
	"context"
	"meet/internal/pkg/app/model"
	"meet/internal/pkg/app/repository"

	"github.com/guregu/null"
)

type UserService interface {
	GetByLogin(login string) (*model.User, error)
	GetByTgID(id int64) (*model.User, error)
	Register(login string, password model.Password) (*model.User, error)
	Create(ctx context.Context, tgID int64, login string) (*model.User, error)
	ChangePassword(user *model.User, password model.Password) error
}

type userService struct {
	userDBRepository    repository.UserRepository
	userRedisRepository repository.UserRepository
}

func NewUserService(userDBRepository repository.UserRepository, userRedisRepository repository.UserRepository) UserService {
	return &userService{userDBRepository, userRedisRepository}
}

func (us *userService) GetByLogin(login string) (*model.User, error) {
	u, err := us.userRedisRepository.GetByLogin(login)
	if err == nil {
		return u, nil
	}
	if err != repository.ErrNotFound {
		return nil, err
	}

	u, err = us.userDBRepository.GetByLogin(login)
	if err != nil {
		return nil, err
	}

	if err := us.userRedisRepository.Add(context.Background(), u); err != nil {
		return nil, err
	}

	return u, nil
}

func (us *userService) GetByTgID(id int64) (*model.User, error) {
	u, err := us.userRedisRepository.GetByTgID(id)
	if err == nil {
		return u, nil
	}
	if err != repository.ErrNotFound {
		return nil, err
	}

	u, err = us.userDBRepository.GetByTgID(id)
	if err != nil {
		return nil, err
	}

	if err := us.userRedisRepository.Add(context.Background(), u); err != nil {
		return nil, err
	}

	return u, nil
}

func (us *userService) Register(login string, password model.Password) (*model.User, error) {
	hasUser, err := us.userDBRepository.HasByLogin(login)
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

	if err := us.userDBRepository.Add(context.Background(), u); err != nil {
		return nil, err
	}

	return u, nil
}

func (us *userService) Create(ctx context.Context, tgID int64, login string) (*model.User, error) {
	hasUser, err := us.userDBRepository.HasByTgID(tgID)
	if err != nil {
		return nil, err
	}
	if hasUser {
		return nil, repository.ErrDuplicate
	}

	u := new(model.User)
	u.Login = login
	u.TgID = null.IntFrom(tgID)

	err = us.userDBRepository.Add(ctx, u)

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

	return us.userDBRepository.Update(user)
}
