package postgres

import (
	"context"
	"fmt"

	catalogservice "github.com/koha90/shopcore/internal/catalog/service"
)

// CreateCity inserts one city row.
func (r *Repository) CreateCity(ctx context.Context, params catalogservice.CreateCityParams) error {
	const op = "catalog postgres repository create city"

	if r == nil {
		return fmt.Errorf("%s: repository is nil", op)
	}
	if r.pool == nil {
		return fmt.Errorf("%s: pool is nil", op)
	}

	const q = `
		insert into cities (
			code,
			name,
			name_latin,
			is_active,
			sort_order,
			created_at,
			updated_at
		)
		values ($1, $2, $3, true, $4, now(), now())
	`

	_, err := r.pool.Exec(
		ctx,
		q,
		params.Code,
		params.Name,
		params.NameLatin,
		params.SortOrder,
	)
	if err != nil {
		return fmt.Errorf("create city %q: %w", params.Code, err)
	}

	return nil
}
