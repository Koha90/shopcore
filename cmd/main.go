package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"botmanager/internal/botconfig"
	botpg "botmanager/internal/botconfig/postgres"
	"botmanager/internal/config"
	"botmanager/pkg/logger"
	"botmanager/pkg/migrator"
)

func main() {
	_ = godotenv.Load()
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbDatabase := os.Getenv("DB_DATABASE")

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser,
		dbPass,
		dbHost,
		dbPort,
		dbDatabase,
	)

	if err := migrator.MigratePostgres(dsn, "./migrations"); err != nil {
		log.Fatalf("Failed to migrate: %v", err)
	}

	cfg := config.MustLoad()

	logger, _ := logger.Setup(cfg.Env)
	logger.Debug("debug mode is enabled")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("failed to create pgx pool: %v", err)
	}

	store := botpg.NewStore(pool)

	_ = botconfig.NewService(
		store.BotRepository(),
		store.DatabaseProfileRepository(),
		logger,
	)

	log.Println("All migrations applied successfully!")
}
