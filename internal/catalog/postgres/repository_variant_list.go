package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/koha90/shopcore/internal/flow"
)

// ListVariants returns active variants for admin selection flows.
func (r *Repository) ListVariants(ctx context.Context) ([]flow.VariantListItem, error) {
	const op = "catalog postgres repository list variants"

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
		from catalog_variants
		where is_active = true
		order by sort_order asc, name asc
	`

	rows, err := r.pool.Query(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("%s: query variants: %w", op, err)
	}
	defer rows.Close()

	items, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (flow.VariantListItem, error) {
		var item flow.VariantListItem
		err := row.Scan(
			&item.ID,
			&item.Code,
			&item.Label,
		)
		if err != nil {
			return flow.VariantListItem{}, err
		}
		return item, nil
	})
	if err != nil {
		return nil, fmt.Errorf("%s: collect variants: %w", op, err)
	}

	return items, nil
}
