package services

import (
	"context"

	"botmanager/internal/domain"
)

// EventBus defines in-process event publishing mechanism.
//
// It allow application services to publish domain events
// and register handlers reacting to them
type EventBus interface {
	// Publish dispatches domain events to subscribe handlers.
	Publish(ctx context.Context, events ...domain.DomainEvent) error

	// Subscribe register a handler for specific event name.
	Subscribe(eventName string, handler EventHandler)
}

type EventHandler func(ctx context.Context, event domain.DomainEvent) error
