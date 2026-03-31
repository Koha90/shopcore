package flow

import "context"

// CatalogProvider returns catalog snapshot used by flow navigation.
//
// Flow does not depend on where catalog comes from.
// Current implementation is static in-memory.
// Future implementations may load catalog from config, storage, or adapters.
type CatalogProvider interface {
	Catalog(ctx context.Context) (Catalog, error)
}

// StaticCatalogProvider returns one fixed catalog snapshot.
type StaticCatalogProvider struct {
	catalog Catalog
}

// NewStaticCatalogProvider constructs static catalog provider.
func NewStaticCatalogProvider(catalog Catalog) *StaticCatalogProvider {
	return &StaticCatalogProvider{catalog: catalog}
}

// Catalog returns configured catalog snapshot.
func (p *StaticCatalogProvider) Catalog(ctx context.Context) (Catalog, error) {
	return p.catalog, nil
}
