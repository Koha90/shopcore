package logger

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

// NewFileLogger creates a logger that writes only to the specified file.
func NewFileLogger(path string, level slog.Level) (*Logger, error) {
	return NewFileLoggerWithHandler(path, level, nil)
}

// NewFileLoggerWithHandler creates a file logger and optionally wraps its handler.
//
// The wrapper is useful for fan-out handlers, for example:
// writing to file and mirroring the same records into an in-memory TUI store.
func NewFileLoggerWithHandler(
	path string,
	level slog.Level,
	wrap func(next slog.Handler) slog.Handler,
) (*Logger, error) {
	if path == "" {
		return nil, fmt.Errorf("log file path is empty")
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("create log directory: %w", err)
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, fmt.Errorf("open log file: %w", err)
	}

	var handler slog.Handler = slog.NewTextHandler(file, &slog.HandlerOptions{
		AddSource: true,
		Level:     level,
	})

	if wrap != nil {
		handler = wrap(handler)
	}

	return &Logger{
		Logger: slog.New(handler),
		closer: file,
	}, nil
}
