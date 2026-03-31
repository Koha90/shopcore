package flow

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStaticCatalogProvider_Catalog(t *testing.T) {
	catalog := DemoCatalog()
	provider := NewStaticCatalogProvider(catalog)

	got, err := provider.Catalog(context.Background())

	require.NoError(t, err)
	require.Equal(t, catalog, got)
}

func TestNewServiceWithCatalogProvider_Defaults(t *testing.T) {
	svc := NewServiceWithCatalogProvider(nil, nil)

	require.NotNil(t, svc)
	require.NotNil(t, svc.store)
	require.NotNil(t, svc.provider)
}

func TestNewServiceWithCatalogProvider_UsesProvidedCatalog(t *testing.T) {
	custom := Catalog{
		Schema: CatalogSchema{
			Levels: []CatalogLevel{LevelCity},
		},
		Roots: []CatalogNode{
			{
				Level: LevelCity,
				ID:    "test-city",
				Label: "Тестовый город",
			},
		},
	}

	svc := NewServiceWithCatalogProvider(
		nil,
		NewStaticCatalogProvider(custom),
	)

	got, err := svc.provider.Catalog(context.Background())
	require.NoError(t, err)
	require.Equal(t, custom, got)
}
