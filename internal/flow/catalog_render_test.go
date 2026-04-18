package flow

import "testing"

func TestCatalogChildButtonLabel_VariantWithPrice(t *testing.T) {
	t.Parallel()

	got := catalogChildButtonLabel(CatalogNode{
		Level:     LevelVariant,
		Label:     "Classic",
		PriceText: "4100 ₽",
	})

	if got != "Classic - 4100 ₽" {
		t.Fatalf("unexpected label %q", got)
	}
}

func TestCatalogChildButtonLabel_VariantWithoutPrice(t *testing.T) {
	t.Parallel()

	got := catalogChildButtonLabel(CatalogNode{
		Level: LevelVariant,
		Label: "Classic",
	})

	if got != "Classic" {
		t.Fatalf("unexpected label %q", got)
	}
}

func TestCatalogChildButtonLabel_NonVariant(t *testing.T) {
	t.Parallel()

	got := catalogChildButtonLabel(CatalogNode{
		Level:     LevelProduct,
		Label:     "Gift Box",
		PriceText: "4100 ₽",
	})

	if got != "Gift Box" {
		t.Fatalf("unexpected label %q", got)
	}
}
