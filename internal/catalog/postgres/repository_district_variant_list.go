package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/koha90/shopcore/internal/flow"
)

// ListDistrictVariants returns active placed variants for one product in one district.
func (r *Repository) ListDistrictVariants(ctx context.Context, districtID, productID int) ([]flow.VariantListItem, error) {
	const op = "catalog postgres repository list district variants"

	if r == nil {
		return nil, fmt.Errorf("%s: repository is nil", op)
	}
	if r.pool == nil {
		return nil, fmt.Errorf("%s: pool is nil", op)
	}

	const q = `
		select
			v.id,
			v.code,
			v.name
		from catalog_district_variants dv
		join catalog_variants v on v.id = dv.variant_id
		join catalog_products p on p.id = v.product_id
		where dv.district_id = $1
			and p.id = $2
			and dv.is_active = true
			and v.is_active = true
			and p.is_active = true
		group by v.id, v.code, v.name, v.sort_order
		order by v.sort_order asc, v.name asc
	`

	rows, err := r.pool.Query(ctx, q, districtID, productID)
	if err != nil {
		return nil, fmt.Errorf("%s: query district variants: %w", op, err)
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
		return nil, fmt.Errorf("%s: collect district variants: %w", op, err)
	}

	return items, nil
}
