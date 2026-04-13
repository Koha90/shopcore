package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRepository_ListDistrictProducts_FiltersByDistrictAndCategory(t *testing.T) {
	pool := openTestPool(t)
	repo := NewRepository(pool)

	districtID := insertTestDistrict(t, pool)

	categoryMatchID := insertTestCategory(t, pool)
	productPlacedID := insertTestProductForCategory(t, pool, categoryMatchID)
	variantPlacedID := insertTestVariantForProduct(t, pool, productPlacedID)
	insertTestDistrictVariant(t, pool, districtID, variantPlacedID, 5900)

	productHiddenID := insertTestProductForCategory(t, pool, categoryMatchID)
	_ = insertTestVariantForProduct(t, pool, productHiddenID)

	categoryOtherID := insertTestCategory(t, pool)
	productOtherID := insertTestProductForCategory(t, pool, categoryOtherID)
	variantOtherID := insertTestVariantForProduct(t, pool, productOtherID)
	insertTestDistrictVariant(t, pool, districtID, variantOtherID, 6100)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	items, err := repo.ListDistrictProducts(ctx, districtID, categoryMatchID)
	require.NoError(t, err)
	require.Len(t, items, 1)
	require.Equal(t, productPlacedID, items[0].ID)
}
