package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"botmanager/pkg/migrator"
)

func main() {
	_ = godotenv.Load()

	dsn := mustPostgresDSN()

	if err := migrator.MigratePostgres(dsn, "./migrations"); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	log.Println("all migrations applied successfully")
}

// mustPostgresDSN builds PostgreSQL DSN from environment variables.
//
// Required environment variables:
//   - DB_HOST
//   - DB_PORT
//   - DB_USER
//   - DB_PASSWORD
//   - DB_DATABASE
//
// Optional:
//   - DB_SSLMODE (defaults to disable)
func mustPostgresDSN() string {
	host := mustEnv("DB_HOST")
	port := mustEnv("DB_PORT")
	user := mustEnv("DB_USER")
	password := mustEnv("DB_PASSWORD")
	database := mustEnv("DB_DATABASE")
	sslmode := envOrDefault("DB_SSLMODE", "disable")

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user,
		password,
		host,
		port,
		database,
		sslmode,
	)
}

// mustEnv returns required environment variable or exits process immediately.
func mustEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("required environment variable %s is not set", key)
	}
	return value
}

// envOrDefault returns environment variable value or fallback if empty.
func envOrDefault(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
