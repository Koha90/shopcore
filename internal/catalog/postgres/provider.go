package postgres

import (
	"context"

	"github.com/koha90/shopcore/internal/flow"
)

// CatalogLoader loads catalog snapshot from storage.
type CatalogLoader interface {
	LoadCatalog(ctx context.Context) (flow.Catalog, error)
}

// CatalogProvider adapts CatalogLoader to flow.CatalogProvider
type CatalogProvider struct {
	loader CatalogLoader
}

// NewCatalogProvider constructs Postgres-backed flow catalog provider.
func NewCatalogProvider(loader CatalogLoader) *CatalogProvider {
	return &CatalogProvider{
		loader: loader,
	}
}

// Catalog returns flow-ready catalog snapshot loaded from Postgres.
func (p *CatalogProvider) Catalog(ctx context.Context) (flow.Catalog, error) {
	return p.loader.LoadCatalog(ctx)
}
