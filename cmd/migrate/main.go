package main

import (
	"log"

	"github.com/joho/godotenv"

	"github.com/koha90/shopcore/internal/app/pgapp"
	"github.com/koha90/shopcore/pkg/migrator"
)

func main() {
	_ = godotenv.Load()

	cfg, err := pgapp.LoadConfigFromEnv()
	if err != nil {
		log.Fatalf("failed to load config from env: %v", err)
	}

	if err := migrator.MigratePostgres(cfg.DSN(), "./migrations"); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	log.Println("all migrations applied successfully")
}
