package postgres

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/koha90/shopcore/internal/flow"
)

func TestFormatPrice(t *testing.T) {
	require.Equal(t, "2500 ₽", formatPrice(2500))
	require.Equal(t, "0 ₽", formatPrice(0))
}

func TestBuildCatalog_BuildsPlacedBranch(t *testing.T) {
	catalog := buildCatalog(
		[]cityRow{
			{ID: 1, Code: "moscow", Name: "Москва"},
		},
		[]categoryRow{
			{ID: 10, Code: "flowers", Name: "Цветы", Description: "Букеты"},
		},
		[]districtRow{
			{ID: 100, CityID: 1, Code: "center", Name: "Центр"},
		},
		[]productRow{
			{ID: 1000, CategoryID: 10, Code: "rose-box", Name: "Rose Box", Description: "Коробка роз", ImageURL: ""},
		},
		[]variantRow{
			{ID: 10000, ProductID: 1000, Code: "large", Name: "L / 25 шт", Description: "Большая упаковка", ImageURL: ""},
		},
		[]districtVariantRow{
			{DistrictID: 100, VariantID: 10000, Price: 5900},
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
	require.Equal(t, "Букеты", category.Description)
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
	require.Nil(t, product.Media)
	require.Len(t, product.Children, 1)

	variant := product.Children[0]
	require.Equal(t, flow.LevelVariant, variant.Level)
	require.Equal(t, "large", variant.ID)
	require.Equal(t, "L / 25 шт", variant.Label)
	require.Equal(t, "Большая упаковка", variant.Description)
	require.Nil(t, product.Media)
	require.Equal(t, "5900 ₽", variant.PriceText)
}

func TestBuildCatalog_SkipsBranchWithoutPlacement(t *testing.T) {
	catalog := buildCatalog(
		[]cityRow{
			{ID: 1, Code: "moscow", Name: "Москва"},
		},
		[]categoryRow{
			{ID: 10, Code: "flowers", Name: "Цветы", Description: "Букеты"},
		},
		[]districtRow{
			{ID: 100, CityID: 1, Code: "center", Name: "Центр"},
		},
		[]productRow{
			{ID: 1000, CategoryID: 10, Code: "rose-box", Name: "Rose Box", Description: "Коробка роз"},
		},
		[]variantRow{
			{ID: 10000, ProductID: 1000, Code: "large", Name: "L / 25 шт", Description: "Большая упаковка"},
		},
		nil,
	)

	require.Empty(t, catalog.Roots)
}

func TestBuildCatalog_SkipsBrokenPlacementReferences(t *testing.T) {
	catalog := buildCatalog(
		[]cityRow{
			{ID: 1, Code: "moscow", Name: "Москва"},
		},
		[]categoryRow{
			{ID: 10, Code: "flowers", Name: "Цветы", Description: "Букеты"},
		},
		[]districtRow{
			{ID: 100, CityID: 1, Code: "center", Name: "Центр"},
		},
		[]productRow{
			{ID: 1000, CategoryID: 10, Code: "rose-box", Name: "Rose Box", Description: "Коробка роз"},
		},
		[]variantRow{
			{ID: 10000, ProductID: 1000, Code: "large", Name: "L / 25 шт", Description: "Большая упаковка"},
		},
		[]districtVariantRow{
			{DistrictID: 999, VariantID: 10000, Price: 5900},
			{DistrictID: 100, VariantID: 99999, Price: 5900},
		},
	)

	require.Empty(t, catalog.Roots)
}

func TestBuildCatalogNodeMedia_EmptySourceReturnsNil(t *testing.T) {
	t.Parallel()

	require.Nil(t, buildCatalogNodeMedia("", "Gift Box"))
}

func TestBuildCatalogNodeMedia_UsesProvidedAlt(t *testing.T) {
	t.Parallel()

	got := buildCatalogNodeMedia("assets/catalog/products/gift-box.png", "Gift Box")
	require.NotNil(t, got)
	require.Equal(t, "assets/catalog/products/gift-box.png", got.ImageSource)
	require.Equal(t, "Gift Box", got.ImageAlt)
}

func TestBuildCatalogNodeMedia_FallbackAlt(t *testing.T) {
	t.Parallel()

	got := buildCatalogNodeMedia("assets/catalog/products/gift-box.png", "")
	require.NotNil(t, got)
	require.Equal(t, "image", got.ImageAlt)
}

func TestBuildCatalog_MapsProductAndVariantMedia(t *testing.T) {
	t.Parallel()

	catalog := buildCatalog(
		[]cityRow{
			{ID: 1, Code: "spb", Name: "СПб"},
		},
		[]categoryRow{
			{ID: 10, Code: "gifts", Name: "Подарки", Description: "Подарочные наборы"},
		},
		[]districtRow{
			{ID: 100, CityID: 1, Code: "south", Name: "Юг"},
		},
		[]productRow{
			{
				ID:          1000,
				CategoryID:  10,
				Code:        "gift-box",
				Name:        "Gift Box",
				Description: "Подарочный набор",
				ImageURL:    "assets/catalog/products/gift-box.png",
			},
		},
		[]variantRow{
			{
				ID:          10000,
				ProductID:   1000,
				Code:        "classic",
				Name:        "Classic",
				Description: "Классический подарочный набор",
				ImageURL:    "assets/catalog/variants/classic.png",
			},
		},
		[]districtVariantRow{
			{DistrictID: 100, VariantID: 10000, Price: 4100},
		},
	)

	require.Len(t, catalog.Roots, 1)

	product := catalog.Roots[0].Children[0].Children[0].Children[0]
	require.Equal(t, flow.LevelProduct, product.Level)
	require.NotNil(t, product.Media)
	require.Equal(t, "assets/catalog/products/gift-box.png", product.Media.ImageSource)
	require.Equal(t, "Gift Box", product.Media.ImageAlt)

	require.Len(t, product.Children, 1)

	variant := product.Children[0]
	require.Equal(t, flow.LevelVariant, variant.Level)
	require.NotNil(t, variant.Media)
	require.Equal(t, "assets/catalog/variants/classic.png", variant.Media.ImageSource)
	require.Equal(t, "Classic", variant.Media.ImageAlt)
	require.Equal(t, "4100 ₽", variant.PriceText)
}
