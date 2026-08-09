package main

import (
	"fmt"
	"log"

	"github.com/ariiiiph/ecommerce/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("E-commerce API is starting...")
	fmt.Println("Environment:", cfg.AppEnv)

}
