package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"

	catalogservice "github.com/koha90/shopcore/internal/catalog/service"
)

func insertTestCity(t *testing.T, pool *pgxpool.Pool) int {
	t.Helper()

	code := testCode("city")
	name := testCode("Город")
	nameLatin := testCode("city")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var id int
	err := pool.QueryRow(ctx, `
		insert into cities (
			code,
			name,
			name_latin,
			is_active,
			sort_order,
			created_at,
			updated_at
		)
		values ($1, $2, $3, true, 0, now(), now())
		returning id
	`, code, name, nameLatin).Scan(&id)
	require.NoError(t, err)

	t.Cleanup(func() {
		deleteCityByCode(t, pool, code)
	})

	return id
}

func TestRepository_CreateDistrict_NilRepository(t *testing.T) {
	t.Parallel()

	var repo *Repository

	err := repo.CreateDistrict(context.Background(), catalogservice.CreateDistrictParams{
		CityID: 1,
		Code:   "center",
		Name:   "Центр",
	})
	require.EqualError(t, err, "catalog postgres repository create district: repository is nil")
}

func TestRepository_CreateDistrict_NilPool(t *testing.T) {
	t.Parallel()

	repo := &Repository{}

	err := repo.CreateDistrict(context.Background(), catalogservice.CreateDistrictParams{
		CityID: 1,
		Code:   "center",
		Name:   "Центр",
	})
	require.EqualError(t, err, "catalog postgres repository create district: pool is nil")
}

func TestRepository_CreateDistrict_Insert(t *testing.T) {
	pool := openTestPool(t)
	repo := NewRepository(pool)

	cityID := insertTestCity(t, pool)

	code := testCode("center")
	name := testCode("Центр")
	nameLatin := testCode("center")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := repo.CreateDistrict(ctx, catalogservice.CreateDistrictParams{
		CityID:    cityID,
		Code:      code,
		Name:      name,
		NameLatin: nameLatin,
		SortOrder: 10,
	})
	require.NoError(t, err)

	var (
		gotCityID    int
		gotName      string
		gotNameLatin string
		isActive     bool
		sortOrder    int
	)

	err = pool.QueryRow(ctx, `
		select
			city_id,
			name,
			name_latin,
			is_active,
			sort_order
		from catalog_districts
		where city_id = $1 and code = $2
	`, cityID, code).Scan(
		&gotCityID,
		&gotName,
		&gotNameLatin,
		&isActive,
		&sortOrder,
	)
	require.NoError(t, err)
	require.Equal(t, cityID, gotCityID)
	require.Equal(t, name, gotName)
	require.Equal(t, nameLatin, gotNameLatin)
	require.True(t, isActive)
	require.Equal(t, 10, sortOrder)
}

func TestRepository_CreateDistrict_DuplicateCodeInSameCity(t *testing.T) {
	pool := openTestPool(t)
	repo := NewRepository(pool)

	cityID := insertTestCity(t, pool)
	code := testCode("center")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := repo.CreateDistrict(ctx, catalogservice.CreateDistrictParams{
		CityID:    cityID,
		Code:      code,
		Name:      testCode("Центр"),
		NameLatin: testCode("center"),
		SortOrder: 10,
	})
	require.NoError(t, err)

	err = repo.CreateDistrict(ctx, catalogservice.CreateDistrictParams{
		CityID:    cityID,
		Code:      code,
		Name:      testCode("Центр-2"),
		NameLatin: testCode("center-2"),
		SortOrder: 20,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "create district")
}

func TestRepository_CreateDistrict_DuplicateNameInSameCity(t *testing.T) {
	pool := openTestPool(t)
	repo := NewRepository(pool)

	cityID := insertTestCity(t, pool)
	name := testCode("Центр")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := repo.CreateDistrict(ctx, catalogservice.CreateDistrictParams{
		CityID:    cityID,
		Code:      testCode("center"),
		Name:      name,
		NameLatin: testCode("center"),
		SortOrder: 10,
	})
	require.NoError(t, err)

	err = repo.CreateDistrict(ctx, catalogservice.CreateDistrictParams{
		CityID:    cityID,
		Code:      testCode("center-2"),
		Name:      name,
		NameLatin: testCode("center-2"),
		SortOrder: 20,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "create district")
}

func TestRepository_CreateDistrict_SameCodeInDifferentCities_IsAllowed(t *testing.T) {
	pool := openTestPool(t)
	repo := NewRepository(pool)

	cityID1 := insertTestCity(t, pool)
	cityID2 := insertTestCity(t, pool)
	code := testCode("center")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := repo.CreateDistrict(ctx, catalogservice.CreateDistrictParams{
		CityID:    cityID1,
		Code:      code,
		Name:      testCode("Центр-1"),
		NameLatin: testCode("center-1"),
		SortOrder: 10,
	})
	require.NoError(t, err)

	err = repo.CreateDistrict(ctx, catalogservice.CreateDistrictParams{
		CityID:    cityID2,
		Code:      code,
		Name:      testCode("Центр-2"),
		NameLatin: testCode("center-2"),
		SortOrder: 20,
	})
	require.NoError(t, err)
}
