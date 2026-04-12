package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	catalogservice "github.com/koha90/shopcore/internal/catalog/service"
)

func TestRepository_UpdateDistrictVariantPrice_NilRepository(t *testing.T) {
	t.Parallel()

	var repo *Repository

	err := repo.UpdateDistrictVariantPrice(context.Background(), catalogservice.UpdateDistrictVariantPriceParams{
		DistrictID: 1,
		VariantID:  2,
		Price:      6100,
	})
	require.EqualError(t, err, "catalog postgres repository update district variant price: repository is nil")
}

func TestRepository_UpdateDistrictVariantPrice_NilPool(t *testing.T) {
	t.Parallel()

	repo := &Repository{}

	err := repo.UpdateDistrictVariantPrice(context.Background(), catalogservice.UpdateDistrictVariantPriceParams{
		DistrictID: 1,
		VariantID:  2,
		Price:      6100,
	})
	require.EqualError(t, err, "catalog postgres repository update district variant price: pool is nil")
}

func TestRepository_UpdateDistrictVariantPrice_Update(t *testing.T) {
	pool := openTestPool(t)
	repo := NewRepository(pool)

	districtID := insertTestDistrict(t, pool)
	variantID := insertTestVariant(t, pool)

	t.Cleanup(func() {
		deleteDistrictVariant(t, pool, districtID, variantID)
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := repo.CreateDistrictVariant(ctx, catalogservice.CreateDistrictVariantParams{
		DistrictID: districtID,
		VariantID:  variantID,
		Price:      5900,
	})
	require.NoError(t, err)

	err = repo.UpdateDistrictVariantPrice(ctx, catalogservice.UpdateDistrictVariantPriceParams{
		DistrictID: districtID,
		VariantID:  variantID,
		Price:      6100,
	})
	require.NoError(t, err)

	var gotPrice int
	err = pool.QueryRow(ctx, `
		select price
		from catalog_district_variants
		where district_id = $1 and variant_id = $2
	`, districtID, variantID).Scan(&gotPrice)
	require.NoError(t, err)
	require.Equal(t, 6100, gotPrice)
}

func TestRepository_UpdateDistrictVariantPrice_NotFound(t *testing.T) {
	pool := openTestPool(t)
	repo := NewRepository(pool)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := repo.UpdateDistrictVariantPrice(ctx, catalogservice.UpdateDistrictVariantPriceParams{
		DistrictID: 999999,
		VariantID:  999999,
		Price:      6100,
	})
	require.ErrorIs(t, err, ErrDistrictVariantNotFound)
}
