package seed

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// EnsureCatalogDemoData inserts or updates demo catalog rows used for local
// development and manual verification of Postgres-backed flow navigation.
func EnsureCatalogDemoData(ctx context.Context, pool *pgxpool.Pool) error {
	const op = "seed catalog"

	if pool == nil {
		return fmt.Errorf("%s: pool is nil", op)
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: begin tx: %w", op, err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	moscowID, err := ensureCity(ctx, tx, citySeed{
		Code:      "moscow",
		Name:      "Москва",
		NameLatin: "Moscow",
		SortOrder: 10,
	})
	if err != nil {
		return err
	}

	spbID, err := ensureCity(ctx, tx, citySeed{
		Code:      "spb",
		Name:      "СПб",
		NameLatin: "Saint Petersburg",
		SortOrder: 20,
	})
	if err != nil {
		return err
	}

	flowersID, err := ensureCategory(ctx, tx, categorySeed{
		Code:        "flowers",
		Name:        "Цветы",
		NameLatin:   "Flowers",
		Description: "Букеты и цветочные композиции.",
		SortOrder:   10,
	})
	if err != nil {
		return err
	}

	giftsID, err := ensureCategory(ctx, tx, categorySeed{
		Code:        "gifts",
		Name:        "Подарки",
		NameLatin:   "Gifts",
		Description: "Подарочные наборы и приятные сюрпризы.",
		SortOrder:   20,
	})
	if err != nil {
		return err
	}

	centerID, err := ensureDistrict(ctx, tx, districtSeed{
		CityID:    moscowID,
		Code:      "center",
		Name:      "Центр",
		NameLatin: "Center",
		SortOrder: 10,
	})
	if err != nil {
		return err
	}

	southID, err := ensureDistrict(ctx, tx, districtSeed{
		CityID:    moscowID,
		Code:      "south",
		Name:      "Юг",
		NameLatin: "South",
		SortOrder: 20,
	})
	if err != nil {
		return err
	}

	petrogradkaID, err := ensureDistrict(ctx, tx, districtSeed{
		CityID:    spbID,
		Code:      "petrogradka",
		Name:      "Петроградка",
		NameLatin: "Petrogradka",
		SortOrder: 10,
	})
	if err != nil {
		return err
	}

	roseBoxID, err := ensureProduct(ctx, tx, productSeed{
		CategoryID:  flowersID,
		Code:        "rose-box",
		Name:        "Rose Box",
		NameLatin:   "Rose Box",
		Description: "Коробка роз.",
		SortOrder:   10,
	})
	if err != nil {
		return err
	}

	tulipMixID, err := ensureProduct(ctx, tx, productSeed{
		CategoryID:  flowersID,
		Code:        "tulip-mix",
		Name:        "Tulip Mix",
		NameLatin:   "Tulip Mix",
		Description: "Микс тюльпанов.",
		SortOrder:   20,
	})
	if err != nil {
		return err
	}

	giftBoxID, err := ensureProduct(ctx, tx, productSeed{
		CategoryID:  giftsID,
		Code:        "gift-box",
		Name:        "Gift Box",
		NameLatin:   "Gift Box",
		Description: "Подарочный набор.",
		ImageURL:    "assets/demo/catalog/products/gift-box.jpg",
		SortOrder:   10,
	})
	if err != nil {
		return err
	}

	peonySetID, err := ensureProduct(ctx, tx, productSeed{
		CategoryID:  flowersID,
		Code:        "peony-set",
		Name:        "Peony Set",
		NameLatin:   "Peony Set",
		Description: "Набор пионов.",
		SortOrder:   30,
	})
	if err != nil {
		return err
	}

	roseBoxSmallID, err := ensureVariant(ctx, tx, variantSeed{
		ProductID:   roseBoxID,
		Code:        "small",
		Name:        "S / 9 шт",
		NameLatin:   "S / 9 pcs",
		Description: "Компактная упаковка.",
		SortOrder:   10,
	})
	if err != nil {
		return err
	}

	roseBoxLargeID, err := ensureVariant(ctx, tx, variantSeed{
		ProductID:   roseBoxID,
		Code:        "large",
		Name:        "L / 25 шт",
		NameLatin:   "L / 25 pcs",
		Description: "Большая упаковка.",
		SortOrder:   20,
	})
	if err != nil {
		return err
	}

	tulipMixStandardID, err := ensureVariant(ctx, tx, variantSeed{
		ProductID:   tulipMixID,
		Code:        "standard",
		Name:        "Standard",
		NameLatin:   "Standard",
		Description: "Стандартный букет.",
		SortOrder:   10,
	})
	if err != nil {
		return err
	}

	giftBoxClassicID, err := ensureVariant(ctx, tx, variantSeed{
		ProductID:   giftBoxID,
		Code:        "classic",
		Name:        "Classic",
		NameLatin:   "Classic",
		Description: "Классический подарочный набор.",
		ImageURL:    "assets/demo/catalog/variants/gift-box.jpg",
		SortOrder:   10,
	})
	if err != nil {
		return err
	}

	peonySetPremiumID, err := ensureVariant(ctx, tx, variantSeed{
		ProductID:   peonySetID,
		Code:        "premium",
		Name:        "Premium",
		NameLatin:   "Premium",
		Description: "Премиальный набор пионов.",
		SortOrder:   10,
	})
	if err != nil {
		return err
	}

	if err := ensureDistrictVariant(ctx, tx, districtVariantSeed{
		DistrictID: centerID,
		VariantID:  roseBoxSmallID,
		Price:      2500,
	}); err != nil {
		return err
	}

	if err := ensureDistrictVariant(ctx, tx, districtVariantSeed{
		DistrictID: centerID,
		VariantID:  roseBoxLargeID,
		Price:      5900,
	}); err != nil {
		return err
	}

	if err := ensureDistrictVariant(ctx, tx, districtVariantSeed{
		DistrictID: southID,
		VariantID:  roseBoxLargeID,
		Price:      6100,
	}); err != nil {
		return err
	}

	if err := ensureDistrictVariant(ctx, tx, districtVariantSeed{
		DistrictID: centerID,
		VariantID:  tulipMixStandardID,
		Price:      3200,
	}); err != nil {
		return err
	}

	if err := ensureDistrictVariant(ctx, tx, districtVariantSeed{
		DistrictID: southID,
		VariantID:  giftBoxClassicID,
		Price:      4100,
	}); err != nil {
		return err
	}

	if err := ensureDistrictVariant(ctx, tx, districtVariantSeed{
		DistrictID: petrogradkaID,
		VariantID:  peonySetPremiumID,
		Price:      6800,
	}); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("%s: commit: %w", op, err)
	}

	return nil
}

