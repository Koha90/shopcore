package service

import (
	"context"

	"github.com/koha90/shopcore/internal/domain"
)

// ProductRepository defines persistence operations required by ProductService.
type ProductRepository interface {
	Save(ctx context.Context, p *domain.Product) error
	ByID(ctx context.Context, id int) (*domain.Product, error)
}

// ProductReader defines read operations required by order use cases.
type ProductReader interface {
	ByID(ctx context.Context, id int) (*domain.Product, error)
}

// OrderRepository defines persistence contain for Order aggregate.
type OrderRepository interface {
	Save(ctx context.Context, order *domain.Order) error
	ByID(ctx context.Context, id int) (*domain.Order, error)
}

// UserRepository defines persistence operations
// required by UserService.
type UserRepository interface {
	Save(ctx context.Context, u *domain.User) error
	ByID(ctx context.Context, id int) (*domain.User, error)
}
