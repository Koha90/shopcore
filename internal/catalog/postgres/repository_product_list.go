package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/koha90/shopcore/internal/flow"
)

// ListProducts returns active products for admin selection flows.
func (r *Repository) ListProducts(ctx context.Context) ([]flow.ProductListItem, error) {
	const op = "catalog postgres repository list products"

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
		from catalog_products
		where is_active = true
		order by sort_order asc, name asc
	`

	rows, err := r.pool.Query(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("%s: query row: %w", op, err)
	}
	defer rows.Close()

	items, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (flow.ProductListItem, error) {
		var item flow.ProductListItem
		err := row.Scan(
			&item.ID,
			&item.Code,
			&item.Label,
		)
		if err != nil {
			return flow.ProductListItem{}, err
		}
		return item, nil
	})
	if err != nil {
		return nil, fmt.Errorf("%s: collect products: %w", op, err)
	}

	return items, nil
}
