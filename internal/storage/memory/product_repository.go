package memory

import (
	"context"
	"sync"

	"github.com/koha90/shopcore/internal/domain"
)

// ProductRepository stores products in process memory.
//
// It is intended of local development and tests.
// Repository assings incremental IDs to new products on first save.
type ProductRepository struct {
	mu       *sync.Mutex
	products map[int]*domain.Product
	nextID   int
}

// NewProductRepository creates a new in-memory product repository.
func NewProductRepository(mu *sync.Mutex) *ProductRepository {
	return &ProductRepository{
		mu:       mu,
		products: make(map[int]*domain.Product),
		nextID:   1,
	}
}

// Save stores product in memory.
//
// If product does not yet have an ID, repository assigns a new one.
func (r *ProductRepository) Save(ctx context.Context, product *domain.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if product.ID() == 0 {
		product.SetID(r.nextID)
		r.nextID++
	}

	r.products[product.ID()] = product

	return nil
}

// ByID returns product by its identifier.
func (r *ProductRepository) ByID(ctx context.Context, id int) (*domain.Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	product, ok := r.products[id]
	if !ok {
		return nil, domain.ErrProductNotFound
	}

	return product, nil
}
