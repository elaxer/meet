package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"meet/internal/config"
	"meet/internal/pkg/api/handler"
	"meet/internal/pkg/api/middleware"
	"meet/internal/pkg/api/router"
	"meet/internal/pkg/app/helper"
	"meet/internal/pkg/app/repository"
	"meet/internal/pkg/app/service"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

func main() {
	_, b, _, _ := runtime.Caller(0)
	rootDir, _ := filepath.Abs(filepath.Dir(b) + "/../..")

	if err := godotenv.Load(rootDir + "/.env"); err != nil {
		slog.Warn(err.Error())
		return
	}

	cfg := config.FromEnv(rootDir)

	logF, err := helper.OpenLogFile(rootDir)
	if err != nil {
		slog.Warn(err.Error())
		return
	}
	defer logF.Close()
	helper.ConfigureSlogger(cfg.Debug, logF)

	db, err := helper.LoadDB(cfg.DB)
	if err != nil {
		slog.Warn(err.Error())
		return
	}
	defer db.Close()

	rdb, err := helper.LoadRDB(cfg.Redis)
	if err != nil {
		slog.Warn(err.Error())
		return
	}
	defer rdb.Close()

	serve(cfg, db, rdb)
}

func serve(cfg *config.Config, db *sql.DB, rdb *redis.Client) {
	http.Handle("/", httpHandler(cfg, db, rdb))

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}

	go func() {
		err := server.ListenAndServe()
		slog.Info(err.Error())
	}()

	slog.Info("Сервер запущен", "address", server.Addr)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	slog.Info("Остановка сервера...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Ошибка при остановке сервера", "err", err)
	}
}

func httpHandler(cfg *config.Config, db *sql.DB, rdb *redis.Client) http.Handler {
	urlHelper := helper.NewURLHelper(cfg.Server, cfg.Path.UploadDirs)
	pathHelper := helper.NewPathHelper(cfg.Path)

	assessmentRepository := repository.NewAssessmentDBRepository(db)
	messageRepository := repository.NewMessageDBRepository(db)
	photoRepository := repository.NewPhotoDBRepository(db)
	questionnaireRepository := repository.NewQuestionnaireDBRepository(db, photoRepository)
	userDBRepository := repository.NewUserDBRepository(db)

	userRedisRepository := repository.NewUserRedisRepository(rdb)
	countryRepository := repository.NewCountryDBRepository(db)
	cityRepository := repository.NewCityDBRepository(db)
	coupleRepository := repository.NewCoupleDBRepository(db)

	assessmentService := service.NewAssessmentService(db, assessmentRepository, coupleRepository, questionnaireRepository)
	userService := service.NewUserService(userDBRepository, userRedisRepository)
	authService := service.NewAuthService(cfg.JWT, userDBRepository, userRedisRepository, userService)
	fileUploaderService := service.NewFileUploaderService(pathHelper, cfg.Path)
	messageService := service.NewMessageService(messageRepository, coupleRepository)
	photoService := service.NewPhotoService(pathHelper, photoRepository, questionnaireRepository, fileUploaderService)
	questionnaireService := service.NewQuestionnaireService(questionnaireRepository)

	authorizeMiddleware := middleware.NewAuthorizeMiddleware(authService)

	assessmentHandler := handler.NewAssessmentHandler(assessmentService)
	authHandler := handler.NewAuthHandler(authService)
	messageHandler := handler.NewMessageHandler(messageRepository, messageService)
	photoHandler := handler.NewPhotoHandler(urlHelper, photoRepository, photoService)
	questionnaireHandler := handler.NewQuestionnaireHandler(urlHelper, questionnaireRepository, questionnaireService)
	swaggerHandler := handler.NewSwaggerHandler(cfg.Path)
	userHandler := handler.NewUserHandler(userDBRepository, userService)
	dictionaryHandler := handler.NewDictionaryHandler(countryRepository, cityRepository)

	rc := router.NewConfigurator(
		cfg.Path,

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
