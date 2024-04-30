package repository

import (
	"context"
	"meet/internal/pkg/app/model"
)

type UserRepository interface {
	GetByLogin(login string) (*model.User, error)
	HasByLogin(login string) (bool, error)
	GetByTgID(id int64) (*model.User, error)
	HasByTgID(id int64) (bool, error)
	Add(ctx context.Context, user *model.User) error
	Update(user *model.User) error
	Remove(user *model.User) error
}

type userRepository struct {
	collectionRepository[*model.User]
}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

func (ur *userRepository) GetByLogin(login string) (*model.User, error) {
	for _, u := range ur.models {
		if u.Login == login {
			return u, nil
		}
	}

	return nil, ErrNotFound
}

func (ur *userRepository) GetByTgID(id int64) (*model.User, error) {
	for _, u := range ur.models {
		if u.TgID.Valid && u.TgID.Int64 == id {
			return u, nil
		}
	}

	return nil, ErrNotFound
}

func (ur *userRepository) HasByTgID(id int64) (bool, error) {
	u, _ := ur.GetByTgID(id)

	return u != nil, nil
}

func (ur *userRepository) HasByLogin(login string) (bool, error) {
	u, _ := ur.GetByLogin(login)

	return u != nil, nil
}
