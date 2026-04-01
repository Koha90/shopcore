package postgres

import (
	"fmt"

	"github.com/koha90/shopcore/internal/flow"
)

func buildCatalog(
	cities []cityRow,
	categories []categoryRow,
	cityCategories []cityCategoryRow,
	districts []districtRow,
	products []productRow,
	variants []variantRow,
) flow.Catalog {
	categoriesByID := make(map[int]categoryRow, len(categories))
	for _, v := range categories {
		categoriesByID[v.ID] = v
	}

	districtsByCityID := make(map[int][]districtRow)
	for _, v := range districts {
		districtsByCityID[v.CityID] = append(districtsByCityID[v.CityID], v)
	}

	productsByDistrictID := make(map[int][]productRow)
	for _, v := range products {
		productsByDistrictID[v.DistrictID] = append(productsByDistrictID[v.DistrictID], v)
	}

	variantsByProductID := make(map[int][]variantRow)
	for _, v := range variants {
		variantsByProductID[v.ProductID] = append(variantsByProductID[v.ProductID], v)
	}

	categoryIDsByCityID := make(map[int][]int)
	for _, v := range cityCategories {
		categoryIDsByCityID[v.CityID] = append(categoryIDsByCityID[v.CityID], v.CategoryID)
	}

	var roots []flow.CatalogNode

	for _, city := range cities {
		categoryIDs := categoryIDsByCityID[city.ID]
		var categoryNodes []flow.CatalogNode

		for _, categoryID := range categoryIDs {
			category, ok := categoriesByID[categoryID]
			if !ok {
				continue
			}

			var districtNodes []flow.CatalogNode

			for _, district := range districtsByCityID[city.ID] {
				var productNodes []flow.CatalogNode

				for _, product := range productsByDistrictID[district.ID] {
					if product.CategoryID != category.ID {
						continue
					}

					var variantNodes []flow.CatalogNode
					for _, variant := range variantsByProductID[product.ID] {
						variantNodes = append(variantNodes, flow.CatalogNode{
							Level:       flow.LevelVariant,
							ID:          variant.Code,
							Label:       variant.Name,
							Description: variant.Description,
							PriceText:   formatPriceMinor(variant.PriceMinor, variant.CurrencyCode),
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

func formatPriceMinor(v int64, currency string) string {
	switch currency {
	case "RUB":
		return fmt.Sprintf("%d ₽", v)
	default:
		return fmt.Sprintf("%d %s", v, currency)
	}
}
