package config

import "os"

type JWTConfig struct {
	SecretKey string
	Expire    int
}

func jwtFromEnv() *JWTConfig {
	return &JWTConfig{
		SecretKey: os.Getenv("JWT_SECRET_KEY"),
		Expire:    mustGetenvInt("JWT_EXPIRE"),
	}
}
