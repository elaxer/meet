package service

import (
	"meet/internal/app/model"
	"meet/internal/app/repository"
)

type UserService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
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
