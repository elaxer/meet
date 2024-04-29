package main

import (
	"log/slog"
	"meet/internal/config"
	"meet/internal/pkg/app/fixture"
	"meet/internal/pkg/app/helper"
	"meet/internal/pkg/app/repository"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	_, b, _, _ := runtime.Caller(0)
	rootDir, _ := filepath.Abs(filepath.Dir(b) + "/../..")

	err := godotenv.Load(rootDir + "/.env")
	if err != nil {
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

	userRepository := repository.NewUserDBRepository(db)
	questionnaireRepository := repository.NewQuestionnaireDBRepository(db, repository.NewPhotoDBRepository(db))
	coupleRepository := repository.NewCoupleDBRepository(db)

	if err := fixture.LoadFixtures(db, userRepository, questionnaireRepository, coupleRepository); err != nil {
		slog.Warn(err.Error())
		return
	}

	slog.Info("Фикстуры успешно загружены в базу данных!")
}
