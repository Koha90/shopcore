package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	catalogservice "github.com/koha90/shopcore/internal/catalog/service"
)

// Repository implements Postrgres-backed catalog write operations.
//
// Read-side catalog loading lives in Loader.
// Write-side operations are introduced gradually through explicit methods
// used by catalog application service.
type Repository struct {
	pool *pgxpool.Pool
}

// NewRepository constructs catalog Postrgres repository.
func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}

// CreateCategory inserts one catalog category row.
func (r *Repository) CreateCategory(ctx context.Context, params catalogservice.CreateCategoryParams) error {
	const op = "catalog postgres repository create category"

	if r == nil {
		return fmt.Errorf("%s: repository is nil", op)
	}
	if r.pool == nil {
		return fmt.Errorf("%s: pool is nil", op)
	}

	const q = `
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
		values ($1, $2, $3, $4, true, $5, now(), now())
	`

	_, err := r.pool.Exec(
		ctx,
		q,
		params.Code,
		params.Name,
		params.NameLatin,
		params.Description,
		// params.IsActive,
		params.SortOrder,
	)
	if err != nil {
		return fmt.Errorf("create category %q: %w", params.Code, err)
	}

	return nil
}
