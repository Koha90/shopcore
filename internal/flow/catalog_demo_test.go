package flow

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDemoCatalog_RootLevel(t *testing.T) {
	catalog := DemoCatalog()

	level, ok := catalog.RootLevel()
	require.True(t, ok)
	require.Equal(t, LevelCity, level)
}

func TestDemoCatalog_RootNodes(t *testing.T) {
	catalog := DemoCatalog()

	roots := catalog.RootNodes()
	require.NotEmpty(t, roots)
	require.Equal(t, LevelCity, roots[0].Level)
}

func TestDemoCatalog_FindNode_City(t *testing.T) {
	catalog := DemoCatalog()

	node, ok := catalog.FindNode(CatalogPath{
		{Level: LevelCity, ID: "moscow"},
	})
	require.True(t, ok)
	require.Equal(t, "Москва", node.Label)
	require.Equal(t, LevelCity, node.Level)
}

func TestDemoCatalog_FindNode_DeepPath(t *testing.T) {
	catalog := DemoCatalog()

	node, ok := catalog.FindNode(CatalogPath{
		{Level: LevelCity, ID: "moscow"},
		{Level: LevelCategory, ID: "flowers"},
		{Level: LevelDistrict, ID: "center"},
		{Level: LevelProduct, ID: "rose-box"},
		{Level: LevelVariant, ID: "large"},
	})
	require.True(t, ok)
	require.Equal(t, LevelVariant, node.Level)
	require.Equal(t, "L / 25 шт", node.Label)
	require.Equal(t, "5900 ₽", node.PriceText)
}

func TestDemoCatalog_FindNode_InvalidPath(t *testing.T) {
	catalog := DemoCatalog()

	_, ok := catalog.FindNode(CatalogPath{
		{Level: LevelCity, ID: "moscow"},
		{Level: LevelCategory, ID: "unknown"},
	})
	require.False(t, ok)
}

func TestDemoCatalog_FindNode_EmptyPath(t *testing.T) {
	catalog := DemoCatalog()

	_, ok := catalog.FindNode(nil)
	require.False(t, ok)
}
