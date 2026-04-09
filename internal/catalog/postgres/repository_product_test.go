package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"

	catalogservice "github.com/koha90/shopcore/internal/catalog/service"
)

func deleteProductByCode(t *testing.T, pool *pgxpool.Pool, categoryID int, code string) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := pool.Exec(ctx, `
		delete from catalog_products
		where category_id = $1 and code = $2
	`, categoryID, code)
	require.NoError(t, err)
}

func insertTestCategory(t *testing.T, pool *pgxpool.Pool) int {
	t.Helper()

	code := testCode("category")
	name := testCode("Категория")
	nameLatin := testCode("category")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var id int
	err := pool.QueryRow(ctx, `
		insert into catalog_categories (
			code,
			name,
			name_latin,
			description,
			is_active,
			sort_order,
			created_at,
			updated_at
		)
		values ($1, $2, $3, '', true, 0, now(), now())
		returning id
	`, code, name, nameLatin).Scan(&id)
	require.NoError(t, err)

	t.Cleanup(func() {
		deleteCategoryByCode(t, pool, code)
	})

	return id
}

func TestRepository_CreateProduct_NilRepository(t *testing.T) {
	t.Parallel()

	var repo *Repository

	err := repo.CreateProduct(context.Background(), catalogservice.CreateProductParams{
		CategoryID: 1,
		Code:       "rose-box",
		Name:       "Розы в коробке",
	})
	require.EqualError(t, err, "catalog postgres repository create product: repository is nil")
}

func TestRepository_CreateProduct_NilPool(t *testing.T) {
	t.Parallel()

	repo := &Repository{}

	err := repo.CreateProduct(context.Background(), catalogservice.CreateProductParams{
		CategoryID: 1,
		Code:       "rose-box",
		Name:       "Розы в коробке",
	})
	require.EqualError(t, err, "catalog postgres repository create product: pool is nil")
}

func TestRepository_CreateProduct_Insert(t *testing.T) {
	pool := openTestPool(t)
	repo := NewRepository(pool)

	categoryID := insertTestCategory(t, pool)
	code := testCode("rose-box")
	name := testCode("Розы в коробке")
	nameLatin := testCode("rose-box")

	t.Cleanup(func() {
		deleteProductByCode(t, pool, categoryID, code)
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := repo.CreateProduct(ctx, catalogservice.CreateProductParams{
		CategoryID:  categoryID,
		Code:        code,
		Name:        name,
		NameLatin:   nameLatin,
		Description: "Композиция из роз",
		SortOrder:   10,
	})
	require.NoError(t, err)

	var (
		gotCategoryID  int
		gotName        string
		gotNameLatin   string
		gotDescription string
		isActive       bool
		sortOrder      int
	)

	err = pool.QueryRow(ctx, `
		select
			category_id,
			name,
			name_latin,
			description,
			is_active,
			sort_order
		from catalog_products
		where category_id = $1 and code = $2
	`, categoryID, code).Scan(
		&gotCategoryID,
		&gotName,
		&gotNameLatin,
		&gotDescription,
		&isActive,
		&sortOrder,
	)
	require.NoError(t, err)
	require.Equal(t, categoryID, gotCategoryID)
	require.Equal(t, name, gotName)
	require.Equal(t, nameLatin, gotNameLatin)
	require.Equal(t, "Композиция из роз", gotDescription)
	require.True(t, isActive)
	require.Equal(t, 10, sortOrder)
}

func TestRepository_CreateProduct_DuplicateCodeInSameCategory(t *testing.T) {
	pool := openTestPool(t)
	repo := NewRepository(pool)

	categoryID := insertTestCategory(t, pool)
	code := testCode("rose-box")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := repo.CreateProduct(ctx, catalogservice.CreateProductParams{
		CategoryID:  categoryID,
		Code:        code,
		Name:        testCode("Розы в коробке"),
		NameLatin:   testCode("rose-box"),
		Description: "Первый",
		SortOrder:   10,
	})
	require.NoError(t, err)

	err = repo.CreateProduct(ctx, catalogservice.CreateProductParams{
		CategoryID:  categoryID,
		Code:        code,
		Name:        testCode("Розы в коробке 2"),
		NameLatin:   testCode("rose-box-2"),
		Description: "Второй",
		SortOrder:   20,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "create product")
}

func TestRepository_CreateProduct_DuplicateNameInSameCategory(t *testing.T) {
	pool := openTestPool(t)
	repo := NewRepository(pool)

	categoryID := insertTestCategory(t, pool)
	name := testCode("Розы в коробке")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := repo.CreateProduct(ctx, catalogservice.CreateProductParams{
		CategoryID:  categoryID,
		Code:        testCode("rose-box"),
		Name:        name,
		NameLatin:   testCode("rose-box"),
		Description: "Первый",
		SortOrder:   10,
	})
	require.NoError(t, err)

	err = repo.CreateProduct(ctx, catalogservice.CreateProductParams{
		CategoryID:  categoryID,
		Code:        testCode("rose-box-2"),
		Name:        name,
		NameLatin:   testCode("rose-box-2"),
		Description: "Второй",
		SortOrder:   20,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "create product")
}

func TestRepository_CreateProduct_SameCodeInDifferentCategories_IsAllowed(t *testing.T) {
	pool := openTestPool(t)
	repo := NewRepository(pool)

	categoryID1 := insertTestCategory(t, pool)
	categoryID2 := insertTestCategory(t, pool)
	code := testCode("rose-box")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := repo.CreateProduct(ctx, catalogservice.CreateProductParams{
		CategoryID:  categoryID1,
		Code:        code,
		Name:        testCode("Розы в коробке 1"),
		NameLatin:   testCode("rose-box-1"),
		Description: "Первый",
		SortOrder:   10,
	})
	require.NoError(t, err)

	err = repo.CreateProduct(ctx, catalogservice.CreateProductParams{
		CategoryID:  categoryID2,
		Code:        code,
		Name:        testCode("Розы в коробке 2"),
		NameLatin:   testCode("rose-box-2"),
		Description: "Второй",
		SortOrder:   20,
	})
	require.NoError(t, err)
}
