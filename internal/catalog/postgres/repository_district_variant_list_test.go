package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRepository_ListDistrictVariants_FiltersByDistrictAndProduct(t *testing.T) {
	pool := openTestPool(t)
	repo := NewRepository(pool)

	districtID := insertTestDistrict(t, pool)

	categoryID := insertTestCategory(t, pool)
	productMatchID := insertTestProductForCategory(t, pool, categoryID)

	variantPlacedID := insertTestVariantForProduct(t, pool, productMatchID)
	insertTestDistrictVariant(t, pool, districtID, variantPlacedID, 5900)

	_ = insertTestVariantForProduct(t, pool, productMatchID)

	productOtherID := insertTestProductForCategory(t, pool, categoryID)
	variantOtherID := insertTestVariantForProduct(t, pool, productOtherID)
	insertTestDistrictVariant(t, pool, districtID, variantOtherID, 6100)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	items, err := repo.ListDistrictVariants(ctx, districtID, productMatchID)
	require.NoError(t, err)
	require.Len(t, items, 1)
	require.Equal(t, variantPlacedID, items[0].ID)
}
