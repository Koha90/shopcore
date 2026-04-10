package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/koha90/shopcore/internal/flow"
)

// Loader reads catalog data from Postgres and builds flow.Catalog.
type Loader struct {
	pool *pgxpool.Pool
}

// NewLoader constructs catalog loader backed by pgx pool.
func NewLoader(pool *pgxpool.Pool) *Loader {
	return &Loader{
		pool: pool,
	}
}

// LoadCatalog loads relational catalog rows and maps them into flow.Catalog.
func (l *Loader) LoadCatalog(ctx context.Context) (flow.Catalog, error) {
	cities, err := l.loadCities(ctx)
	if err != nil {
		return flow.Catalog{}, fmt.Errorf("load cities: %w", err)
	}

	categories, err := l.loadCategories(ctx)
	if err != nil {
		return flow.Catalog{}, fmt.Errorf("load categories: %w", err)
	}

	districts, err := l.loadDistricts(ctx)
	if err != nil {
		return flow.Catalog{}, fmt.Errorf("load districts: %w", err)
	}

	products, err := l.loadProducts(ctx)
	if err != nil {
		return flow.Catalog{}, fmt.Errorf("load products: %w", err)
	}

	variants, err := l.loadVariants(ctx)
	if err != nil {
		return flow.Catalog{}, fmt.Errorf("load variants: %w", err)
	}

	districtVariants, err := l.loadDistrictVariants(ctx)
	if err != nil {
		return flow.Catalog{}, fmt.Errorf("load district variants: %w", err)
	}

	catalog := buildCatalog(
		cities,
		categories,
		districts,
		products,
		variants,
		districtVariants,
	)

	return catalog, nil
}
