package postgres

import (
	"fmt"

	"github.com/koha90/shopcore/internal/flow"
)

// buildCatalog maps flat relational rows into flow.Catalog tree.
//
// Only valid branches are included:
//   - placement without referenced district is skipped
//   - placement without referenced variant is skipped
//   - variant without referenced product is skipped
//   - product without referenced category is skipped
//   - district without referenced city is skipped
//   - products without placed variants are skipped
//   - districts without placed products are skipped
//   - categories without district/product branch are skipped
//   - cities without children are skipped
func buildCatalog(
	cities []cityRow,
	categories []categoryRow,
	districts []districtRow,
	products []productRow,
	variants []variantRow,
	districtVariants []districtVariantRow,
) flow.Catalog {
	citiesByID := make(map[int]cityRow, len(cities))
	for _, v := range cities {
		citiesByID[v.ID] = v
	}

	categoriesByID := make(map[int]categoryRow, len(categories))
	for _, v := range categories {
		categoriesByID[v.ID] = v
	}

	districtsByID := make(map[int]districtRow, len(districts))
	districtsByCityID := make(map[int][]districtRow)
	for _, v := range districts {
		districtsByID[v.ID] = v
		districtsByCityID[v.CityID] = append(districtsByCityID[v.CityID], v)
	}

	productsByID := make(map[int]productRow, len(products))
	for _, v := range products {
		productsByID[v.ID] = v
	}

	variantsByID := make(map[int]variantRow, len(variants))
	for _, v := range variants {
		variantsByID[v.ID] = v
	}

	type productBucket struct {
		variantPrices map[int]int
	}

	type districtBucket struct {
		products map[int]*productBucket
	}

	type categoryBucket struct {
		districts map[int]*districtBucket
	}

	type cityBucket struct {
		categories map[int]*categoryBucket
	}

	cityBuckets := make(map[int]*cityBucket)

	for _, placement := range districtVariants {
		district, ok := districtsByID[placement.DistrictID]
		if !ok {
			continue
		}

		variant, ok := variantsByID[placement.VariantID]
		if !ok {
			continue
		}

		product, ok := productsByID[variant.ProductID]
		if !ok {
			continue
		}

		category, ok := categoriesByID[product.CategoryID]
		if !ok {
			continue
		}

		_, ok = citiesByID[district.CityID]
		if !ok {
			continue
		}

		cb, ok := cityBuckets[district.CityID]
		if !ok {
			cb = &cityBucket{
				categories: make(map[int]*categoryBucket),
			}
			cityBuckets[district.CityID] = cb
		}

		catb, ok := cb.categories[category.ID]
		if !ok {
			catb = &categoryBucket{
				districts: make(map[int]*districtBucket),
			}
			cb.categories[category.ID] = catb
		}

		db, ok := catb.districts[district.ID]
		if !ok {
			db = &districtBucket{
				products: make(map[int]*productBucket),
			}
			catb.districts[district.ID] = db
		}

		pb, ok := db.products[product.ID]
		if !ok {
			pb = &productBucket{
				variantPrices: make(map[int]int),
			}
			db.products[product.ID] = pb
		}

		pb.variantPrices[variant.ID] = placement.Price
	}

	var roots []flow.CatalogNode

	for _, city := range cities {
		cb, ok := cityBuckets[city.ID]
		if !ok {
			continue
		}

		var categoryNodes []flow.CatalogNode

		for _, category := range categories {
			catb, ok := cb.categories[category.ID]
			if !ok {
				continue
			}

			var districtNodes []flow.CatalogNode

			for _, district := range districtsByCityID[city.ID] {
				db, ok := catb.districts[district.ID]
				if !ok {
					continue
				}

				var productNodes []flow.CatalogNode

				for _, product := range products {
					if product.CategoryID != category.ID {
						continue
					}

					pb, ok := db.products[product.ID]
					if !ok {
						continue
					}

					var variantNodes []flow.CatalogNode

					for _, variant := range variants {
						if variant.ProductID != product.ID {
							continue
						}

						price, ok := pb.variantPrices[variant.ID]
						if !ok {
							continue
						}

						variantNodes = append(variantNodes, flow.CatalogNode{
							Level:       flow.LevelVariant,
							ID:          variant.Code,
							Label:       variant.Name,
							Description: variant.Description,
							ImageURL:    variant.ImageURL,
							PriceText:   formatPrice(price),
						})
					}

					if len(variantNodes) == 0 {
						continue
					}

					productNodes = append(productNodes, flow.CatalogNode{
						Level:       flow.LevelProduct,
						ID:          product.Code,
						Label:       product.Name,
						Description: product.Description,
						ImageURL:    product.ImageURL,
						Children:    variantNodes,
					})
				}

				if len(productNodes) == 0 {
					continue
				}

				districtNodes = append(districtNodes, flow.CatalogNode{
					Level:    flow.LevelDistrict,
					ID:       district.Code,
					Label:    district.Name,
					Children: productNodes,
				})
			}

			if len(districtNodes) == 0 {
				continue
			}

			categoryNodes = append(categoryNodes, flow.CatalogNode{
				Level:       flow.LevelCategory,
				ID:          category.Code,
				Label:       category.Name,
				Description: category.Description,
				Children:    districtNodes,
			})
		}

		if len(categoryNodes) == 0 {
			continue
		}

		roots = append(roots, flow.CatalogNode{
			Level:    flow.LevelCity,
			ID:       city.Code,
			Label:    city.Name,
			Children: categoryNodes,
		})
	}

	return flow.Catalog{
		Schema: flow.DemoCatalogSchema(),
		Roots:  roots,
	}
}

// formatPrice formats integer price storage into display text.
func formatPrice(v int) string {
	return fmt.Sprintf("%d ₽", v)
}
