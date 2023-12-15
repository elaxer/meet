package service

import "meet/internal/app/model"

type MessageService struct {
}

func (ms *MessageService) Text(from, to *model.User, text string) error {
	return nil
}
