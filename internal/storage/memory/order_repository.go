package memory

import (
	"context"
	"sync"

	"botmanager/internal/domain"
)

// OrderRepository stores orders in process memory.
//
// It is intended for local development, tests and simple runtime scenarios.
// Repository assigns incremental IDs to new orders on first save.
type OrderRepository struct {
	mu     sync.Mutex
	orders map[int]*domain.Order
	nextID int
}

// NewOrderRepository creates a new in-memory order repository.
func NewOrderRepository() *OrderRepository {
	return &OrderRepository{
		orders: make(map[int]*domain.Order),
		nextID: 1,
	}
}

// Save stores order in memory.
//
// If order does not yet have an ID, repository assigns a new one.
func (r *OrderRepository) Save(ctx context.Context, order *domain.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if order.ID() == 0 {
		order.SetID(r.nextID)
		r.nextID++
	}

	r.orders[order.ID()] = order
	return nil
}

// ByID returns order by its identifier.
func (r *OrderRepository) ByID(ctx context.Context, id int) (*domain.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	order, ok := r.orders[id]
	if !ok {
		return nil, domain.ErrOrderNotFound
	}

	return order, nil
}
