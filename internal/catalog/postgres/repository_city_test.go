package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"

	catalogservice "github.com/koha90/shopcore/internal/catalog/service"
)

func deleteCityByCode(t *testing.T, pool *pgxpool.Pool, code string) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := pool.Exec(ctx, `delete from cities where code = $1`, code)
	require.NoError(t, err)
}

func TestRepository_CreateCity_NelRepository(t *testing.T) {
	t.Parallel()

	var repo *Repository

	err := repo.CreateCity(context.Background(), catalogservice.CreateCityParams{
		Code: "moscow",
		Name: "Москва",
	})
	require.EqualError(t, err, "catalog postgres repository create city: repository is nil")
}

func TestRepository_CreateCity_Insert(t *testing.T) {
	pool := openTestPool(t)
	repo := NewRepository(pool)

	code := testCode("moscow")
	name := testCode("Москва")
	nameLatin := testCode("Moscow")
	t.Cleanup(func() {
		deleteCategoryByCode(t, pool, code)
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := repo.CreateCity(ctx, catalogservice.CreateCityParams{
		Code:      code,
		Name:      name,
		NameLatin: nameLatin,
		SortOrder: 10,
	})
	require.NoError(t, err)

	var (
		gotName      string
		gotNameLatin string
		isActive     bool
		sortOrder    int
	)

	err = pool.QueryRow(ctx, `
  	select
			name,
			name_latin,
			is_active,
			sort_order
		from cities
		where code = $1
	`, code).Scan(
		&gotName,
		&gotNameLatin,
		&isActive,
		&sortOrder,
	)
	require.NoError(t, err)
	require.Equal(t, name, gotName)
	require.Equal(t, nameLatin, gotNameLatin)
	require.True(t, isActive)
	require.Equal(t, 10, sortOrder)
}

func TestRepository_CreateCity_DuplicateCode(t *testing.T) {
	pool := openTestPool(t)
	repo := NewRepository(pool)

	code := testCode("spb")
	t.Cleanup(func() {
		deleteCategoryByCode(t, pool, code)
	})

	name := testCode("Санкт-Петербург")
	nameLatin := testCode("Saint-Petersburg")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := repo.CreateCity(ctx, catalogservice.CreateCityParams{
		Code:      code,
		Name:      name,
		NameLatin: nameLatin,
		SortOrder: 20,
	})
	require.NoError(t, err)

	err = repo.CreateCity(ctx, catalogservice.CreateCityParams{
		Code:      code,
		Name:      "Санкт-Петербург 2",
		NameLatin: "Saint-Petersburg-2",
		SortOrder: 30,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "create city")
}
