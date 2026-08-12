package main

import (
	"fmt"
	"log"

	"github.com/ariiiiph/ecommerce/internal/app"
	"github.com/ariiiiph/ecommerce/internal/config"
	"github.com/ariiiiph/ecommerce/internal/db"
	"github.com/ariiiiph/ecommerce/internal/redis"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	database, err := db.NewPostgres(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	redisClient, err := redis.NewClient(
		cfg.Redis.Host,
		cfg.Redis.Port,
		cfg.Redis.Password,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer redisClient.Close()

	dependencies := &app.Dependencies{
		Config: cfg,
		DB:     database,
		Redis:  redisClient,
	}

	application := app.New(dependencies)

	fmt.Println("E-commerce API is starting...")
	fmt.Println("Environment:", application.Dependencies.Config.AppEnv)
	fmt.Println("PostgreSQL connection: OK")
	fmt.Println("Redis connection: OK")
}
