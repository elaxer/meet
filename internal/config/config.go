package config

import (
	"os"
	"path/filepath"
	"strconv"
)

const (
	uploadsDir      = "/uploads"
	photosDir       = "/photos"
	swaggerFilePath = "/api/swagger.yml"
)

type PathConfig struct {
	RootDir         string
	UploadDirs      *UploadDirs
	SwaggerFilePath string
}

type UploadDirs struct {
	UploadDir string
	PhotoDir  string
}

func (pc *PathConfig) FullPath(paths ...string) string {
	path := filepath.Join(paths...)

	return filepath.Join(pc.RootDir, path)
}

type Config struct {
	DBConfig     *DBConfig
	ServerConfig *ServerConfig
	JWTConfig    *JWTConfig
	PathConfig   *PathConfig
}

func NewConfig(rootDir string) *Config {
	return &Config{
		newDBConfig(),
		&ServerConfig{
			Host: os.Getenv("SERVER_HOST"),
			Port: mustGetenvInt("SERVER_PORT"),
		},
		&JWTConfig{
			SecretKey: os.Getenv("JWT_SECRET_KEY"),
			Expire:    mustGetenvInt("JWT_EXPIRE"),
		},
		&PathConfig{
			rootDir,
			&UploadDirs{uploadsDir, photosDir},
			swaggerFilePath,
		},
	}
}

func mustGetenvInt(key string) int {
	val, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		panic(err)
	}

	return val
}
