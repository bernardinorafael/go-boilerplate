package main

import (
	"gulg/internal/infra/database/pg"
	"gulg/pkg/config"
	"log"
)

func main() {
	cfg := config.GetConfig()

	conn, err := pg.NewConnection(cfg.PostgresDSN)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer conn.Close()
}
