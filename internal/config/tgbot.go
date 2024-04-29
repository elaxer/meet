package config

import "os"

type TgBotConfig struct {
	Token         string
	UpdateTimeout int
}

func tgBotFromEnv() *TgBotConfig {
	return &TgBotConfig{
		Token:         os.Getenv("TG_BOT_TOKEN"),
		UpdateTimeout: mustGetenvInt("TG_BOT_UPDATE_TIMEOUT"),
	}
}
