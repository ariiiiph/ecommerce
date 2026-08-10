package main

import (
	"fmt"
	"log"

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

	fmt.Println("E-commerce API is starting...")
	fmt.Println("Environment:", cfg.AppEnv)
	fmt.Println("PostgreSQL connection: OK")
	fmt.Println("Redis connection: OK")
}
