package telegram

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/koha90/shopcore/internal/flow"
	"github.com/koha90/shopcore/internal/manager"
	ordersvc "github.com/koha90/shopcore/internal/order/service"
)

type orderServiceStub struct {
	called       bool
	createParams ordersvc.CreateOrderParams
	createResult ordersvc.CreateResult
	order        ordersvc.Order
	err          error
}

func (s *orderServiceStub) Create(ctx context.Context, params ordersvc.CreateOrderParams) (ordersvc.CreateResult, error) {
	s.called = true
	s.createParams = params
	return s.createResult, s.err
}

func (s *orderServiceStub) ByID(ctx context.Context, id int64) (ordersvc.Order, error) {
	return s.order, s.err
}

func (s *orderServiceStub) UpdateStatus(ctx context.Context, id int64, status ordersvc.OrderStatus) error {
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

	store.Put(key, flow.Session{
		Current: flow.ScreenOrderConfirm,
		History: []flow.ScreenID{
			flow.ScreenID("catalog:screen:city=moscow;category=flowers;district=center;product=rose-box;variant=small"),
		},
	})

	orders := &orderServiceStub{
		createResult: ordersvc.CreateResult{
			ID:     42,
			Status: ordersvc.OrderStatusNew,
		},
		order: ordersvc.Order{
			ID:           42,
			BotID:        "shop-main",
			BotName:      "Shop Main",
			ChatID:       101,
			UserID:       202,
			UserName:     "Алексей",
			UserUsername: "koha90",
			CityID:       "moscow",
			CityName:     "Москва",
			DistrictID:   "center",
			DistrictName: "Центр",
			ProductID:    "rose-box",
			ProductName:  "Rose Box",
			VariantID:    "small",
			VariantName:  "S / 9 шт",
			PriceText:    "2500 ₽",
			Status:       ordersvc.OrderStatusNew,
		},
	}

	runner := NewRunnerWithDeps(
		Config{},
		nil,
		nil,
		func(spec manager.BotSpec) (ordersvc.RuntimeService, error) {
			return orders, nil
		},
		nil,
		nil,
	)

	got, err := runner.persistConfirmedOrder(
		context.Background(),
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

	require.True(t, orders.called)
	require.Equal(t, "shop-main", orders.createParams.BotID)
	require.Equal(t, int64(101), orders.createParams.ChatID)
	require.Equal(t, int64(202), orders.createParams.UserID)

	require.Equal(t, int64(42), got.ID)
	require.Equal(t, ordersvc.OrderStatusNew, got.Status)
	require.Equal(t, "Москва", got.CityName)
}

