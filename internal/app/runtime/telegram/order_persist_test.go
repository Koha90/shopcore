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
	called bool
	params ordersvc.CreateOrderParams
	err    error
}

func (s *orderCreatorStub) Create(ctx context.Context, params ordersvc.CreateOrderParams) error {
	s.called = true
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

	// Prepare order-confirmation state with catalog leaf in history.
	store.Put(key, flow.Session{
		Current: flow.ScreenOrderConfirm,
		History: []flow.ScreenID{
			flow.ScreenID("catalog:screen:city=moscow;category=flowers;district=center;product=rose-box;variant=small"),
		},
	})

	ctx := context.Background()

	creator := &orderCreatorStub{}
	runner := NewRunnerWithDeps(
		Config{},
		nil,
		nil,
		func(spec manager.BotSpec) (ordersvc.OrderCreator, error) {
			return creator, nil
		},
		nil,
	)

	err := runner.persistConfirmedOrder(
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

	require.True(t, creator.called)
	require.Equal(t, "shop-main", creator.params.BotID)
	require.Equal(t, "Shop Main", creator.params.BotName)
	require.Equal(t, int64(101), creator.params.ChatID)
	require.Equal(t, int64(202), creator.params.UserID)
	require.Equal(t, "Алексей", creator.params.UserName)
	require.Equal(t, "koha90", creator.params.UserUsername)

	require.Equal(t, "moscow", creator.params.CityID)
	require.Equal(t, "Москва", creator.params.CityName)
	require.Equal(t, "center", creator.params.DistrictID)
	require.Equal(t, "Центр", creator.params.DistrictName)
	require.Equal(t, "rose-box", creator.params.ProductID)
	require.Equal(t, "Rose Box", creator.params.ProductName)
	require.Equal(t, "small", creator.params.VariantID)
	require.Equal(t, "S / 9 шт", creator.params.VariantName)
	require.Equal(t, "2500 ₽", creator.params.PriceText)
}

