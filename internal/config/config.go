package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv string

	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
	PostgresPort     string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	return &Config{
		AppEnv:           os.Getenv("APP_ENV"),
		PostgresUser:     os.Getenv("POSTGRES_USER"),
		PostgresPassword: os.Getenv("POSTGRES_PASSWORD"),
		PostgresDB:       os.Getenv("POSTGRES_DB"),
		PostgresPort:     os.Getenv("POSTGRES_PORT"),
	}, nil
}
