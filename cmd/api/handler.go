package main

import (
	"database/sql"
	"meet/internal/pkg/api/handler"
	"meet/internal/pkg/api/middleware"
	"meet/internal/pkg/api/router"
	"meet/internal/pkg/app/helper"
	"meet/internal/pkg/app/repository"
	"meet/internal/pkg/app/service"
	"net/http"

	"github.com/gorilla/mux"
)

func httpHandler(db *sql.DB) http.Handler {
	urlHelper := helper.NewURLHelper(cfg.ServerConfig, cfg.PathConfig.UploadDirs)
	pathHelper := helper.NewPathHelper(cfg.PathConfig)

	assessmentRepository := repository.NewAssessmentDBRepository(db)
	messageRepository := repository.NewMessageDBRepository(db)
	photoRepository := repository.NewPhotoDBRepository(db)
	questionnaireRepository := repository.NewQuestionnaireDBRepository(db, photoRepository)
	userRepository := repository.NewUserDBRepository(db)
	countryRepository := repository.NewCountryDBRepository(db)
	cityRepository := repository.NewCityDBRepository(db)
	coupleRepository := repository.NewCoupleDBRepository(db)

	assessmentService := service.NewAssessmentService(db, assessmentRepository, coupleRepository, questionnaireRepository)
	authService := service.NewAuthService(cfg.JWTConfig, userRepository)
	fileUploaderService := service.NewFileUploaderService(pathHelper, cfg.PathConfig)
	messageService := service.NewMessageService(messageRepository, coupleRepository)
	photoService := service.NewPhotoService(pathHelper, photoRepository, questionnaireRepository, fileUploaderService)
	questionnaireService := service.NewQuestionnaireService(questionnaireRepository)
	userService := service.NewUserService(userRepository)

	authorizeMiddleware := middleware.NewAuthorizeMiddleware(authService)

	assessmentHandler := handler.NewAssessmentHandler(assessmentService)
	authHandler := handler.NewAuthHandler(authService)
	messageHandler := handler.NewMessageHandler(messageRepository, messageService)
	photoHandler := handler.NewPhotoHandler(urlHelper, photoRepository, photoService)
	questionnaireHandler := handler.NewQuestionnaireHandler(urlHelper, questionnaireRepository, questionnaireService)
	swaggerHandler := handler.NewSwaggerHandler(cfg.PathConfig)
	userHandler := handler.NewUserHandler(userRepository, userService)
	dictionaryHandler := handler.NewDictionaryHandler(countryRepository, cityRepository)

	rc := router.NewConfigurator(
		cfg.PathConfig,

		authorizeMiddleware,

		assessmentHandler,
		authHandler,
		messageHandler,
		photoHandler,
		questionnaireHandler,
		swaggerHandler,
		userHandler,
		dictionaryHandler,
	)

	return rc.Configure(mux.NewRouter())
}
