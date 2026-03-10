package service

import (
	"context"

	"botmanager/internal/domain"
)

type stubTxManager struct {
	err   error
	calls int
}

func (s *stubTxManager) WithinTransaction(
	ctx context.Context,
	fn func(ctx context.Context) error,
) error {
	s.calls++
	if s.err != nil {
		return s.err
	}
	return fn(ctx)
}

type stubEventBus struct {
	err   error
	calls int
}

func (s *stubEventBus) Publish(ctx context.Context, events ...domain.Event) error {
	s.calls++
	return s.err
}

type stubUserRepository struct {
	user      *domain.User
	byIDErr   error
	saveErr   error
	savedUser *domain.User
	saveCalls int
}

func (r *stubUserRepository) ByID(ctx context.Context, id int) (*domain.User, error) {
	return r.user, r.byIDErr
}

func (r *stubUserRepository) Save(ctx context.Context, user *domain.User) error {
	r.saveCalls++
	r.savedUser = user
	return r.saveErr
}

type stubProductReader struct {
	product *domain.Product
	err     error
}

func (r *stubProductReader) ByID(ctx context.Context, id int) (*domain.Product, error) {
	return r.product, r.err
}

type stubOrderRepository struct {
	order      *domain.Order
	byIDErr    error
	saveErr    error
	saveCalls  int
	savedOrder *domain.Order
}

func (r *stubOrderRepository) ByID(ctx context.Context, id int) (*domain.Order, error) {
	return r.order, r.byIDErr
}

func (r *stubOrderRepository) Save(ctx context.Context, order *domain.Order) error {
	r.saveCalls++
	r.savedOrder = order
	return r.saveErr
}