type citySeed struct {
	Code      string
	Name      string
	NameLatin string
	SortOrder int
}

type categorySeed struct {
	Code        string
	Name        string
	NameLatin   string
	Description string
	SortOrder   int
}

type districtSeed struct {
	CityID    int
	Code      string
	Name      string
	NameLatin string
	SortOrder int
}

type productSeed struct {
	CategoryID  int
	Code        string
	Name        string
	NameLatin   string
	Description string
	ImageURL    string
	SortOrder   int
}

type variantSeed struct {
	ProductID   int
	Code        string
	Name        string
	NameLatin   string
	Description string
	ImageURL    string
	SortOrder   int
}

type districtVariantSeed struct {
	DistrictID int
	VariantID  int
	Price      int
}

func ensureCity(ctx context.Context, tx pgx.Tx, v citySeed) (int, error) {
	const q = `
		insert into cities (
			code, name, name_latin, is_active, sort_order, created_at, updated_at
		)
		values ($1, $2, $3, true, $4, now(), now())
		on conflict (code) do update set
			name = excluded.name,
			name_latin = excluded.name_latin,
			is_active = true,
			sort_order = excluded.sort_order,
			updated_at = now()
		returning id
	`

	var id int
	if err := tx.QueryRow(ctx, q, v.Code, v.Name, v.NameLatin, v.SortOrder).Scan(&id); err != nil {
		return 0, fmt.Errorf("seed city %q: %w", v.Code, err)
	}

	return id, nil
}

