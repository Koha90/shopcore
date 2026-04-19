// Package runtimelog stores structured runtime log entries for TUI views.
//
// The package is intentionally transport-agnostic.
// Runtimes may write records here, and TUI may read them for per-bot log panel.
package runtimelog

import (
	"log/slog"
	"time"
)

// Field describes one flattened structured log attribute.
type Field struct {
	Key   string
	Value string
}

// Entry represents one runtime one log record prepared for UI consumption.
//
// BotID is the primary routing key used by TUI to show logs for one bot.
// BotName is optional and may be used for display.
// Fields keeps additional structured attributes in flattened form.
type Entry struct {
	Time    time.Time
	Level   slog.Level
	BotID   string
	BotName string
	Message string
	Fields  []Field
}
