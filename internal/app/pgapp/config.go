package pgapp

import (
	"fmt"
	"os"
)

// Config contains PostgreSQL connection parameters required by application
// entrypoints.
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

// LoadConfigFromEnv loads PostgreSQL connection config from environment.
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
//
// The function panics if a required environment variable is missing.
func LoadConfigFromEnv() Config {
	return Config{
		Host:     mustEnv("DB_HOST"),
		Port:     mustEnv("DB_PORT"),
		User:     mustEnv("DB_USER"),
		Password: mustEnv("DB_PASSWORD"),
		Database: mustEnv("DB_DATABASE"),
		SSLMode:  envOrDefault("DB_SSLMODE", "disable"),
	}
}

// DSN returns PostgreSQL connection string.
func (c Config) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
		c.SSLMode,
	)
}

func mustEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("required environment variable %s is not set", key))
	}
	return value
}

func envOrDefault(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
