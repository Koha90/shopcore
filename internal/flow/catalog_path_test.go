package flow

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCatalogPath_AppendAndLast(t *testing.T) {
	var path CatalogPath

	path = path.Append(LevelCity, "moscow")
	path = path.Append(LevelCategory, "flowers")

	last, ok := path.Last()
	require.True(t, ok)
	require.Equal(t, CatalogSelection{
		Level: LevelCategory,
		ID:    "flowers",
	}, last)

	require.Len(t, path, 2)
	require.Equal(t, CatalogSelection{
		Level: LevelCity,
		ID:    "moscow",
	}, path[0])
}

func TestCatalogPath_Last_Empty(t *testing.T) {
	var path CatalogPath

	last, ok := path.Last()
	require.False(t, ok)
	require.Equal(t, CatalogSelection{}, last)
}

func TestEncodeParseCatalogPath_RoundTrip(t *testing.T) {
	path := CatalogPath{
		{Level: LevelCity, ID: "moscow"},
		{Level: LevelCategory, ID: "flowers"},
		{Level: LevelDistrict, ID: "center"},
		{Level: LevelProduct, ID: "rose-box"},
		{Level: LevelVariant, ID: "large"},
	}

	raw := encodeCatalogPath(path)
	require.Equal(t, "city=moscow;category=flowers;district=center;product=rose-box;variant=large", raw)

	parsed, ok := parseCatalogPath(raw)
	require.True(t, ok)
	require.Equal(t, path, parsed)
}

func TestEncodeCatalogPath_InvalidSelection(t *testing.T) {
	raw := encodeCatalogPath(CatalogPath{
		{Level: LevelCity, ID: "moscow"},
		{Level: "", ID: "broken"},
	})

	require.Equal(t, "", raw)
}

func TestParseCatalogPath_Invalid(t *testing.T) {
	tests := []string{
		"",
		"city",
		"city=",
		"=moscow",
		"city=moscow;broken",
		"city=moscow;category=",
	}

	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			path, ok := parseCatalogPath(tt)
			require.False(t, ok)
			require.Nil(t, path)
		})
	}
}
