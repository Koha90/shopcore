package postgres

import (
	"context"
	"fmt"

	catalogservice "github.com/koha90/shopcore/internal/catalog/service"
)

// CreateDistrict inserts one catalog district row.
func (r *Repository) CreateDistrict(ctx context.Context, params catalogservice.CreateDistrictParams) error {
	const op = "catalog postgres repository create district"

	if r == nil {
		return fmt.Errorf("%s: repository is nil", op)
	}
	if r.pool == nil {
		return fmt.Errorf("%s: pool is nil", op)
	}

	const q = `
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
		values ($1, $2, $3, $4, true, $5, now(), now())
	`

	_, err := r.pool.Exec(
		ctx,
		q,
		params.CityID,
		params.Code,
		params.Name,
		params.NameLatin,
		params.SortOrder,
	)
	if err != nil {
		return fmt.Errorf("create district %q in city %d: %w", params.Code, params.CityID, err)
	}

	return nil
}
