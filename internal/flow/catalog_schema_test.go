package flow

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDemoCatalogSchema_First(t *testing.T) {
	schema := DemoCatalogSchema()

	level, ok := schema.First()
	require.True(t, ok)
	require.Equal(t, LevelCity, level)
}

func TestCatalogSchema_First_Empty(t *testing.T) {
	var schema CatalogSchema

	level, ok := schema.First()
	require.False(t, ok)
	require.Empty(t, level)
}

func TestDemoCatalogSchema_Next(t *testing.T) {
	schema := DemoCatalogSchema()

	next, ok := schema.Next(LevelCity)
	require.True(t, ok)
	require.Equal(t, LevelCategory, next)

	next, ok = schema.Next(LevelCategory)
	require.True(t, ok)
	require.Equal(t, LevelDistrict, next)

	next, ok = schema.Next(LevelDistrict)
	require.True(t, ok)
	require.Equal(t, LevelProduct, next)

	next, ok = schema.Next(LevelProduct)
	require.True(t, ok)
	require.Equal(t, LevelVariant, next)

	next, ok = schema.Next(LevelVariant)
	require.False(t, ok)
	require.Empty(t, next)
}

func TestCatalogSchema_Next_UnknownLevel(t *testing.T) {
	schema := DemoCatalogSchema()

	next, ok := schema.Next(CatalogLevel("unknown"))
	require.False(t, ok)
	require.Empty(t, next)
}
