package config

import (
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type Config struct {
	AppAddress string `env:"APP_ADDRESS"`

	DatabaseHost     string `env:"DATABASE_HOST"`
	DatabaseName     string `env:"DATABASE_NAME"`
	DatabaseUser     string `env:"DATABASE_USER"`
	DatabasePassword string `env:"DATABASE_PASSWORD"`
	DatabasePort     int    `env:"DATABASE_PORT"`
}

func NewConfig(isDebug bool) (*Config, error) {
	if isDebug {
		err := godotenv.Load(".env")
		if err != nil {
			return nil, err
		}
	}

	var config Config
	err := env.Parse(&config)

	if err != nil {
		return nil, err
	}

	return &config, nil
}
