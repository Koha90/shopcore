package telegram

import (
	"os"
	"time"
)

// Config contains telegram runtime settings.
type Config struct {
	ProxyURL         string
	PollTimeout      time.Duration
	CheckInitTimeout time.Duration
	Debug            bool
}

// LoadConfigFromEnv loads Telegram runtime config from environment.
func LoadConfigFromEnv() Config {
	return Config{
		ProxyURL:         os.Getenv("TELEGRAM_PROXY_URL"),
		PollTimeout:      60 * time.Second,
		CheckInitTimeout: 10 * time.Second,
		Debug:            os.Getenv("TELEGRAM_DEBUG") == "1",
	}
}
