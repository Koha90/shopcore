package telegram

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLoadConfigFromEnv_Defaults(t *testing.T) {
	t.Setenv("TELEGRAM_PROXY_URL", "")
	t.Setenv("TELEGRAM_DEBUG", "")

	cfg := LoadConfigFromEnv()

	require.Equal(t, "", cfg.ProxyURL)
	require.Equal(t, 60*time.Second, cfg.PollTimeout)
	require.Equal(t, 10*time.Second, cfg.CheckInitTimeout)
	require.False(t, cfg.Debug)
}

func TestLoadConfigFromEnv_DebugAndProxy(t *testing.T) {
	t.Setenv("TELEGRAM_PROXY_URL", "socks5://127.0.0.1:1080")
	t.Setenv("TELEGRAM_DEBUG", "1")

	cfg := LoadConfigFromEnv()

	require.Equal(t, "socks5://127.0.0.1:1080", cfg.ProxyURL)
	require.Equal(t, 60*time.Second, cfg.PollTimeout)
	require.Equal(t, 10*time.Second, cfg.CheckInitTimeout)
	require.True(t, cfg.Debug)
}
