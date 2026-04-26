package telegram

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/koha90/shopcore/internal/flow"
	"github.com/koha90/shopcore/internal/manager"
	ordersvc "github.com/koha90/shopcore/internal/order/service"
)

type orderCreatorStub struct {
	params ordersvc.CreateOrderParams
	err    error
}

func (s *orderCreatorStub) Create(ctx context.Context, params ordersvc.CreateOrderParams) error {
	s.params = params
	return s.err
}

func TestRunnerPersistConfirmedOrder(t *testing.T) {
	t.Parallel()

	store := flow.NewMemoryStore()
	svc := flow.NewService(store)

	key := flow.SessionKey{
		BotID:  "shop-main",
		ChatID: 101,
		UserID: 202,
	}

	ctx := context.Background()

	_, err := svc.Start(ctx, flow.StartRequest{
		SessionKey:    key,
		StartScenario: "inline_catalog",
	})
	require.NoError(t, err)

	// use existing helper from flow tests or inline drill-down if you already have one
	// to reach variant leaf and order confirm screen

	creator := &orderCreatorStub{}
	runner := NewRunnerWithDeps(Config{}, nil, nil, func(spec manager.BotSpec) (ordersvc.OrderCreator, error) {
		return creator, nil
	}, nil)

	err = runner.persistConfirmedOrder(
		ctx,
		manager.BotSpec{
			ID:   "shop-main",
			Name: "Shop Main",
		},
		svc,
		key,
		OrderNotificationMeta{
			UserID:    202,
			ChatID:    101,
			UserName:  "Алексей",
			UserLogin: "koha90",
		},
	)
	require.NoError(t, err)

	require.Equal(t, "shop-main", creator.params.BotID)
	require.Equal(t, int64(101), creator.params.ChatID)
	require.Equal(t, int64(202), creator.params.UserID)
}
