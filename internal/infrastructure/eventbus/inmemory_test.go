package eventbus

import (
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"

	"botmanager/internal/domain"
)

func TestInMemoryBus_MultipleEvents(t *testing.T) {
	bus := New(slog.Default())

	count := 0

	e1 := domain.NewOrderPaid(1)
	e2 := domain.NewOrderPaid(2)

	bus.Subscribe(
		e1.Name(),
		func(ctx context.Context, event domain.Event) error {
			count++
			return nil
		},
	)

	err := bus.Publish(context.Background(), e1, e2)
	require.NoError(t, err)
	require.Equal(t, 2, count)
}

func TestInMemoryBus_MultipleHandlers(t *testing.T) {
	bus := New(slog.Default())

	count := 0
	event := domain.NewOrderPaid(1)

	bus.Subscribe(event.Name(), func(ctx context.Context, event domain.Event) error {
		count++
		return nil
	})

	bus.Subscribe(event.Name(), func(ctx context.Context, event domain.Event) error {
		count++
		return nil
	})

	err := bus.Publish(context.Background(), event)
	require.NoError(t, err)
	require.Equal(t, 2, count)
}

func TestInMemoryBus_HandlerErrorDoesNotStopDelivery(t *testing.T) {
	bus := New(slog.Default())

	count := 0
	event := domain.NewOrderPaid(1)

	bus.Subscribe(event.Name(), func(ctx context.Context, event domain.Event) error {
		return errors.New("handler failed")
	})

	bus.Subscribe(event.Name(), func(ctx context.Context, event domain.Event) error {
		count++
		return nil
	})

	err := bus.Publish(context.Background(), event)
	require.NoError(t, err)
	require.Equal(t, 1, count)
}

func TestInMemoryBus_PublishWithoutSubscribers(t *testing.T) {
	bus := New(slog.Default())

	event := domain.NewOrderPaid(1)

	err := bus.Publish(context.Background(), event)
	require.NoError(t, err)
}
