package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"botmanager/internal/domain"
)

func TestOrderService_CreateForVariant(t *testing.T) {
	product, err := domain.NewProduct(
		"Amnesia",
		1,
		"good stuff",
		"/tmp/img.png",
	)
	require.NoError(t, err)

	err = product.AddVariant("1g", 10, 1500)
	require.NoError(t, err)

	variants := product.Variants()
	require.Len(t, variants, 1)
	require.NoError(t, err)

	variant := variants[0]

	products := &stubProductReader{
		product: product,
	}
	orders := &stubOrderRepository{}
	users := &stubUserRepository{}
	tx := &stubTxManager{}
	bus := &stubEventBus{}

	svc := NewOrderService(products, orders, users, bus, tx, nil)

	order, err := svc.CreateForVariant(context.Background(), 42, product.ID(), variant.ID())
	require.NoError(t, err)
	require.NotNil(t, order)

	require.Equal(t, 1, tx.calls)
	require.Equal(t, 1, orders.saveCalls)
	require.NotNil(t, orders.savedOrder)
	require.Same(t, order, orders.savedOrder)
}

func TestOrderService_CreateForVariant_ProductError(t *testing.T) {
	products := &stubProductReader{
		err: errors.New("product load failed"),
	}
	orders := &stubOrderRepository{}
	users := &stubUserRepository{}
	tx := &stubTxManager{}
	bus := &stubEventBus{}

	svc := NewOrderService(products, orders, users, bus, tx, nil)

	order, err := svc.CreateForVariant(context.Background(), 42, 100, 200)

	require.Nil(t, order)
	require.EqualError(t, err, "load product: product load failed")

	require.Equal(t, 1, tx.calls)
	require.Equal(t, 0, orders.saveCalls)
}

func TestOrderService_CreateForVariant_SaveError(t *testing.T) {
	product, err := domain.NewProduct(
		"Amnesia",
		1,
		"good stuff",
		"/tmp/img.png",
	)
	require.NoError(t, err)

	err = product.AddVariant("1g", 10, 1500)
	require.NoError(t, err)

	variants := product.Variants()
	require.Len(t, variants, 1)
	require.NoError(t, err)

	variant := variants[0]

	products := &stubProductReader{
		product: product,
	}
	orders := &stubOrderRepository{
		saveErr: errors.New("save failed"),
	}
	users := &stubUserRepository{}
	tx := &stubTxManager{}
	bus := &stubEventBus{}

	svc := NewOrderService(products, orders, users, bus, tx, nil)

	order, err := svc.CreateForVariant(context.Background(), 42, product.ID(), variant.ID())

	require.Nil(t, order)
	require.EqualError(t, err, "save failed")

	require.Equal(t, 1, tx.calls)
	require.Equal(t, 1, orders.saveCalls)
	require.NotNil(t, orders.savedOrder)
}

func TestOrderService_ConfirmPayment(t *testing.T) {
	order, err := domain.NewOrder(
		42,
		[]domain.OrderItem{
			domain.NewOrderItem(1, 1, 1, 1500),
		},
		time.Now(),
	)
	require.NoError(t, err)

	products := &stubProductReader{}
	orders := &stubOrderRepository{
		order: order,
	}
	users := &stubUserRepository{}
	tx := &stubTxManager{}
	bus := &stubEventBus{}

	svc := NewOrderService(products, orders, users, bus, tx, nil)

	err = svc.ConfirmPayment(context.Background(), order.ID())
	require.NoError(t, err)

	require.Equal(t, 1, tx.calls)
	require.Equal(t, 1, orders.saveCalls)
	require.NotNil(t, orders.savedOrder)
	require.Equal(t, domain.OrderStatusPaid, orders.savedOrder.Status())
	require.Equal(t, 1, bus.calls)
}

func TestOrderService_ConfirmPayment_OrderNotFound(t *testing.T) {
	products := &stubProductReader{}
	orders := &stubOrderRepository{
		byIDErr: domain.ErrOrderNotFound,
	}
	users := &stubUserRepository{}
	tx := &stubTxManager{}
	bus := &stubEventBus{}

	svc := NewOrderService(products, orders, users, bus, tx, nil)

	err := svc.ConfirmPayment(context.Background(), 100)

	require.ErrorIs(t, err, domain.ErrOrderNotFound)
	require.Equal(t, 1, tx.calls)
	require.Equal(t, 0, orders.saveCalls)
	require.Equal(t, 0, bus.calls)
}

