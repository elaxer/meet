package main

import (
	"context"
	"database/sql"
	"log/slog"
	"meet/internal/config"
	"meet/internal/pkg/app"
	"meet/internal/pkg/app/database"
	"meet/internal/pkg/app/model"
	"meet/internal/pkg/app/rdatabase"
	"meet/internal/pkg/app/repository"
	"meet/internal/pkg/app/repository/dbrepository"
	"meet/internal/pkg/app/repository/rdbrepository"
	"meet/internal/pkg/app/service"
	"meet/internal/pkg/app/slogger"
	"meet/internal/pkg/tgbot/command"
	"meet/internal/pkg/tgbot/router"
	routercommand "meet/internal/pkg/tgbot/router/command"
	routerstate "meet/internal/pkg/tgbot/router/state"
	"path/filepath"
	"runtime"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var (
	questionnaireRepository repository.QuestionnaireRepository
	userService             service.UserService
)

func main() {
	_, b, _, _ := runtime.Caller(0)
	rootDir, _ := filepath.Abs(filepath.Dir(b) + "/../..")

	err := godotenv.Load(rootDir + "/.env")
	if err != nil {
		panic(err)
	}

	cfg := config.FromEnv(rootDir)

	logF := slogger.MustOpenLog(rootDir)
	defer logF.Close()

	slogger.Configure(cfg.Debug, logF)

	db := database.MustConnect(cfg.DB)
	rdb := rdatabase.MustConnect(cfg.Redis)

	questionnaireRepository = dbrepository.NewQuestionnaireRepository(db)
	userService = service.NewUserService(
		dbrepository.NewUserRepository(db),
		rdbrepository.NewUserRedisRepository(rdb),
	)

	bot, err := tgbotapi.NewBotAPI(cfg.TgBot.Token)
	if err != nil {
		panic(err)
	}

	bot.Debug = cfg.Debug

	u := tgbotapi.NewUpdate(0)
	u.Timeout = cfg.TgBot.UpdateTimeout

	serve(bot, u, configuredRouter(db))
}

func serve(bot *tgbotapi.BotAPI, u tgbotapi.UpdateConfig, r router.Matcher) {
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		u, q, err := userWithQuestionnaire(update.Message.From.ID)
		if err != nil {
			slog.Warn(err.Error(), "user", u, "questionnaire", q)
			continue
		}

		ident := identifier(update.Message, q)
		if ident == "" {
			continue
		}

		ctx := context.WithValue(context.Background(), app.CtxKeyUser, u)
		ctx = context.WithValue(ctx, app.CtxKeyQuestionnaire, q)

		if handler := r.Match(ident); handler != nil {
			handler(ctx, bot, update)
		}
	}
}

func userWithQuestionnaire(tgID int64) (*model.User, *model.Questionnaire, error) {
	u, err := userService.GetByTgID(tgID)
	if err == repository.ErrNotFound {
		return nil, nil, nil
	} else if err != nil {
		return nil, nil, err
	}

	q, err := questionnaireRepository.GetByUserID(u.ID)
	if err != nil {
		return nil, nil, err
	}

	return u, q, nil
}

func identifier(message *tgbotapi.Message, questionnaire *model.Questionnaire) string {
	if message.IsCommand() {
		return message.Command()
	} else if questionnaire != nil {
		return questionnaire.FSM.Current()
	}

	return ""
}

func configuredRouter(db *sql.DB) router.Matcher {
	var (
		questionnaireService = service.NewQuestionnaireService(questionnaireRepository)
	)

	var (
		startCommand         = command.NewStartCommand(db, userService, questionnaireService)
		questionnaireCommand = command.NewQuestionnaireCommand(questionnaireRepository, questionnaireService)
	)

	commandR := router.New()
	routercommand.NewConfigurator(startCommand).Configure(commandR)

	stateR := router.New()
	routerstate.NewConfigurator(questionnaireCommand).Configure(stateR)

	return router.NewCompositeRouter(commandR, stateR)
}
