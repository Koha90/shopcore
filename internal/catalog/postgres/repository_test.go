package postgres

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"

	catalogservice "github.com/koha90/shopcore/internal/catalog/service"
)

func openTestPool(t *testing.T) *pgxpool.Pool {
	t.Helper()

	dsn := strings.TrimSpace(os.Getenv("TEST_DATABASE_URL"))
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	require.NoError(t, err)

	t.Cleanup(func() {
		pool.Close()
	})

	return pool
}

func testCode(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano())
}

func deleteCategoryByCode(t *testing.T, pool *pgxpool.Pool, code string) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := pool.Exec(ctx, `delete from catalog_categories where code = $1`, code)
	require.NoError(t, err)
}

func TestNewRepository(t *testing.T) {
	t.Parallel()

	repo := NewRepository(nil)
	require.NotNil(t, repo)
}

func TestRepository_CreateCategory_NilRepository(t *testing.T) {
	t.Parallel()

	var repo *Repository

	err := repo.CreateCategory(context.Background(), catalogservice.CreateCategoryParams{
		Code: "flowers",
		Name: "Цветы",
	})
	require.EqualError(t, err, "catalog postgres repository create category: repository is nil")
}

func TestRepository_CreateCategory_NilPool(t *testing.T) {
	t.Parallel()

	repo := &Repository{}

	err := repo.CreateCategory(context.Background(), catalogservice.CreateCategoryParams{
		Code: "flowers",
		Name: "Цветы",
	})
	require.EqualError(t, err, "catalog postgres repository create category: pool is nil")
}

func TestRepository_CreateCategory_Insert(t *testing.T) {
	pool := openTestPool(t)
	repo := NewRepository(pool)

	code := testCode("flowers")
	name := testCode("Цветы")
	nameLatin := testCode("flowers")
	t.Cleanup(func() {
		deleteCategoryByCode(t, pool, code)
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := repo.CreateCategory(ctx, catalogservice.CreateCategoryParams{
		Code:        code,
		Name:        name,
		NameLatin:   nameLatin,
		Description: "Букеты и композиции",
		SortOrder:   10,
	})
	require.NoError(t, err)

	var (
		gotName        string
		gotNameLatin   string
		gotDescription string
		isActive       bool
		sortOrder      int
	)

	err = pool.QueryRow(ctx, `
		select
			name,
			name_latin,
			description,
			is_active,
			sort_order
		from catalog_categories
		where code = $1
	`, code).Scan(
		&gotName,
		&gotNameLatin,
		&gotDescription,
		&isActive,
		&sortOrder,
	)
	require.NoError(t, err)
	require.Equal(t, name, gotName)
	require.Equal(t, nameLatin, gotNameLatin)
	require.Equal(t, "Букеты и композиции", gotDescription)
	require.True(t, isActive)
	require.Equal(t, 10, sortOrder)
}

func TestRepository_CreateCategory_DuplicateCode(t *testing.T) {
	pool := openTestPool(t)
	repo := NewRepository(pool)

	code := testCode("gifts")
	t.Cleanup(func() {
		deleteCategoryByCode(t, pool, code)
	})

	name := testCode("Подарки")
	nameLatin := testCode("Gifts")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := repo.CreateCategory(ctx, catalogservice.CreateCategoryParams{
		Code:        code,
		Name:        name,
		NameLatin:   nameLatin,
		Description: "Подарочные наборы",
		SortOrder:   20,
	})
	require.NoError(t, err)

	err = repo.CreateCategory(ctx, catalogservice.CreateCategoryParams{
		Code:        code,
		Name:        "Подарки 2",
		NameLatin:   "Gifts 2",
		Description: "Дубликат",
		SortOrder:   30,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "create category")
}
