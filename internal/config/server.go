package config

import "os"

type ServerConfig struct {
	Host string
	Port int
}

func serverFromEnv() *ServerConfig {
	return &ServerConfig{
		Host: os.Getenv("SERVER_HOST"),
		Port: mustGetenvInt("SERVER_PORT"),
	}
}
