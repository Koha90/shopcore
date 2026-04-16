package flow

import "testing"

func TestNewImageMedia_EmptySourceReturnsNil(t *testing.T) {
	t.Parallel()

	got := NewImageMedia("", "welcome")
	if got != nil {
		t.Fatalf("expected nil media, got %#v", got)
	}
}

func TestNewImageMedia_BuildsImageMedia(t *testing.T) {
	t.Parallel()

	got := NewImageMedia("catalog/product/rose-box.jpg", "Rose Box")
	if got == nil {
		t.Fatal("expected media, got nil")
	}
	if got.Kind != MediaKindImage {
		t.Fatalf("expected media kind %q, got %q", MediaKindImage, got.Kind)
	}
	if got.Source != "catalog/product/rose-box.jpg" {
		t.Fatalf("unexpected source %q", got.Source)
	}
	if got.Alt != "Rose Box" {
		t.Fatalf("unexpected alt %q", got.Alt)
	}
}
