package tgbot

import (
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func ResponseError(bot *tgbotapi.BotAPI, chatID int64, err error) {
	slog.Warn(err.Error())
	bot.Send(tgbotapi.NewMessage(chatID, "Произошла ошибка, попробуйте снова"))
}
