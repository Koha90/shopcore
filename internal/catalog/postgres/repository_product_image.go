package postgres

import (
	"context"
	"fmt"

	catalogservice "github.com/koha90/shopcore/internal/catalog/service"
)

// UpdateProductImage updates image_url for one active catalog product.
func (r *Repository) UpdateProductImage(
	ctx context.Context,
	params catalogservice.UpdateProductImageParams,
) error {
	const op = "catalog postgres repository update product image"

	if r == nil {
		return fmt.Errorf("%s: repository is nil", op)
	}
	if r.pool == nil {
		return fmt.Errorf("%s: pool is nil", op)
	}

	const q = `
		update catalog_products
		set
			image_url = $2,
			updated_at = now()
		where id = $1
			and is_active = true
	`

	tag, err := r.pool.Exec(ctx, q, params.ProductID, params.ImageURL)
	if err != nil {
		return fmt.Errorf("%s: update product %d image: %w", op, params.ProductID, err)
	}
	if tag.RowsAffected() == 0 {
		return ErrProductNotFound
	}

	return nil
}
