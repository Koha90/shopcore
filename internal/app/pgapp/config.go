package pgapp

import (
	"errors"
	"fmt"
	"os"
	"strings"
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
func LoadConfigFromEnv() (Config, error) {
	cfg := Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     envOrDefault("DB_PORT", "5432"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: os.Getenv("DB_DATABASE"),
		SSLMode:  envOrDefault("DB_SSLMODE", "disable"),
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

// Validate checks whether the config contains required fields.
func (c Config) Validate() error {
	switch {
	case strings.TrimSpace(c.Host) == "":
		return errors.New("DB_HOST is required")
	case strings.TrimSpace(c.Port) == "":
		return errors.New("DB_PORT is required")
	case strings.TrimSpace(c.User) == "":
		return errors.New("DB_USER is required")
	case strings.TrimSpace(c.Password) == "":
		return errors.New("DB_PASSWORD is required")
	case strings.TrimSpace(c.Database) == "":
		return errors.New("DB_DATABASE is required")
	default:
		return nil
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

func envOrDefault(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
