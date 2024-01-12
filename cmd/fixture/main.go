package main

import (
	"database/sql"
	"log"
	"meet/internal/config"
	"meet/internal/pkg/app/fixture"
	"meet/internal/pkg/app/helper"
	"meet/internal/pkg/app/repository"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB

var dbConfig *config.DBConfig

var (
	userRepository          = repository.NewUserDBRepository(db)
	photoRepository         = repository.NewPhotoDBRepository(db)
	questionnaireRepository = repository.NewQuestionnaireDBRepository(db, photoRepository)
	coupleRepository        = repository.NewCoupleDBRepository(db)
)

func init() {
	_, b, _, _ := runtime.Caller(0)
	rootDir, _ := filepath.Abs(filepath.Dir(b) + "/../..")

	err := godotenv.Load(rootDir + "/.env")
	if err != nil {
		log.Fatal(err)
	}

	dbConfig = config.NewConfig(rootDir).DBConfig

	db, err = helper.LoadDB(dbConfig)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	if err := fixture.LoadFixtures(db, userRepository, questionnaireRepository, coupleRepository); err != nil {
		log.Fatal(err)
	}

	log.Println("Фикстуры успешно загружены в базу данных!")
}
