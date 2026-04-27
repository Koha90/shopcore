package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

type repositoryStub struct {
	record      OrderRecord
	createRes   CreateResult
	order       Order
	updateID    int64
	updateState OrderStatus
	err         error
}

func (s *repositoryStub) Create(ctx context.Context, record OrderRecord) (CreateResult, error) {
	s.record = record
	return s.createRes, s.err
}

func (s *repositoryStub) ByID(ctx context.Context, id int64) (Order, error) {
	return s.order, s.err
}

func (s *repositoryStub) UpdateStatus(ctx context.Context, id int64, status OrderStatus) error {
	s.updateID = id
	s.updateState = status
	return s.err
}

func TestServiceCreate(t *testing.T) {
	t.Parallel()

	repo := &repositoryStub{
		createRes: CreateResult{
			ID:     42,
			Status: OrderStatusNew,
		},
	}
	svc := New(repo)

	got, err := svc.Create(context.Background(), CreateOrderParams{
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
		VariantID:    "large",
		VariantName:  "L / 25 шт",
		PriceText:    "5900 ₽",
	})
	require.NoError(t, err)

	require.Equal(t, int64(42), got.ID)
	require.Equal(t, OrderStatusNew, got.Status)
	require.Equal(t, OrderStatusNew, repo.record.Status)
}

func TestServiceByID(t *testing.T) {
	t.Parallel()

	repo := &repositoryStub{
		order: Order{
			ID:     42,
			Status: OrderStatusNew,
		},
	}
	svc := New(repo)

	got, err := svc.ByID(context.Background(), 42)
	require.NoError(t, err)
	require.Equal(t, int64(42), got.ID)
}

func TestServiceUpdateStatus(t *testing.T) {
	t.Parallel()

	repo := &repositoryStub{
		order: Order{
			ID:     42,
			Status: OrderStatusNew,
		},
	}
	svc := New(repo)

	err := svc.UpdateStatus(context.Background(), 42, OrderStatusInProgress)
	require.NoError(t, err)

	require.Equal(t, int64(42), repo.updateID)
	require.Equal(t, OrderStatusInProgress, repo.updateState)
}

func TestServiceUpdateStatus_InvalidStatus(t *testing.T) {
	t.Parallel()

	repo := &repositoryStub{}
	svc := New(repo)

	err := svc.UpdateStatus(context.Background(), 42, "abracadabra")
	require.ErrorIs(t, err, ErrOrderStatusInvalid)
}

func TestServiceUpdateStatus_AllowsNewToInProgress(t *testing.T) {
	t.Parallel()

	repo := &repositoryStub{
		order: Order{
			ID:     42,
			Status: OrderStatusNew,
		},
	}
	svc := New(repo)

	err := svc.UpdateStatus(context.Background(), 42, OrderStatusInProgress)
	require.NoError(t, err)
	require.Equal(t, int64(42), repo.updateID)
	require.Equal(t, OrderStatusInProgress, repo.updateState)
}

func TestServiceUpdateStatus_AllowsNewToClosed(t *testing.T) {
	t.Parallel()

	repo := &repositoryStub{
		order: Order{
			ID:     42,
			Status: OrderStatusNew,
		},
	}
	svc := New(repo)

	err := svc.UpdateStatus(context.Background(), 42, OrderStatusClosed)
	require.NoError(t, err)
	require.Equal(t, OrderStatusClosed, repo.updateState)
}

func TestServiceUpdateStatus_RejectsClosedToInProgress(t *testing.T) {
	t.Parallel()

	repo := &repositoryStub{
		order: Order{
			ID:     42,
			Status: OrderStatusClosed,
		},
	}
	svc := New(repo)

	err := svc.UpdateStatus(context.Background(), 42, OrderStatusInProgress)
	require.ErrorIs(t, err, ErrOrderStatusTransitionDead)
}
