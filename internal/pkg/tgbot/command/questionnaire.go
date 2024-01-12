package command

import (
	"context"
	"meet/internal/pkg/app/model"
	"meet/internal/pkg/app/repository"
	"meet/internal/pkg/app/service"
	"meet/internal/pkg/tgbot"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	genderMale   = "Мужской"
	genderFemale = "Женский"
)

var (
	orientationHetero = "Гетеро"
	orientationGay    = "Гей"
	orientationLesbi  = "Лесби"
	orientationBi     = "Би"
)

type QuestionnaireCommand interface {
	FillName(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update)
	FillBirthDate(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update)
	FillGender(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update)
}

type questionnaireCommand struct {
	questionnaireRepository repository.QuestionnaireRepository
	questionnaireService    service.QuestionnaireService
}

func NewQuestionnaireCommand(
	questionnaireRepository repository.QuestionnaireRepository,
	questionnaireService service.QuestionnaireService,
) QuestionnaireCommand {
	return &questionnaireCommand{questionnaireRepository, questionnaireService}
}

func (qc *questionnaireCommand) FillName(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if update.Message.Text == skipAction {
		return
	}

	q, _ := ctx.Value(tgbot.CtxKeyQuestionnaire).(*model.Questionnaire)

	if err := qc.questionnaireService.UpdateName(q, update.Message.Text); err != nil {
		tgbot.ResponseError(bot, update.Message.From.ID, err)

		return
	}

	msg := tgbotapi.NewMessage(update.Message.From.ID, "Укажите вашу дату рождения в формате \"дд.мм.ггг\"")
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
	bot.Send(msg)
}

func (qc *questionnaireCommand) FillBirthDate(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	q, _ := ctx.Value(tgbot.CtxKeyQuestionnaire).(*model.Questionnaire)
	if !q.FSM.Is(model.StateQuestionnaireFillingBirthDate) {
		return
	}

	bd, err := time.Parse("02.01.2006", update.Message.Text)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(update.Message.From.ID, "Пожалуйста, введите дату в формате \"дд.мм.ггг\""))
	}
	q.BirthDate = model.BirthDateFrom(bd)

	if err := q.FSM.Event(context.Background(), model.EventQuestionnaireFillGender); err != nil {
		tgbot.ResponseError(bot, update.Message.From.ID, err)

		return
	}

	if err := qc.questionnaireRepository.Update(q); err != nil {
		if err, ok := err.(*model.ValidationErrors); ok {
			bot.Send(tgbotapi.NewMessage(update.Message.From.ID, err.First().Error()))

			return
		}

		tgbot.ResponseError(bot, update.Message.From.ID, err)

		return
	}

	msg := tgbotapi.NewMessage(update.Message.From.ID, "Укажите ваш пол")
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(genderMale),
			tgbotapi.NewKeyboardButton(genderFemale),
		),
	)

	bot.Send(msg)
}

func (qc *questionnaireCommand) FillGender(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	q, _ := ctx.Value(tgbot.CtxKeyQuestionnaire).(*model.Questionnaire)
	if !q.FSM.Is(model.StateQuestionnaireFillingGender) {
		return
	}

	gender := update.Message.Text
	switch strings.TrimSpace(gender) {
	case genderMale:
		q.Gender = model.GenderMale
	case genderFemale:
		q.Gender = model.GenderFemale
	default:
		return
	}

	if err := q.FSM.Event(context.Background(), model.EventQuestionnaireFillOrientation); err != nil {
		tgbot.ResponseError(bot, update.Message.From.ID, err)

		return
	}

	if err := qc.questionnaireRepository.Update(q); err != nil {
		if err, ok := err.(*model.ValidationErrors); ok {
			bot.Send(tgbotapi.NewMessage(update.Message.From.ID, err.First().Error()))

			return
		}
	}

	orientationHomo := orientationGay
	if q.Gender == model.GenderFemale {
		orientationHomo = orientationLesbi
	}

	msg := tgbotapi.NewMessage(update.Message.From.ID, "Укажите вашу ориентацию")
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(orientationHetero),
			tgbotapi.NewKeyboardButton(orientationHomo),
			tgbotapi.NewKeyboardButton(orientationBi),
		),
	)

	bot.Send(msg)
}
