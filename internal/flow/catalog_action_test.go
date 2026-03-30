package flow

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCatalogSelectAction(t *testing.T) {
	got := catalogSelectAction(LevelCity, "moscow")
	require.Equal(t, ActionID("catalog:select:city:moscow"), got)
}

func TestParseCatalogSelectAction(t *testing.T) {
	level, id, ok := parseCatalogSelectAction(ActionID("catalog:select:product:rose-box"))
	require.True(t, ok)
	require.Equal(t, LevelProduct, level)
	require.Equal(t, "rose-box", id)
}

func TestParseCatalogSelectAction_Invalid(t *testing.T) {
	tests := []ActionID{
		"",
		"catalog:select:",
		"catalog:select:city",
		"catalog:select:city:moscow:extra",
		"city:moscow",
	}

	for _, tt := range tests {
		t.Run(string(tt), func(t *testing.T) {
			level, id, ok := parseCatalogSelectAction(tt)
			require.False(t, ok)
			require.Empty(t, level)
			require.Empty(t, id)
		})
	}
}

func TestCatalogScreen_ParseRoundTrip(t *testing.T) {
	path := CatalogPath{
		{Level: LevelCity, ID: "moscow"},
		{Level: LevelCategory, ID: "flowers"},
		{Level: LevelDistrict, ID: "center"},
	}

	screen := catalogScreen(path)
	require.Equal(t, ScreenID("catalog:screen:city=moscow;category=flowers;district=center"), screen)

	parsed, ok := parseCatalogScreen(screen)
	require.True(t, ok)
	require.Equal(t, path, parsed)
}

func TestParseCatalogScreen_Invalid(t *testing.T) {
	tests := []ScreenID{
		"",
		"catalog:screen:",
		"screen:city=moscow",
	}

	for _, tt := range tests {
		t.Run(string(tt), func(t *testing.T) {
			path, ok := parseCatalogScreen(tt)
			require.False(t, ok)
			require.Nil(t, path)
		})
	}
}
