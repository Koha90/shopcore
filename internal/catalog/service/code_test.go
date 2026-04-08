package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSuggestCategoryCode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "empty",
			in:   "",
			want: "",
		},
		{
			name: "spaces only",
			in:   "   ",
			want: "",
		},
		{
			name: "simple cyrillic",
			in:   "Цветы",
			want: "tsvety",
		},
		{
			name: "phrase",
			in:   "Тестовая категория",
			want: "testovaya-kategoriya",
		},
		{
			name: "yo and spaces",
			in:   "Ёлки и подарки",
			want: "elki-i-podarki",
		},
		{
			name: "mixed chars",
			in:   "VIP Букеты 24/7",
			want: "vip-bukety-24-7",
		},
		{
			name: "collapse separators",
			in:   "Цветы --- и   подарки",
			want: "tsvety-i-podarki",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := SuggestCode(tt.in)
			require.Equal(t, tt.want, got)
		})
	}
}
