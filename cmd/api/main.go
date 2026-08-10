package main

import (
	"fmt"
	"log"

	"github.com/ariiiiph/ecommerce/internal/config"
	"github.com/ariiiiph/ecommerce/internal/db"
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

	fmt.Println("E-commerce API is starting...")
	fmt.Println("Environment:", cfg.AppEnv)
	fmt.Println("PostgreSQL connection: OK")

}
