package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Environment string `env:"ENV" env-default:"development"`
	AppEnv      string

	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
	PostgresPort     string

	Redis RedisConfig
	JWT   JWTConfig
}

type RedisConfig struct {
	Host     string `env:"REDIS_HOST" env-default:"localhost"`
	Port     int    `env:"REDIS_PORT" env-default:"6379"`
	Password string `env:"REDIS_PASSWORD"`
	DB       int    `env:"REDIS_DB" env-default:"0"`
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	return &Config{
		AppEnv: os.Getenv("APP_ENV"),

		PostgresUser:     os.Getenv("POSTGRES_USER"),
		PostgresPassword: os.Getenv("POSTGRES_PASSWORD"),
		PostgresDB:       os.Getenv("POSTGRES_DB"),
		PostgresPort:     os.Getenv("POSTGRES_PORT"),

		Redis: RedisConfig{
			Host:     os.Getenv("REDIS_HOST"),
			Port:     getEnvInt("REDIS_PORT", 6379),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       getEnvInt("REDIS_DB", 0),
		},
	}, nil
}

func getEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	result, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return result
}
