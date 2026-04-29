package postgres

import (
	"context"
	"fmt"

	catalogservice "github.com/koha90/shopcore/internal/catalog/service"
)

// UpdateVariantImage updates image_url for one variant.
func (r *Repository) UpdateVariantImage(
	ctx context.Context,
	params catalogservice.UpdateVariantImageParams,
) error {
	const op = "catalog postgres repository update variant image"

	if r == nil {
		return fmt.Errorf("%s: repository is nil", op)
	}
	if r.pool == nil {
		return fmt.Errorf("%s: pool is nil", op)
	}

	const q = `
    update catalog_variants
		set
			image_url = $2,
			updated_at = now()
		where id = $1
			and is_active = true
	`

	tag, err := r.pool.Exec(ctx, q, params.VariantID, params.ImageURL)
	if err != nil {
		return fmt.Errorf("%s: update variant %d image: %w", op, params.VariantID, err)
	}
	if tag.RowsAffected() == 0 {
		return ErrVariantNotFound
	}

	return nil
}
