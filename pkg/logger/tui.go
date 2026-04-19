package logger

import (
	"fmt"
	"log/slog"
	"os"
)

// SetupForTUI initializes application logger for TUI mode.
//
// Unlike Setup, this logger does not write to stdout becouse Bubble Tea owns
// the terminal screen. Writing regular logs to stdout during TUI rendering
// corrupts the interface.
//
// Local and dev environments use debug level.
// Prod uses info level.
//
// The logger is also installed as slog default logger so accidental
// slog.Default() calls stay out of stdout in TUI mode.
func SetupForTUI(env string) (*Logger, error) {
	if err := os.MkdirAll("logs", 0o755); err != nil {
		return nil, fmt.Errorf("cannot create logs directory: %w", err)
	}

	file, err := os.OpenFile("logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, fmt.Errorf("cannot open log file: %w", err)
	}

	level := slog.LevelInfo
	switch env {
	case envLocal, envDev:
		level = slog.LevelDebug
	case envProd:
		level = slog.LevelInfo
	}

	handler := slog.NewJSONHandler(file, &slog.HandlerOptions{
		AddSource: true,
		Level:     level,
	})

	base := slog.New(handler)
	slog.SetDefault(base)

	return &Logger{
		Logger: base,
		closer: file,
	}, nil
}