func TestOrderService_PayFromBalance(t *testing.T) {
	order, err := domain.NewOrder(
		42,
		[]domain.OrderItem{
			domain.NewOrderItem(1, 1, 1, 1500),
		},
		time.Now(),
	)
	require.NoError(t, err)

	user, err := domain.NewUser(domain.NewUserParams{
		Email:        "user@site.dev",
		PasswordHash: "hash",
	})
	require.NoError(t, err)

	err = user.AddBalance(2000)
	require.NoError(t, err)

	products := &stubProductReader{}
	orders := &stubOrderRepository{
		order: order,
	}
	users := &stubUserRepository{
		user: user,
	}
	tx := &stubTxManager{}
	bus := &stubEventBus{}

	svc := NewOrderService(products, orders, users, bus, tx, nil)

	err = svc.PayFromBalance(context.Background(), order.ID())
	require.NoError(t, err)

	require.Equal(t, 1, tx.calls)
	require.Equal(t, 1, users.saveCalls)
	require.Equal(t, 1, orders.saveCalls)

	require.NotNil(t, users.savedUser)
	require.NotNil(t, orders.savedOrder)

	require.Equal(t, domain.OrderStatusPaid, orders.savedOrder.Status())
	require.EqualValues(t, 500, users.savedUser.Balance())
}

func TestOrderService_PayFromBalance_InsufficientBalance(t *testing.T) {
	order, err := domain.NewOrder(
		42,
		[]domain.OrderItem{
			domain.NewOrderItem(1, 1, 1, 1500),
		},
		time.Now(),
	)
	require.NoError(t, err)

	user, err := domain.NewUser(domain.NewUserParams{
		Email:        "user@site.dev",
		PasswordHash: "hash",
	})
	require.NoError(t, err)

	products := &stubProductReader{}
	orders := &stubOrderRepository{
		order: order,
	}
	users := &stubUserRepository{
		user: user,
	}
	tx := &stubTxManager{}
	bus := &stubEventBus{}

	svc := NewOrderService(products, orders, users, bus, tx, nil)

	err = svc.PayFromBalance(context.Background(), order.ID())

	require.ErrorIs(t, err, domain.ErrInsufficientBalance)
	require.Equal(t, 1, tx.calls)
	require.Equal(t, 0, users.saveCalls)
	require.Equal(t, 0, orders.saveCalls)
	require.Equal(t, 0, bus.calls)
}

func TestOrderService_Cancel(t *testing.T) {
	order, err := domain.NewOrder(
		42,
		[]domain.OrderItem{
			domain.NewOrderItem(1, 1, 1, 1500),
		},
		time.Now(),
	)
	require.NoError(t, err)

	user, err := domain.NewUser(domain.NewUserParams{
		Email:        "user@site.dev",
		PasswordHash: "hash",
	})
	require.NoError(t, err)

	products := &stubProductReader{}
	orders := &stubOrderRepository{
		order: order,
	}
	users := &stubUserRepository{
		user: user,
	}
	tx := &stubTxManager{}
	bus := &stubEventBus{}

	svc := NewOrderService(products, orders, users, bus, tx, nil)

	err = svc.Cancel(context.Background(), order.ID())
	require.NoError(t, err)

	require.Equal(t, 1, tx.calls)
	require.Equal(t, 1, orders.saveCalls)
	require.NotNil(t, orders.savedOrder)
	require.Equal(t, domain.OrderStatusCancelled, orders.savedOrder.Status())
	require.Equal(t, 1, bus.calls)
}

func TestOrderService_Cancel_OrderNotFound(t *testing.T) {
	products := &stubProductReader{}
	orders := &stubOrderRepository{
		byIDErr: domain.ErrOrderNotFound,
	}
	users := &stubUserRepository{}
	tx := &stubTxManager{}
	bus := &stubEventBus{}

	svc := NewOrderService(products, orders, users, bus, tx, nil)

	err := svc.Cancel(context.Background(), 100)

	require.ErrorIs(t, err, domain.ErrOrderNotFound)
	require.Equal(t, 1, tx.calls)
	require.Equal(t, 0, orders.saveCalls)
	require.Equal(t, 0, bus.calls)
}
