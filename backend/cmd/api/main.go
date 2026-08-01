package main

import (
	"log"

	"pascal-tech-dev/calczzle-backend/internal/app"
	"pascal-tech-dev/calczzle-backend/internal/config"
)

func main() {
	cfg := config.Load()
	server := app.New(cfg)

	if err := server.Listen(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
