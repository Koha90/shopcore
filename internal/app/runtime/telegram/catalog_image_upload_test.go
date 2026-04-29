package telegram

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/koha90/shopcore/internal/flow"
)

func TestCatalogImageUploadPathProduct(t *testing.T) {
	t.Parallel()

	got, err := catalogImageUploadPath(
		flow.CatalogImageInputTarget{
			Kind:       flow.CatalogImageTargetProduct,
			EntityID:   123,
			EntityCode: "rose-box",
		},
		time.Unix(1710000000, 0),
	)

	require.NoError(t, err)
	assert.Equal(t, "assets/catalog/products/123-rose-box-1710000000.jpg", got)
}

func TestCatalogImageUploadPathVariant(t *testing.T) {
	t.Parallel()

	got, err := catalogImageUploadPath(
		flow.CatalogImageInputTarget{
			Kind:       flow.CatalogImageTargetVariant,
			EntityID:   456,
			EntityCode: "large-red",
		},
		time.Unix(1710000000, 0),
	)

	require.NoError(t, err)
	assert.Equal(t, "assets/catalog/variants/456-large-red-1710000000.jpg", got)
}

func TestCatalogImageUploadPathRejectsInvalidEntityID(t *testing.T) {
	t.Parallel()

	_, err := catalogImageUploadPath(
		flow.CatalogImageInputTarget{
			Kind:     flow.CatalogImageTargetProduct,
			EntityID: 0,
		},
		time.Unix(1710000000, 0),
	)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "entity id is invalid")
}

func TestCatalogImageUploadPathRejectsUnknownKind(t *testing.T) {
	t.Parallel()

	_, err := catalogImageUploadPath(
		flow.CatalogImageInputTarget{
			Kind:       flow.CatalogImageTargetKind("unknown"),
			EntityID:   123,
			EntityCode: "rose-box",
		},
		time.Unix(1710000000, 0),
	)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown catalog image target kind")
}

func TestCatalogImageUploadPathFallsBackToIDWithoutCode(t *testing.T) {
	t.Parallel()

	got, err := catalogImageUploadPath(
		flow.CatalogImageInputTarget{
			Kind:     flow.CatalogImageTargetProduct,
			EntityID: 123,
		},
		time.Unix(1710000000, 0),
	)

	require.NoError(t, err)
	assert.Equal(t, "assets/catalog/products/123-1710000000.jpg", got)
}