func ensureCategory(ctx context.Context, tx pgx.Tx, v categorySeed) (int, error) {
	const q = `
		insert into catalog_categories (
			code, name, name_latin, description, is_active, sort_order, created_at, updated_at
		)
		values ($1, $2, $3, $4, true, $5, now(), now())
		on conflict (code) do update set
			name = excluded.name,
			name_latin = excluded.name_latin,
			description = excluded.description,
			is_active = true,
			sort_order = excluded.sort_order,
			updated_at = now()
		returning id
	`

	var id int
	if err := tx.QueryRow(ctx, q, v.Code, v.Name, v.NameLatin, v.Description, v.SortOrder).Scan(&id); err != nil {
		return 0, fmt.Errorf("seed category %q: %w", v.Code, err)
	}

	return id, nil
}

func ensureDistrict(ctx context.Context, tx pgx.Tx, v districtSeed) (int, error) {
	const q = `
		insert into catalog_districts (
			city_id, code, name, name_latin, is_active, sort_order, created_at, updated_at
		)
		values ($1, $2, $3, $4, true, $5, now(), now())
		on conflict (city_id, code) do update set
			name = excluded.name,
			name_latin = excluded.name_latin,
			is_active = true,
			sort_order = excluded.sort_order,
			updated_at = now()
		returning id
	`

	var id int
	if err := tx.QueryRow(ctx, q, v.CityID, v.Code, v.Name, v.NameLatin, v.SortOrder).Scan(&id); err != nil {
		return 0, fmt.Errorf("seed district %q: %w", v.Code, err)
	}

	return id, nil
}

func ensureProduct(ctx context.Context, tx pgx.Tx, v productSeed) (int, error) {
	const q = `
		insert into catalog_products (
			category_id, code, name, name_latin, description, image_url, is_active, sort_order, created_at, updated_at
		)
		values ($1, $2, $3, $4, $5, $6, true, $7, now(), now())
		on conflict (category_id, code) do update set
			name = excluded.name,
			name_latin = excluded.name_latin,
			description = excluded.description,
			image_url = excluded.image_url,
			is_active = true,
			sort_order = excluded.sort_order,
			updated_at = now()
		returning id
	`

	var id int
	if err := tx.QueryRow(
		ctx,
		q,
		v.CategoryID,
		v.Code,
		v.Name,
		v.NameLatin,
		v.Description,
		v.ImageURL,
		v.SortOrder,
	).Scan(&id); err != nil {
		return 0, fmt.Errorf("seed product %q: %w", v.Code, err)
	}

	return id, nil
}

func ensureVariant(ctx context.Context, tx pgx.Tx, v variantSeed) (int, error) {
	const q = `
		insert into catalog_variants (
			product_id, code, name, name_latin, description, image_url, is_active, sort_order, created_at, updated_at
		)
		values ($1, $2, $3, $4, $5, $6, true, $7, now(), now())
		on conflict (product_id, code) do update set
			name = excluded.name,
			name_latin = excluded.name_latin,
			description = excluded.description,
			image_url = excluded.image_url,
			is_active = true,
			sort_order = excluded.sort_order,
			updated_at = now()
		returning id
	`

	var id int
	if err := tx.QueryRow(
		ctx,
		q,
		v.ProductID,
		v.Code,
		v.Name,
		v.NameLatin,
		v.Description,
		v.ImageURL,
		v.SortOrder,
	).Scan(&id); err != nil {
		return 0, fmt.Errorf("seed variant %q: %w", v.Code, err)
	}

	return id, nil
}

func ensureDistrictVariant(ctx context.Context, tx pgx.Tx, v districtVariantSeed) error {
	const q = `
		insert into catalog_district_variants (
			district_id, variant_id, price, is_active, created_at, updated_at
		)
		values ($1, $2, $3, true, now(), now())
		on conflict (district_id, variant_id) do update set
			price = excluded.price,
			is_active = true,
			updated_at = now()
	`

	if _, err := tx.Exec(ctx, q, v.DistrictID, v.VariantID, v.Price); err != nil {
		return fmt.Errorf(
			"seed district variant district=%d variant=%d: %w",
			v.DistrictID,
			v.VariantID,
			err,
		)
	}

	return nil
}
