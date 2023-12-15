package service

import (
	"meet/internal/app/model"
	"meet/internal/app/repository"
)

type UserService struct {
	userRepository repository.UserRepository
}

func newUserService(userRepository repository.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (us *UserService) Register(login string, password model.Password) (*model.User, error) {
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

	u := new(model.User)
	u.Login = login
	u.PasswordHash, err = password.GetHash()
	if err != nil {
		return nil, err
	}

	if err := us.userRepository.Add(u); err != nil {
		return nil, err
	}

	return u, nil
}

func (us *UserService) ChangePassword(user *model.User, password model.Password) error {
	err := password.Validate()
	if err != nil {
		return err
	}

	user.PasswordHash, err = password.GetHash()
	if err != nil {
		return err
	}

	err = us.userRepository.Update(user)

	return err
}
