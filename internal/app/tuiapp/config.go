package tuiapp

import (
	"time"

	"botmanager/internal/app/pgapp"
)

// Config contains infrastructures settings required to build the TUI
// application container.
type Config struct {
	Postgres      pgapp.Config
	OpenDBTimeout time.Duration
}

// LoadConfigFromEnv loads tuiapp configuration from environment.
//
// It delegates PostgreSQL parsing to pgapp.LoadConfigFromEnv.
func LoadConfigFromEnv() Config {
	return Config{
		Postgres:      pgapp.LoadConfigFromEnv(),
		OpenDBTimeout: 10 * time.Second,
	}
}
