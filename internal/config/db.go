package config

import "os"

type DBConfig struct {
	DriverName string
	Host       string
	Port       int
	SSLMode    string
	User       string
	Password   string
	DBName     string
}

func newDBConfig() *DBConfig {
	return &DBConfig{
		DriverName: os.Getenv("DB_DRIVER_NAME"),
		Host:       os.Getenv("DB_HOST"),
		Port:       mustGetenvInt("DB_PORT"),
		SSLMode:    os.Getenv("DB_SSL_MODE"),
		User:       os.Getenv("DB_USER"),
		Password:   os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
	}
}
