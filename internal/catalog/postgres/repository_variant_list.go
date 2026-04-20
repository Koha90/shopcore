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
			v.id,
			v.code,
			v.name,
			p.name as product_name
		from catalog_variants v
		join catalog_products p on p.id = v.product_id
		where v.is_active = true
			and p.is_active = true
		order by p.sort_order asc, p.name asc, v.sort_order asc, v.name asc
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
			&item.ProductLabel,
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

func (r *Repository) ListVariantsByProduct(ctx context.Context, productID int) ([]flow.VariantListItem, error) {
	const op = "catalog postgres repository list variants by product"

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
			v.name,
			p.name as product_name
		from catalog_variants v
		join catalog_products p on p.id = v.product_id
		where v.is_active = true
			and p.is_active = true
			and v.product_id = $1
		order by v.sort_order asc, v.name asc
	`

	rows, err := r.pool.Query(ctx, q, productID)
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
			&item.ProductLabel,
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
