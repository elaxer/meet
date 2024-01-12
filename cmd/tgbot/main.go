package main

import (
	"context"
	"log"
	"meet/internal/config"
	"meet/internal/pkg/app/helper"
	"meet/internal/pkg/app/model"
	"meet/internal/pkg/app/repository"
	"meet/internal/pkg/tgbot"
	"meet/internal/pkg/tgbot/router"

	"path/filepath"
	"runtime"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var (
	_, b, _, _ = runtime.Caller(0)
	rootDir, _ = filepath.Abs(filepath.Dir(b) + "/../..")
)

var (
	userRepository          repository.UserRepository
	questionnaireRepository repository.QuestionnaireRepository
)

var r router.Matcher

func init() {
	err := godotenv.Load(rootDir + "/.env")
	if err != nil {
		panic(err)
	}

	cfg := config.NewConfig(rootDir)

	db, err := helper.LoadDB(cfg.DBConfig)
	if err != nil {
		panic(err)
	}

	questionnaireRepository = repository.NewQuestionnaireDBRepository(db, repository.NewPhotoDBRepository(db))
	userRepository = repository.NewUserDBRepository(db)
	r = configuredRouter(db, userRepository, questionnaireRepository)
}

func main() {
	bot, err := tgbotapi.NewBotAPI("")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			return
		}

		var q *model.Questionnaire
		u, _ := userRepository.GetByTgID(update.Message.From.ID)
		if u != nil {
			q, _ = questionnaireRepository.GetByUserID(u.ID)
		}

		var identifier string
		if update.Message.IsCommand() {
			identifier = update.Message.Command()
		} else if q != nil {
			identifier = q.FSM.Current()
		} else {
			continue
		}

		ctx := context.WithValue(context.Background(), tgbot.CtxKeyUser, u)
		ctx = context.WithValue(ctx, tgbot.CtxKeyQuestionnaire, q)

		if handler := r.Match(identifier); handler != nil {
			handler(ctx, bot, update)
		}
	}
}
