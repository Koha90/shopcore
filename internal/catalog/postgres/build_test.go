package postgres

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/koha90/shopcore/internal/flow"
)

func TestFormatPriceMinor(t *testing.T) {
	require.Equal(t, "5900 ₽", formatPriceMinor(5900, "RUB"))
	require.Equal(t, "42 USD", formatPriceMinor(42, "USD"))
}

func TestBuildCatalog_HappyPath(t *testing.T) {
	catalog := buildCatalog(
		[]cityRow{
			{ID: 1, Code: "moscow", Name: "Москва"},
		},
		[]categoryRow{
			{ID: 10, Code: "flowers", Name: "Цветы", Description: "Категория цветов"},
		},
		[]cityCategoryRow{
			{CityID: 1, CategoryID: 10},
		},
		[]districtRow{
			{ID: 100, CityID: 1, Code: "center", Name: "Центр"},
		},
		[]productRow{
			{ID: 1000, CategoryID: 10, DistrictID: 100, Code: "rose-box", Name: "Rose Box", Description: "Коробка роз"},
		},
		[]variantRow{
			{ID: 10000, ProductID: 1000, Code: "large", Name: "L / 25 шт", Description: "Большая упаковка", PriceMinor: 5900, CurrencyCode: "RUB"},
		},
	)

	require.Equal(t, flow.DemoCatalogSchema(), catalog.Schema)
	require.Len(t, catalog.Roots, 1)

	city := catalog.Roots[0]
	require.Equal(t, flow.LevelCity, city.Level)
	require.Equal(t, "moscow", city.ID)
	require.Equal(t, "Москва", city.Label)

	require.Len(t, city.Children, 1)
	category := city.Children[0]
	require.Equal(t, flow.LevelCategory, category.Level)
	require.Equal(t, "flowers", category.ID)
	require.Equal(t, "Цветы", category.Label)
	require.Equal(t, "Категория цветов", category.Description)

	require.Len(t, category.Children, 1)
	district := category.Children[0]
	require.Equal(t, flow.LevelDistrict, district.Level)
	require.Equal(t, "center", district.ID)
	require.Equal(t, "Центр", district.Label)

	require.Len(t, district.Children, 1)
	product := district.Children[0]
	require.Equal(t, flow.LevelProduct, product.Level)
	require.Equal(t, "rose-box", product.ID)
	require.Equal(t, "Rose Box", product.Label)
	require.Equal(t, "Коробка роз", product.Description)

	require.Len(t, product.Children, 1)
	variant := product.Children[0]
	require.Equal(t, flow.LevelVariant, variant.Level)
	require.Equal(t, "large", variant.ID)
	require.Equal(t, "L / 25 шт", variant.Label)
	require.Equal(t, "Большая упаковка", variant.Description)
	require.Equal(t, "5900 ₽", variant.PriceText)
}

func TestBuildCatalog_SkipsProductWithoutVariants(t *testing.T) {
	catalog := buildCatalog(
		[]cityRow{
			{ID: 1, Code: "moscow", Name: "Москва"},
		},
		[]categoryRow{
			{ID: 10, Code: "flowers", Name: "Цветы"},
		},
		[]cityCategoryRow{
			{CityID: 1, CategoryID: 10},
		},
		[]districtRow{
			{ID: 100, CityID: 1, Code: "center", Name: "Центр"},
		},
		[]productRow{
			{ID: 1000, CategoryID: 10, DistrictID: 100, Code: "rose-box", Name: "Rose Box"},
		},
		nil,
	)

	require.Empty(t, catalog.Roots)
}

func TestBuildCatalog_SkipsCategoryWithoutValidDistrictBranch(t *testing.T) {
	catalog := buildCatalog(
		[]cityRow{
			{ID: 1, Code: "moscow", Name: "Москва"},
		},
		[]categoryRow{
			{ID: 10, Code: "flowers", Name: "Цветы"},
		},
		[]cityCategoryRow{
			{CityID: 1, CategoryID: 10},
		},
		[]districtRow{
			{ID: 100, CityID: 1, Code: "center", Name: "Центр"},
		},
		nil,
		nil,
	)

	require.Empty(t, catalog.Roots)
}

func TestBuildCatalog_SkipsUnknownCategoryLinks(t *testing.T) {
	catalog := buildCatalog(
		[]cityRow{
			{ID: 1, Code: "moscow", Name: "Москва"},
		},
		nil,
		[]cityCategoryRow{
			{CityID: 1, CategoryID: 10},
		},
		nil,
		nil,
		nil,
	)

	require.Empty(t, catalog.Roots)
}

func TestBuildCatalog_SkipsCityWithoutChildren(t *testing.T) {
	catalog := buildCatalog(
		[]cityRow{
			{ID: 1, Code: "moscow", Name: "Москва"},
			{ID: 2, Code: "spb", Name: "СПб"},
		},
		[]categoryRow{
			{ID: 10, Code: "flowers", Name: "Цветы"},
		},
		[]cityCategoryRow{
			{CityID: 1, CategoryID: 10},
		},
		[]districtRow{
			{ID: 100, CityID: 1, Code: "center", Name: "Центр"},
		},
		[]productRow{
			{ID: 1000, CategoryID: 10, DistrictID: 100, Code: "rose-box", Name: "Rose Box"},
		},
		[]variantRow{
			{ID: 10000, ProductID: 1000, Code: "large", Name: "L / 25 шт", PriceMinor: 5900, CurrencyCode: "RUB"},
		},
	)

	require.Len(t, catalog.Roots, 1)
	require.Equal(t, "moscow", catalog.Roots[0].ID)
}
