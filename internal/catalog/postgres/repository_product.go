package postgres

import (
	"context"
	"fmt"

	catalogservice "github.com/koha90/shopcore/internal/catalog/service"
)

// CreateProduct inserts one catalog product row.
func (r *Repository) CreateProduct(ctx context.Context, params catalogservice.CreateProductParams) error {
	const op = "catalog postgres repository create product"

	if r == nil {
		return fmt.Errorf("%s: repository is nil", op)
	}
	if r.pool == nil {
		return fmt.Errorf("%s: pool is nil", op)
	}

	const q = `
		inser into catalog_products (
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
		values ($1, $2, $3, $4, $5, true, $6, now(), now())
	`

	_, err := r.pool.Exec(
		ctx,
		q,
		params.CategoryID,
		params.Code,
		params.Name,
		params.NameLatin,
		params.Description,
		params.SortOrder,
	)
	if err != nil {
		return fmt.Errorf("create product %q in category %d: %w", params.Code, params.CategoryID, err)
	}

	return nil
}
