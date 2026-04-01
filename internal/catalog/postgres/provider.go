// Package postgres ...
package postgres

import (
	"context"

	"github.com/koha90/shopcore/internal/flow"
)

type CatalogLoader interface {
	LoadCatalog(ctx context.Context) (flow.Catalog, error)
}

type CatalogProvider struct {
	loader CatalogLoader
}

func NewCatalogProvider(loader CatalogLoader) *CatalogProvider {
	return &CatalogProvider{
		loader: loader,
	}
}

func (p *CatalogProvider) Catalog(ctx context.Context) (flow.Catalog, error) {
	return p.loader.LoadCatalog(ctx)
}
