package repository

import (
	"context"
	"errors"
	"meet/internal/pkg/app/model"
)

var (
	ErrNotFound  = errors.New("не удалось найти модель в репозитории")
	ErrDuplicate = errors.New("в репозитории уже сущетсвует модель с таким идентификатором")
)

type collectionRepository[T model.Model] struct {
	models []T
}

func (cr *collectionRepository[T]) Add(ctx context.Context, model T) error {
	model.BeforeAdd()

	if err := model.Validate(); err != nil {
		return err
	}

	cr.models = append(cr.models, model)

	return nil
}

func (cr *collectionRepository[T]) Update(model T) error {
	if !cr.has(model) {
		return ErrNotFound
	}

	model.BeforeUpdate()

	err := model.Validate()

	return err
}

func (cr *collectionRepository[T]) Remove(model T) error {
	if !cr.has(model) {
		return ErrNotFound
	}

	var index int
	for i, m := range cr.models {
		if m == model {
			index = i
		}
	}

	if err := model.Validate(); err != nil {
		return err
	}

	cr.models = append(cr.models[:index], cr.models[index+1:]...)

	return nil
}

func (cr *collectionRepository[T]) has(model T) bool {
	for _, m := range cr.models {
		if m == model {
			return true
		}
	}

	return false
}
