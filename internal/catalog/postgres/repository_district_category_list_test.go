package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

func insertTestProductForCategory(t *testing.T, pool *pgxpool.Pool, categoryID int) int {
	t.Helper()

	code := testCode("product")
	name := testCode("Товар")
	nameLatin := testCode("product")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var id int
	err := pool.QueryRow(ctx, `
		insert into catalog_products (
			category_id,
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
	`, categoryID, code, name, nameLatin).Scan(&id)
	require.NoError(t, err)

	return id
}

func insertTestVariantForProduct(t *testing.T, pool *pgxpool.Pool, productID int) int {
	t.Helper()

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

func insertTestDistrictVariant(t *testing.T, pool *pgxpool.Pool, districtID, variantID, price int) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := pool.Exec(ctx, `
		insert into catalog_district_variants (
			district_id,
			variant_id,
			price,
			is_active,
			created_at,
			updated_at
		)
		values ($1, $2, $3, true, now(), now())
	`, districtID, variantID, price)
	require.NoError(t, err)

	t.Cleanup(func() {
		deleteDistrictVariant(t, pool, districtID, variantID)
	})
}

func TestRepository_ListDistrictCategories_FiltersByDistrictPlacement(t *testing.T) {
	pool := openTestPool(t)
	repo := NewRepository(pool)

	districtID := insertTestDistrict(t, pool)

	categoryPlacedID := insertTestCategory(t, pool)
	productPlacedID := insertTestProductForCategory(t, pool, categoryPlacedID)
	variantPlacedID := insertTestVariantForProduct(t, pool, productPlacedID)
	insertTestDistrictVariant(t, pool, districtID, variantPlacedID, 5900)

	categoryHiddenID := insertTestCategory(t, pool)
	productHiddenID := insertTestProductForCategory(t, pool, categoryHiddenID)
	_ = insertTestVariantForProduct(t, pool, productHiddenID)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	items, err := repo.ListDistrictCategories(ctx, districtID)
	require.NoError(t, err)
	require.Len(t, items, 1)
	require.Equal(t, categoryPlacedID, items[0].ID)
}
