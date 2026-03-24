package tuiapp

import (
	"time"

	"github.com/koha90/shopcore/internal/app/pgapp"
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
func LoadConfigFromEnv() (Config, error) {
	pgcfg, err := pgapp.LoadConfigFromEnv()
	if err != nil {
		return Config{}, err
	}
	return Config{
		Postgres:      pgcfg,
		OpenDBTimeout: 10 * time.Second,
	}, nil
}
