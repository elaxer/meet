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

func NewConfig() (*Config, error) {
	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		return nil, err
	}
	jwtExpire, err := strconv.Atoi(os.Getenv("JWT_EXPIRE"))
	if err != nil {
		return nil, err
	}

	c := &Config{
		&DBConfig{
			DriverName: os.Getenv("DB_DRIVER_NAME"),
			Host:       os.Getenv("DB_HOST"),
			Port:       dbPort,
			SSLMode:    os.Getenv("DB_SSL_MODE"),
			User:       os.Getenv("DB_USER"),
			Password:   os.Getenv("DB_PASSWORD"),
			DBName:     os.Getenv("DB_NAME"),
		},
		&ServerConfig{
			Port: os.Getenv("SERVER_PORT"),
		},
		&JWTConfig{
			SecretKey: os.Getenv("JWT_SECRET_KEY"),
			Expire:    jwtExpire,
		},
		os.Getenv("UPLOAD_DIR"),
		swaggerFile,
	}

	return c, nil
}
