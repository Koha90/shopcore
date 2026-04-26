package telegram

import "testing"

func TestNormalizeTelegramUsername(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in   string
		want string
	}{
		{"shop_main_bot", "shop_main_bot"},
		{"@shop_main_bot", "shop_main_bot"},
		{"  @shop_main_bot  ", "shop_main_bot"},
		{"", ""},
	}

	for _, tt := range tests {
		if got := normalizeTelegramUsername(tt.in); got != tt.want {
			t.Fatalf("normalizeTelegramUsername(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}
