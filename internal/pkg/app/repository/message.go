package repository

import "meet/internal/pkg/app/model"

type MessageRepository interface {
	Get(id int) (*model.Message, error)
	GetList(usersDirection model.Direction, limit, offset int) ([]*model.Message, error)
	UnreadCount(userID int) (int, error)
	Add(message *model.Message) error
	Update(message *model.Message) error
}
