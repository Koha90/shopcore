package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"

	"github.com/koha90/shopcore/internal/flow"
)

func insertTestVariantWithOrder(t *testing.T, pool *pgxpool.Pool, sortOrder int, isActive bool) (int, string, string) {
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
		values ($1, $2, $3, $4, '', $5, $6, now(), now())
		returning id
	`, productID, code, name, nameLatin, isActive, sortOrder).Scan(&id)
	require.NoError(t, err)

	t.Cleanup(func() {
		deleteVariantByCode(t, pool, productID, code)
	})

	return id, code, name
}

func TestRepository_ListVariants_NilRepository(t *testing.T) {
	t.Parallel()

	var repo *Repository

	items, err := repo.ListVariants(context.Background())
	require.Nil(t, items)
	require.EqualError(t, err, "catalog postgres repository list variants: repository is nil")
}

func TestRepository_ListVariants_NilPool(t *testing.T) {
	t.Parallel()

	repo := &Repository{}

	items, err := repo.ListVariants(context.Background())
	require.Nil(t, items)
	require.EqualError(t, err, "catalog postgres repository list variants: pool is nil")
}

func TestRepository_ListVariants_ReturnsOnlyActiveVariantsOrdered(t *testing.T) {
	pool := openTestPool(t)
	repo := NewRepository(pool)

	id3, code3, name3 := insertTestVariantWithOrder(t, pool, 30, true)
	_, _, _ = insertTestVariantWithOrder(t, pool, 5, false)
	id1, code1, name1 := insertTestVariantWithOrder(t, pool, 10, true)
	id2, code2, name2 := insertTestVariantWithOrder(t, pool, 20, true)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	items, err := repo.ListVariants(ctx)
	require.NoError(t, err)

	var got []flow.VariantListItem
	for _, item := range items {
		switch item.ID {
		case id1, id2, id3:
			got = append(got, item)
		}
	}

	require.Len(t, got, 3)

	require.Equal(t, flow.VariantListItem{
		ID:    id1,
		Code:  code1,
		Label: name1,
	}, got[0])

	require.Equal(t, flow.VariantListItem{
		ID:    id2,
		Code:  code2,
		Label: name2,
	}, got[1])

	require.Equal(t, flow.VariantListItem{
		ID:    id3,
		Code:  code3,
		Label: name3,
	}, got[2])
}
