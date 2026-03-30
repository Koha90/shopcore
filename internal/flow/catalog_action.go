package flow

import "strings"

const (
	catalogSelectPrefix = "catalog:select:"
	catalogScreenPrefix = "catalog:screen:"
)

// catalogSelectAction builds generic catalog selection action.
//
// Examples:
//
//	catalog:select:city:moscow
//	catalog:select:product:sku-001
func catalogSelectAction(level CatalogLevel, id string) ActionID {
	return ActionID(catalogSelectPrefix + string(level) + ":" + id)
}

// parseCatalogSelectAction parses generic catalog selection action payload.
func parseCatalogSelectAction(id ActionID) (CatalogLevel, string, bool) {
	raw := strings.TrimPrefix(string(id), catalogSelectPrefix)
	if raw == string(id) || raw == "" {
		return "", "", false
	}

	parts := strings.Split(raw, ":")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", false
	}
	return CatalogLevel(parts[0]), parts[1], true
}

// catalogScreen builds screen identifier for catalog path state.
func catalogScreen(path CatalogPath) ScreenID {
	return ScreenID(catalogScreenPrefix + encodeCatalogPath(path))
}

// parseCatalogScreen parses catalog screen payload back into path state.
func parseCatalogScreen(id ScreenID) (CatalogPath, bool) {
	raw := strings.TrimPrefix(string(id), catalogScreenPrefix)
	if raw == string(id) || raw == "" {
		return nil, false
	}
	return parseCatalogPath(raw)
}
