package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"

	catalogservice "github.com/koha90/shopcore/internal/catalog/service"
)

func deleteVariantByCode(t *testing.T, pool *pgxpool.Pool, productID int, code string) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := pool.Exec(ctx, `
		delete from catalog_variants
		where product_id = $1 and code = $2
	`, productID, code)
	require.NoError(t, err)
}

func insertTestProduct(t *testing.T, pool *pgxpool.Pool) int {
	t.Helper()

	categoryID := insertTestCategory(t, pool)
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

	t.Cleanup(func() {
		deleteProductByCode(t, pool, categoryID, code)
	})

	return id
}

func TestRepository_CreateVariant_NilRepository(t *testing.T) {
	t.Parallel()

	var repo *Repository

	err := repo.CreateVariant(context.Background(), catalogservice.CreateVariantParams{
		ProductID: 1,
		Code:      "l-25",
		Name:      "L / 25 шт",
	})
	require.EqualError(t, err, "catalog postgres repository create variant: repository is nil")
}

func TestRepository_CreateVariant_NilPool(t *testing.T) {
	t.Parallel()

	repo := &Repository{}

	err := repo.CreateVariant(context.Background(), catalogservice.CreateVariantParams{
		ProductID: 1,
		Code:      "l-25",
		Name:      "L / 25 шт",
	})
	require.EqualError(t, err, "catalog postgres repository create variant: pool is nil")
}

func TestRepository_CreateVariant_Insert(t *testing.T) {
	pool := openTestPool(t)
	repo := NewRepository(pool)

	productID := insertTestProduct(t, pool)
	code := testCode("l-25")
	name := testCode("L / 25 шт")
	nameLatin := testCode("l-25")

	t.Cleanup(func() {
		deleteVariantByCode(t, pool, productID, code)
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := repo.CreateVariant(ctx, catalogservice.CreateVariantParams{
		ProductID:   productID,
		Code:        code,
		Name:        name,
		NameLatin:   nameLatin,
		Description: "Большая упаковка",
		SortOrder:   10,
	})
	require.NoError(t, err)

	var (
		gotProductID   int
		gotName        string
		gotNameLatin   string
		gotDescription string
		isActive       bool
		sortOrder      int
	)

	err = pool.QueryRow(ctx, `
		select
			product_id,
			name,
			name_latin,
			description,
			is_active,
			sort_order
		from catalog_variants
		where product_id = $1 and code = $2
	`, productID, code).Scan(
		&gotProductID,
		&gotName,
		&gotNameLatin,
		&gotDescription,
		&isActive,
		&sortOrder,
	)
	require.NoError(t, err)
	require.Equal(t, productID, gotProductID)
	require.Equal(t, name, gotName)
	require.Equal(t, nameLatin, gotNameLatin)
	require.Equal(t, "Большая упаковка", gotDescription)
	require.True(t, isActive)
	require.Equal(t, 10, sortOrder)
}

func TestRepository_CreateVariant_DuplicateCodeInSameProduct(t *testing.T) {
	pool := openTestPool(t)
	repo := NewRepository(pool)

	productID := insertTestProduct(t, pool)
	code := testCode("l-25")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := repo.CreateVariant(ctx, catalogservice.CreateVariantParams{
		ProductID:   productID,
		Code:        code,
		Name:        testCode("L / 25 шт"),
		NameLatin:   testCode("l-25"),
		Description: "Первый",
		SortOrder:   10,
	})
	require.NoError(t, err)

	err = repo.CreateVariant(ctx, catalogservice.CreateVariantParams{
		ProductID:   productID,
		Code:        code,
		Name:        testCode("L / 25 шт 2"),
		NameLatin:   testCode("l-25-2"),
		Description: "Второй",
		SortOrder:   20,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "create variant")
}

func TestRepository_CreateVariant_DuplicateNameInSameProduct(t *testing.T) {
	pool := openTestPool(t)
	repo := NewRepository(pool)

	productID := insertTestProduct(t, pool)
	name := testCode("L / 25 шт")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := repo.CreateVariant(ctx, catalogservice.CreateVariantParams{
		ProductID:   productID,
		Code:        testCode("l-25"),
		Name:        name,
		NameLatin:   testCode("l-25"),
		Description: "Первый",
		SortOrder:   10,
	})
	require.NoError(t, err)

	err = repo.CreateVariant(ctx, catalogservice.CreateVariantParams{
		ProductID:   productID,
		Code:        testCode("l-25-2"),
		Name:        name,
		NameLatin:   testCode("l-25-2"),
		Description: "Второй",
		SortOrder:   20,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "create variant")
}

func TestRepository_CreateVariant_SameCodeInDifferentProducts_IsAllowed(t *testing.T) {
	pool := openTestPool(t)
	repo := NewRepository(pool)

	productID1 := insertTestProduct(t, pool)
	productID2 := insertTestProduct(t, pool)
	code := testCode("l-25")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := repo.CreateVariant(ctx, catalogservice.CreateVariantParams{
		ProductID:   productID1,
		Code:        code,
		Name:        testCode("L / 25 шт 1"),
		NameLatin:   testCode("l-25-1"),
		Description: "Первый",
		SortOrder:   10,
	})
	require.NoError(t, err)

	err = repo.CreateVariant(ctx, catalogservice.CreateVariantParams{
		ProductID:   productID2,
		Code:        code,
		Name:        testCode("L / 25 шт 2"),
		NameLatin:   testCode("l-25-2"),
		Description: "Второй",
		SortOrder:   20,
	})
	require.NoError(t, err)
}
