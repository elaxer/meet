package config

import "os"

type RedisConfig struct {
	Host      string
	Port      int
	Password  string
	Databases int
}

func redisFromEnv() *RedisConfig {
	return &RedisConfig{
		Host:      os.Getenv("REDIS_HOST"),
		Port:      mustGetenvInt("REDIS_PORT"),
		Password:  os.Getenv("REDIS_PASSWORD"),
		Databases: mustGetenvInt("REDIS_DATABASES"),
	}
}
