package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/koha90/shopcore/internal/flow"
)

// ListDistrictProducts returns active products from one category that have placed variants in one district.
func (r *Repository) ListDistrictProducts(ctx context.Context, districtID, categoryID int) ([]flow.ProductListItem, error) {
	const op = "catalog postgres repository list district products"

	if r == nil {
		return nil, fmt.Errorf("%s: repository is nil", op)
	}
	if r.pool == nil {
		return nil, fmt.Errorf("%s: pool is nil", op)
	}

	const q = `
		select
			p.id,
			p.code,
			p.name
		from catalog_district_variants dv
		join catalog_variants v on v.id = dv.variant_id
		join catalog_products p on p.id = v.product_id
		where dv.district_id = $1
			and p.category_id = $2
			and dv.is_active = true
			and v.is_active = true
			and p.is_active = true
		group by p.id, p.code, p.name, p.sort_order
		order by p.created_at asc, p.id asc
	`

	rows, err := r.pool.Query(ctx, q, districtID, categoryID)
	if err != nil {
		return nil, fmt.Errorf("%s: query district products: %w", op, err)
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
		return item, err
	})
	if err != nil {
		return nil, fmt.Errorf("%s: collect district products: %w", op, err)
	}

	return items, nil
}
