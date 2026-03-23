package main

import (
	"log"

	"github.com/joho/godotenv"

	"botmanager/internal/app/pgapp"
	"botmanager/pkg/migrator"
)

func main() {
	_ = godotenv.Load()

	cfg := pgapp.LoadConfigFromEnv()

	if err := migrator.MigratePostgres(cfg.DSN(), "./migrations"); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	log.Println("all migrations applied successfully")
}
