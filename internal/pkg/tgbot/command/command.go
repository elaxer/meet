package command

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	skipAction = "Пропустить"
)

type Command interface {
	Handle(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update)
}

type CommandFunc func(update tgbotapi.Update) error
