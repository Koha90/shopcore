package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/koha90/shopcore/internal/flow"
)

// ListDistricts returns active districts for admin selection flows.
func (r *Repository) ListDistricts(ctx context.Context) ([]flow.DistrictListItem, error) {
	const op = "catalog postgres repository list districts"

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
		from catalog_districts
		where is_active = true
		order by sort_order asc, name asc
	`

	rows, err := r.pool.Query(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("%s: query districts: %w", op, err)
	}
	defer rows.Close()

	items, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (flow.DistrictListItem, error) {
		var item flow.DistrictListItem
		err := row.Scan(
			&item.ID,
			&item.Code,
			&item.Label,
		)
		if err != nil {
			return flow.DistrictListItem{}, err
		}
		return item, nil
	})
	if err != nil {
		return nil, fmt.Errorf("%s: collect districts: %w", op, err)
	}

	return items, nil
}

// ListDistrictsByCity returns active districts of city for admin selection flows.
func (r *Repository) ListDistrictsByCity(ctx context.Context, cityID int) ([]flow.DistrictListItem, error) {
	const op = "catalog postgres repository list districts by city id"

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
		from catalog_districts
		where is_active = true and city_id = $1
		order by sort_order asc, name asc
	`

	rows, err := r.pool.Query(ctx, q, cityID)
	if err != nil {
		return nil, fmt.Errorf("%s: query districts: %w", op, err)
	}
	defer rows.Close()

	items, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (flow.DistrictListItem, error) {
		var item flow.DistrictListItem
		err := row.Scan(
			&item.ID,
			&item.Code,
			&item.Label,
		)
		if err != nil {
			return flow.DistrictListItem{}, err
		}
		return item, nil
	})
	if err != nil {
		return nil, fmt.Errorf("%s: collect districts: %w", op, err)
	}

	return items, nil
}
