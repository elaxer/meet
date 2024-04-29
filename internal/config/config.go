package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Debug  bool
	DB     *DBConfig
	Server *ServerConfig
	JWT    *JWTConfig
	Path   *PathConfig
	Redis  *RedisConfig
	TgBot  *TgBotConfig
}

func FromEnv(rootDir string) *Config {
	return &Config{
		mustGetenvBool("DEBUG"),
		dbFromEnv(),
		serverFromEnv(),
		jwtFromEnv(),
		pathFromEnv(rootDir),
		redisFromEnv(),
		tgBotFromEnv(),
	}
}

func mustGetenvBool(key string) bool {
	strVal := os.Getenv(key)

	intVal, err := strconv.Atoi(strVal)
	if err == nil {
		return intVal > 0
	}

	return strings.ToLower(strVal) == "true"
}

func mustGetenvInt(key string) int {
	val, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		panic(err)
	}

	return val
}
