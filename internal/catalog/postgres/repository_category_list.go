package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/koha90/shopcore/internal/flow"
)

// ListCategories returns active categories for admin selection flows.
func (r *Repository) ListCategories(ctx context.Context) ([]flow.CategoryListItem, error) {
	const op = "catalog postgres repository list categories"

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
		from catalog_categories
		where is_active = true
		order by created_at asc, id asc
	`

	rows, err := r.pool.Query(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("%s: query categories: %w", op, err)
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
		return nil, fmt.Errorf("%s: collect categories: %w", op, err)
	}

	return items, nil
}
