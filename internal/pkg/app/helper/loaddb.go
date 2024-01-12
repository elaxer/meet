package helper

import (
	"database/sql"
	"fmt"
	"meet/internal/config"
)

func LoadDB(dbConfig *config.DBConfig) (*sql.DB, error) {
	db, err := sql.Open(
		dbConfig.DriverName,
		fmt.Sprintf(
			"host=%s port=%d sslmode=%s user=%s password=%s dbname=%s",
			dbConfig.Host,
			dbConfig.Port,
			dbConfig.SSLMode,
			dbConfig.User,
			dbConfig.Password,
			dbConfig.DBName,
		),
	)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, err
}
