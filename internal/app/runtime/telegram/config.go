package telegram

import (
	"errors"
	"os"
	"strings"
	"time"
)

// Config contains telegram runtime settings.
type Config struct {
	Token            string
	Debug            bool
	ProxyURL         string
	PollTimeout      time.Duration
	CheckInitTimeout time.Duration
}

// LoadConfigFromEnv loads Telegram runtime config from environment.
func LoadConfigFromEnv() (Config, error) {
	cfg := Config{
		Token:            os.Getenv("TELEGRAM_BOT_TOKEN"),
		ProxyURL:         os.Getenv("TELEGRAM_PROXY_URL"),
		PollTimeout:      60 * time.Second,
		CheckInitTimeout: 10 * time.Second,
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

// Validate checks required Telegram settings.
func (c Config) Validate() error {
	if strings.TrimSpace(c.Token) == "" {
		return errors.New("TELEGRAM_BOT_TOKEN is required")
	}
	return nil
}
