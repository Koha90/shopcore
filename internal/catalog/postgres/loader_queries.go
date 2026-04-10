package postgres

import "context"

func (l *Loader) loadCities(ctx context.Context) ([]cityRow, error) {
	rows, err := l.pool.Query(ctx, `
		select id, code, name, name_latin, sort_order
		from cities
		where is_active = true
		order by sort_order, id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []cityRow
	for rows.Next() {
		var v cityRow
		if err := rows.Scan(
			&v.ID,
			&v.Code,
			&v.Name,
			&v.NameLatin,
			&v.SortOrder,
		); err != nil {
			return nil, err
		}
		out = append(out, v)
	}

	return out, rows.Err()
}

func (l *Loader) loadCategories(ctx context.Context) ([]categoryRow, error) {
	rows, err := l.pool.Query(ctx, `
		select id, code, name, name_latin, description, sort_order
		from catalog_categories
		where is_active = true
		order by sort_order, id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []categoryRow
	for rows.Next() {
		var v categoryRow
		if err := rows.Scan(
			&v.ID,
			&v.Code,
			&v.Name,
			&v.NameLatin,
			&v.Description,
			&v.SortOrder,
		); err != nil {
			return nil, err
		}
		out = append(out, v)
	}

	return out, rows.Err()
}

func (l *Loader) loadDistricts(ctx context.Context) ([]districtRow, error) {
	rows, err := l.pool.Query(ctx, `
		select id, city_id, code, name, name_latin, sort_order
		from catalog_districts
		where is_active = true
		order by sort_order, id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []districtRow
	for rows.Next() {
		var v districtRow
		if err := rows.Scan(
			&v.ID,
			&v.CityID,
			&v.Code,
			&v.Name,
			&v.NameLatin,
			&v.SortOrder,
		); err != nil {
			return nil, err
		}
		out = append(out, v)
	}

	return out, rows.Err()
}

func (l *Loader) loadDistrictVariants(ctx context.Context) ([]districtVariantRow, error) {
	rows, err := l.pool.Query(ctx, `
		select district_id, variant_id, price
		from catalog_district_variants
		where is_active = true
		order by district_id, variant_id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []districtVariantRow
	for rows.Next() {
		var v districtVariantRow
		if err := rows.Scan(
			&v.DistrictID,
			&v.VariantID,
			&v.Price,
		); err != nil {
			return nil, err
		}
		out = append(out, v)
	}

	return out, rows.Err()
}

func (l *Loader) loadProducts(ctx context.Context) ([]productRow, error) {
	rows, err := l.pool.Query(ctx, `
		select id, category_id, code, name, name_latin, description, sort_order
		from catalog_products
		where is_active = true
		order by sort_order, id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []productRow
	for rows.Next() {
		var v productRow
		if err := rows.Scan(
			&v.ID,
			&v.CategoryID,
			&v.Code,
			&v.Name,
			&v.NameLatin,
			&v.Description,
			&v.SortOrder,
		); err != nil {
			return nil, err
		}
		out = append(out, v)
	}

	return out, rows.Err()
}

func (l *Loader) loadVariants(ctx context.Context) ([]variantRow, error) {
	rows, err := l.pool.Query(ctx, `
		select id, product_id, code, name, name_latin, description, sort_order
		from catalog_variants
		where is_active = true
		order by sort_order, id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []variantRow
	for rows.Next() {
		var v variantRow
		if err := rows.Scan(
			&v.ID,
			&v.ProductID,
			&v.Code,
			&v.Name,
			&v.NameLatin,
			&v.Description,
			&v.SortOrder,
		); err != nil {
			return nil, err
		}
		out = append(out, v)
	}

	return out, rows.Err()
}
