// Package eventbus provides in-memory event bus implementation.
package eventbus

import (
	"context"
	"log/slog"
	"sync"

	"botmanager/internal/domain"
	"botmanager/internal/service"
)

// InMemoryBus delivers domain events to subscribed handlers
// within current process memory.
//
// It is safe for concurrent use.
type InMemoryBus struct {
	mu       sync.RWMutex
	handlers map[string][]service.EventHandler
	logger   *slog.Logger
}

// New creates a new in-memory event bus.
//
// logger may be nil. In that case slog.Default() is used.
func New(logger *slog.Logger) *InMemoryBus {
	if logger == nil {
		logger = slog.Default()
	}

	return &InMemoryBus{
		handlers: make(map[string][]service.EventHandler),
		logger:   logger,
	}
}

// Subscribe registers handler for events with the provided name.
func (b *InMemoryBus) Subscribe(name string, h service.EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.handlers[name] = append(b.handlers[name], h)
}

// Publish delivers provided events to all subscribed handlers.
//
// Handler errors are logged and do not stop delivery to other handlers.
func (b *InMemoryBus) Publish(ctx context.Context, events ...domain.Event) error {
	for _, e := range events {
		name := e.Name()

		b.mu.RLock()
		handlers := b.handlers[name]
		b.mu.RUnlock()

		for _, h := range handlers {
			if err := h(ctx, e); err != nil {
				b.logger.Error(
					"event handler failed",
					"event", name,
					"error", err,
				)
			}
		}
	}

	return nil
}
