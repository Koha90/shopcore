package postgres

import (
	"context"
	"fmt"

	catalogservice "github.com/koha90/shopcore/internal/catalog/service"
)

// CreateVariant inserts one catalog variant row.
func (r *Repository) CreateVariant(ctx context.Context, params catalogservice.CreateVariantParams) error {
	const op = "catalog postgres repository create variant"

	if r == nil {
		return fmt.Errorf("%s: repository is nil", op)
	}
	if r.pool == nil {
		return fmt.Errorf("%s: pool is nil", op)
	}

	const q = `
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
		values ($1, $2, $3, $4, $5, true, $6, now(), now())
	`

	_, err := r.pool.Exec(
		ctx,
		q,
		params.ProductID,
		params.Code,
		params.Name,
		params.NameLatin,
		params.Description,
		params.SortOrder,
	)
	if err != nil {
		return fmt.Errorf("create variant %q for product %d: %w", params.Code, params.ProductID, err)
	}

	return nil
}
