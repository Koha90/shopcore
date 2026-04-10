package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"

	catalogservice "github.com/koha90/shopcore/internal/catalog/service"
)

func deleteDistrictVariant(t *testing.T, pool *pgxpool.Pool, districtID, variantID int) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := pool.Exec(ctx, `
		delete from catalog_district_variants
		where district_id = $1 and variant_id = $2
	`, districtID, variantID)
	require.NoError(t, err)
}

func insertTestDistrict(t *testing.T, pool *pgxpool.Pool) int {
	t.Helper()

	cityID := insertTestCity(t, pool)
	code := testCode("district")
	name := testCode("Район")
	nameLatin := testCode("district")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var id int
	err := pool.QueryRow(ctx, `
		insert into catalog_districts (
			city_id,
			code,
			name,
			name_latin,
			is_active,
			sort_order,
			created_at,
			updated_at
		)
		values ($1, $2, $3, $4, true, 0, now(), now())
		returning id
	`, cityID, code, name, nameLatin).Scan(&id)
	require.NoError(t, err)

	return id
}

func insertTestVariant(t *testing.T, pool *pgxpool.Pool) int {
	t.Helper()

	productID := insertTestProduct(t, pool)
	code := testCode("variant")
	name := testCode("Вариант")
	nameLatin := testCode("variant")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var id int
	err := pool.QueryRow(ctx, `
		insert into catalog_variants (
			product_id,
			code,
			name,
			name_latin,
			description,
			is_active,
			sort_order,
			created_at,
			updated_at
		)
		values ($1, $2, $3, $4, '', true, 0, now(), now())
		returning id
	`, productID, code, name, nameLatin).Scan(&id)
	require.NoError(t, err)

	return id
}

func TestRepository_CreateDistrictVariant_NilRepository(t *testing.T) {
	t.Parallel()

	var repo *Repository

	err := repo.CreateDistrictVariant(context.Background(), catalogservice.CreateDistrictVariantParams{
		DistrictID: 1,
		VariantID:  2,
		Price:      5900,
	})
	require.EqualError(t, err, "catalog postgres repository create district variant: repository is nil")
}

func TestRepository_CreateDistrictVariant_NilPool(t *testing.T) {
	t.Parallel()

	repo := &Repository{}

	err := repo.CreateDistrictVariant(context.Background(), catalogservice.CreateDistrictVariantParams{
		DistrictID: 1,
		VariantID:  2,
		Price:      5900,
	})
	require.EqualError(t, err, "catalog postgres repository create district variant: pool is nil")
}

func TestRepository_CreateDistrictVariant_Insert(t *testing.T) {
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

	var (
		gotDistrictID int
		gotVariantID  int
		gotPrice      int
		isActive      bool
	)

	err = pool.QueryRow(ctx, `
		select
			district_id,
			variant_id,
			price,
			is_active
		from catalog_district_variants
		where district_id = $1 and variant_id = $2
	`, districtID, variantID).Scan(
		&gotDistrictID,
		&gotVariantID,
		&gotPrice,
		&isActive,
	)
	require.NoError(t, err)
	require.Equal(t, districtID, gotDistrictID)
	require.Equal(t, variantID, gotVariantID)
	require.Equal(t, 5900, gotPrice)
	require.True(t, isActive)
}

func TestRepository_CreateDistrictVariant_DuplicatePlacement(t *testing.T) {
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

	err = repo.CreateDistrictVariant(ctx, catalogservice.CreateDistrictVariantParams{
		DistrictID: districtID,
		VariantID:  variantID,
		Price:      6100,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "create district variant")
}
