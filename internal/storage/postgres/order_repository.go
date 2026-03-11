package postgres

import (
	"context"

	"botmanager/internal/domain"
	"botmanager/internal/service"
)

var _ service.OrderRepository = (*OrderRepo)(nil)

type OrderRepo struct{}

func NewOrderRepo() *OrderRepo {
	return &OrderRepo{}
}

func (r *OrderRepo) Save(ctx context.Context, order *domain.Order) error {
	if order.ID() == 0 {
		return r.Create(ctx, order)
	}
	return r.Update(ctx, order)
}

func (r *OrderRepo) Create(ctx context.Context, order *domain.Order) error {
	panic("unimplemented")
}

// ByID implements [service.OrderRepository].
func (r *OrderRepo) ByID(ctx context.Context, id int) (*domain.Order, error) {
	panic("unimplemented")
}

// ListByCustomer implements [service.OrderRepository].
func (r *OrderRepo) ListByCustomer(ctx context.Context, customerID int) ([]domain.Order, error) {
	panic("unimplemented")
}

// Update implements [service.OrderRepository].
func (r *OrderRepo) Update(ctx context.Context, order *domain.Order) error {
	panic("unimplemented")
}
