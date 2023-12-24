package app

import (
	"os"
	"strconv"
)

const (
	swaggerFile = "api/swagger.yml"
)

type DBConfig struct {
	DriverName string
	Host       string
	Port       int
	SSLMode    string
	User       string
	Password   string
	DBName     string
}

type ServerConfig struct {
	Host string
	Port string
}

type JWTConfig struct {
	SecretKey string
	Expire    int
}

type Config struct {
	DBConfig     *DBConfig
	ServerConfig *ServerConfig
	JWTConfig    *JWTConfig
	UploadDir    string
	SwaggerFile  string
}

func NewConfig() *Config {
	return &Config{
		&DBConfig{
			DriverName: os.Getenv("DB_DRIVER_NAME"),
			Host:       os.Getenv("DB_HOST"),
			Port:       mustGetenvInt("DB_PORT"),
			SSLMode:    os.Getenv("DB_SSL_MODE"),
			User:       os.Getenv("DB_USER"),
			Password:   os.Getenv("DB_PASSWORD"),
			DBName:     os.Getenv("DB_NAME"),
		},
		&ServerConfig{
			Host: os.Getenv("SERVER_HOST"),
			Port: os.Getenv("SERVER_PORT"),
		},
		&JWTConfig{
			SecretKey: os.Getenv("JWT_SECRET_KEY"),
			Expire:    mustGetenvInt("JWT_EXPIRE"),
		},
		os.Getenv("UPLOAD_DIR"),
		swaggerFile,
	}
}

func mustGetenvInt(key string) int {
	val, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		panic(err)
	}

	return val
}
