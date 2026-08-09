package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	return &Config{
		AppEnv: os.Getenv("APP_ENV"),
	}, nil
}
