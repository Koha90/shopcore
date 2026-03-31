package flow

// CatalogProvider returns catalog snapshot used by flow navigation.
//
// Current implementation is static in-memory.
// Later it can be replaced with config-backed or storage-backed provider.
type CatalogProvider interface {
	Catalog() Catalog
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
func (p *StaticCatalogProvider) Catalog() Catalog {
	return p.catalog
}
