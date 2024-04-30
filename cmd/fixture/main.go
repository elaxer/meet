package main

import (
	"log/slog"
	"meet/internal/config"
	"meet/internal/pkg/app/database"
	"meet/internal/pkg/app/fixture"
	"meet/internal/pkg/app/repository/dbrepository"
	"meet/internal/pkg/app/slogger"
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
		panic(err)
	}

	cfg := config.FromEnv(rootDir)

	logF := slogger.MustOpenLog(rootDir)
	defer logF.Close()

	slogger.Configure(cfg.Debug, logF)

	db := database.MustConnect(cfg.DB)

	var (
		userRepository          = dbrepository.NewUserRepository(db)
		questionnaireRepository = dbrepository.NewQuestionnaireRepository(db)
		coupleRepository        = dbrepository.NewCoupleRepository(db)
	)

	if err := fixture.LoadFixtures(db, userRepository, questionnaireRepository, coupleRepository); err != nil {
		panic(err)
	}

	slog.Info("Фикстуры успешно загружены в базу данных!")
}
