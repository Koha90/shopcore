package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/koha90/shopcore/internal/flow"
)

// ListDistrictCategories returns active categories that have placed variants in one district.
func (r *Repository) ListDistrictCategories(ctx context.Context, districtID int) ([]flow.CategoryListItem, error) {
	const op = "catalog postgres repository list district categories"

	if r == nil {
		return nil, fmt.Errorf("%s: repository is nil", op)
	}
	if r.pool == nil {
		return nil, fmt.Errorf("%s: pool is nil", op)
	}

	const q = `
		select
			c.id,
			c.code,
			c.name
		from catalog_district_variants dv
		join catalog_variants v on v.id = dv.variant_id
		join catalog_products p on p.id = v.product_id
		join catalog_categories c on c.id = p.category_id
		where dv.district_id = $1
			and dv.is_active = true
			and v.is_active = true
			and p.is_active = true
			and c.is_active = true
		group by c.id, c.code, c.name, c.sort_order
		order by c.created_at asc, c.id asc
	`

	rows, err := r.pool.Query(ctx, q, districtID)
	if err != nil {
		return nil, fmt.Errorf("%s: query district categories: %w", op, err)
	}
	defer rows.Close()

	items, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (flow.CategoryListItem, error) {
		var item flow.CategoryListItem
		err := row.Scan(
			&item.ID,
			&item.Code,
			&item.Label,
		)
		if err != nil {
			return flow.CategoryListItem{}, err
		}
		return item, nil
	})
	if err != nil {
		return nil, fmt.Errorf("%s: collect district categories: %w", op, err)
	}

	return items, nil
}
