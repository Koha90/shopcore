package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/koha90/shopcore/internal/flow"
)

// ListCities returns active cities for admin selection flows.
func (r *Repository) ListCities(ctx context.Context) ([]flow.CityListItem, error) {
	const op = "catalog postgres repository list cities"

	if r == nil {
		return nil, fmt.Errorf("%s: repository is nil", op)
	}
	if r.pool == nil {
		return nil, fmt.Errorf("%s: pool is nil", op)
	}

	const q = `
		select
			id,
			code,
			name
		from cities
		where is_active = true
		order by sort_order asc, name asc
	`

	rows, err := r.pool.Query(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("%s: query cities: %w", op, err)
	}
	defer rows.Close()

	items, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (flow.CityListItem, error) {
		var item flow.CityListItem
		err := row.Scan(
			&item.ID,
			&item.Code,
			&item.Label,
		)
		if err != nil {
			return flow.CityListItem{}, err
		}
		return item, nil
	})
	if err != nil {
		return nil, fmt.Errorf("%s: collect cities: %w", op, err)
	}

	return items, nil
}
