package flow

import "strings"

// CatalogSelection stores one selected node inside catalog path.
type CatalogSelection struct {
	Level CatalogLevel
	ID    string
}

// CatalogPath describes currently selected catalog route.
//
// Example:
//
//	city=moscow;category=flowers;district=center
type CatalogPath []CatalogSelection

// Append returns a new path with one extra selection.
func (p CatalogPath) Append(level CatalogLevel, id string) CatalogPath {
	next := make(CatalogPath, 0, len(p)+1)
	next = append(next, p...)
	next = append(next, CatalogSelection{
		Level: level,
		ID:    id,
	})
	return next
}

// Last returns the last selected path item.
func (p CatalogPath) Last() (CatalogSelection, bool) {
	if len(p) == 0 {
		return CatalogSelection{}, false
	}
	return p[len(p)-1], true
}

// encodeCatalogPath serializes catalog path into compact screen payload form.
func encodeCatalogPath(path CatalogPath) string {
	if len(path) == 0 {
		return ""
	}

	parts := make([]string, 0, len(path))
	for _, sel := range path {
		if sel.Level == "" || strings.TrimSpace(sel.ID) == "" {
			return ""
		}
		parts = append(parts, string(sel.Level)+"="+sel.ID)
	}
	return strings.Join(parts, ";")
}

// parseCatalogPath parses serialized catalog path payload.
func parseCatalogPath(raw string) (CatalogPath, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, false
	}

	pairs := strings.Split(raw, ";")
	path := make(CatalogPath, 0, len(pairs))

	for _, pair := range pairs {
		parts := strings.Split(pair, "=")
		if len(parts) != 2 {
			return nil, false
		}
		level := strings.TrimSpace(parts[0])
		id := strings.TrimSpace(parts[1])
		if level == "" || id == "" {
			return nil, false
		}

		path = append(path, CatalogSelection{
			Level: CatalogLevel(level),
			ID:    id,
		})
	}

	return path, true
}
