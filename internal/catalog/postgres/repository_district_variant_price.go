package postgres

import (
	"context"
	"fmt"

	catalogservice "github.com/koha90/shopcore/internal/catalog/service"
)

// UpdateDistrictVariantPrice updates price for one district-variant placement row.
func (r *Repository) UpdateDistrictVariantPrice(
	ctx context.Context,
	params catalogservice.UpdateDistrictVariantPriceParams,
) error {
	const op = "catalog postgres repository update district variant price"

	if r == nil {
		return fmt.Errorf("%s: repository is nil", op)
	}
	if r.pool == nil {
		return fmt.Errorf("%s: pool is nil", op)
	}

	const q = `
		update catalog_district_variants
		set
			price = $3
			updated_at = now()
		where district_id = $1
			and variant_id = $2
	`

	tag, err := r.pool.Exec(
		ctx,
		q,
		params.DistrictID,
		params.VariantID,
		params.Price,
	)
	if err != nil {
		return fmt.Errorf(
			"update district variant price for district %d and variant %d: %w",
			params.DistrictID,
			params.VariantID,
			err,
		)
	}
	if tag.RowsAffected() == 0 {
		return ErrDistrictVariantNotFound
	}

	return nil
}
