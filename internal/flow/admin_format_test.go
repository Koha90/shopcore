package flow

import "testing"

func TestFormatAdminFieldLine(t *testing.T) {
	t.Parallel()

	if got := formatAdminFieldLine("Район", "Центр"); got != "Район: Центр" {
		t.Fatalf("unexpected field line: %q", got)
	}

	if got := formatAdminFieldLine("Район", ""); got != "" {
		t.Fatalf("expected empty field line, got %q", got)
	}
}

func TestFormatAdminPriceValue(t *testing.T) {
	t.Parallel()

	if got := formatAdminPriceValue(5900); got != "5900 ₽" {
		t.Fatalf("unexpected price value: %q", got)
	}

	if got := formatAdminPriceValue(0); got != "" {
		t.Fatalf("expected empty price value, got %q", got)
	}
}

func TestFormatDistrictPlacementVariantActionLabel(t *testing.T) {
	t.Parallel()

	got := formatDistrictPlacementVariantActionLabel("L / 25 шт", 5900, "")
	if got != "L / 25 шт - 5900 ₽" {
		t.Fatalf("unexpected action label: %q", got)
	}
}
