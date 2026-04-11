package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"

	"github.com/koha90/shopcore/internal/flow"
)

func deleteDistrictByCode(t *testing.T, pool *pgxpool.Pool, cityID int, code string) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := pool.Exec(ctx, `
		delete from catalog_districts
		where city_id = $1 and code = $2
	`, cityID, code)
	require.NoError(t, err)
}

func insertTestDistrictWithOrder(t *testing.T, pool *pgxpool.Pool, sortOrder int, isActive bool) (int, string, string) {
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
		values ($1, $2, $3, $4, $5, $6, now(), now())
		returning id
	`, cityID, code, name, nameLatin, isActive, sortOrder).Scan(&id)
	require.NoError(t, err)

	t.Cleanup(func() {
		deleteDistrictByCode(t, pool, cityID, code)
	})

	return id, code, name
}

func TestRepository_ListDistricts_NilRepository(t *testing.T) {
	t.Parallel()

	var repo *Repository

	items, err := repo.ListDistricts(context.Background())
	require.Nil(t, items)
	require.EqualError(t, err, "catalog postgres repository list districts: repository is nil")
}

func TestRepository_ListDistricts_NilPool(t *testing.T) {
	t.Parallel()

	repo := &Repository{}

	items, err := repo.ListDistricts(context.Background())
	require.Nil(t, items)
	require.EqualError(t, err, "catalog postgres repository list districts: pool is nil")
}

func TestRepository_ListDistricts_ReturnsOnlyActiveDistrictsOrdered(t *testing.T) {
	pool := openTestPool(t)
	repo := NewRepository(pool)

	id3, code3, name3 := insertTestDistrictWithOrder(t, pool, 30, true)
	_, _, _ = insertTestDistrictWithOrder(t, pool, 5, false)
	id1, code1, name1 := insertTestDistrictWithOrder(t, pool, 10, true)
	id2, code2, name2 := insertTestDistrictWithOrder(t, pool, 20, true)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	items, err := repo.ListDistricts(ctx)
	require.NoError(t, err)

	var got []flow.DistrictListItem
	for _, item := range items {
		switch item.ID {
		case id1, id2, id3:
			got = append(got, item)
		}
	}

	require.Len(t, got, 3)

	require.Equal(t, flow.DistrictListItem{
		ID:    id1,
		Code:  code1,
		Label: name1,
	}, got[0])

	require.Equal(t, flow.DistrictListItem{
		ID:    id2,
		Code:  code2,
		Label: name2,
	}, got[1])

	require.Equal(t, flow.DistrictListItem{
		ID:    id3,
		Code:  code3,
		Label: name3,
	}, got[2])
}
