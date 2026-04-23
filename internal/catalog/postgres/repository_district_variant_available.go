package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/koha90/shopcore/internal/flow"
)

func (r *Repository) ListAvailableVariantsForDistrictProduct(
	ctx context.Context,
	districtID, productID int,
) ([]flow.VariantListItem, error) {
	const op = "catalog postgres repository list available variants for district product"

	if r == nil {
		return nil, fmt.Errorf("%s: repository is nil", op)
	}
	if r.pool == nil {
		return nil, fmt.Errorf("%s: pool is nil", op)
	}
	if districtID <= 0 {
		return nil, fmt.Errorf("%s: district id is invalid", op)
	}
	if productID <= 0 {
		return nil, fmt.Errorf("%s: product id is invalid", op)
	}

	const q = `
		select
			v.id,
			v.code,
			v.name,
			p.name as product_name
		from catalog_variants v
		join catalog_products p on p.id = v.product_id
		where v.is_active = true
			and p.is_active = true
			and v.product_id = $2
			and not exists (
				select 1
				from catalog_district_variants dv
				where dv.district_id = $1
					and dv.variant_id = v.id
					and dv.is_active = true
			)
		order by v.sort_order asc, v.name asc
	`

	rows, err := r.pool.Query(ctx, q, districtID, productID)
	if err != nil {
		return nil, fmt.Errorf("%s: query available variants: %w", op, err)
	}
	defer rows.Close()

	items, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (flow.VariantListItem, error) {
		var item flow.VariantListItem

		if err := row.Scan(
			&item.ID,
			&item.Code,
			&item.Label,
			&item.ProductLabel,
		); err != nil {
			return flow.VariantListItem{}, err
		}

		return item, nil
	})
	if err != nil {
		return nil, fmt.Errorf("%s: collect available variants: %w", op, err)
	}

	return items, nil
}
