package main

import (
	"database/sql"
	"meet/internal/pkg/app/repository"
	"meet/internal/pkg/app/service"
	"meet/internal/pkg/tgbot/command"
	"meet/internal/pkg/tgbot/router"
	routercommand "meet/internal/pkg/tgbot/router/command"
	routerstate "meet/internal/pkg/tgbot/router/state"
)

func configuredRouter(
	db *sql.DB,
	userRepository repository.UserRepository,
	questionnaireRepository repository.QuestionnaireRepository,
) router.Matcher {
	userService := service.NewUserService(userRepository)
	questionnaireService := service.NewQuestionnaireService(questionnaireRepository)

	startCommand := command.NewStartCommand(db, userService, questionnaireService)
	questionnaireCommand := command.NewQuestionnaireCommand(questionnaireRepository, questionnaireService)

	commandR := router.New()
	routercommand.NewConfigurator(startCommand).Configure(commandR)

	stateR := router.New()
	routerstate.NewConfigurator(questionnaireCommand).Configure(stateR)

	return router.NewCompositeRouter(commandR, stateR)
}
