package command

import (
	"context"
	"database/sql"
	"meet/internal/pkg/app/database"
	"meet/internal/pkg/app/model"
	"meet/internal/pkg/app/service"
	"meet/internal/pkg/tgbot"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type startCommand struct {
	db                   *sql.DB
	userService          service.UserService
	questionnaireService service.QuestionnaireService
}

func NewStartCommand(db *sql.DB, userService service.UserService, questionnaireService service.QuestionnaireService) Command {
	return &startCommand{db, userService, questionnaireService}
}

func (sc *startCommand) Handle(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	ctx, tx, err := database.BeginTx(ctx, sc.db)
	if err != nil {
		tgbot.ResponseError(bot, update.Message.From.ID, err)

		return
	}

	u, err := sc.userService.Create(ctx, update.Message.From.ID, update.Message.From.UserName)
	if err != nil {
		tx.Rollback()

		tgbot.ResponseError(bot, update.Message.From.ID, err)

		return
	}

	q := model.NewQuestionnaire(u.ID, update.Message.From.FirstName)
	if err := sc.questionnaireService.Add(ctx, q); err != nil {
		tx.Rollback()

		tgbot.ResponseError(bot, update.Message.From.ID, err)

		return
	}

	tx.Commit()

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Добро пожаловать в сервис знакомств! Введите ваше имя, которое будет отображаться в анкете")
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(skipAction)))

	bot.Send(msg)
}
