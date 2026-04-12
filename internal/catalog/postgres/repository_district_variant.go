package postgres

import (
	"context"
	"fmt"

	catalogservice "github.com/koha90/shopcore/internal/catalog/service"
)

// CreateDistrictVariant inserts one district-variant placement row.
func (r *Repository) CreateDistrictVariant(
	ctx context.Context,
	params catalogservice.CreateDistrictVariantParams,
) error {
	const op = "catalog postgres repository create district variant"

	if r == nil {
		return fmt.Errorf("%s: repository is nil", op)
	}
	if r.pool == nil {
		return fmt.Errorf("%s: pool is nil", op)
	}

	const q = `
		insert into catalog_district_variants (
			district_id,
			variant_id, 
			price,
			is_active,
			created_at,
			updated_at
		)
		values ($1, $2, $3, true, now(), now())
	`

	_, err := r.pool.Exec(
		ctx,
		q,
		params.DistrictID,
		params.VariantID,
		params.Price,
	)
	if err != nil {
		return fmt.Errorf(
			"create district variant for district %d and variant %d: %w",
			params.DistrictID,
			params.VariantID,
			err,
		)
	}

	return nil
}
