package telegram

import (
	"fmt"
	"strings"
	"time"

	"github.com/koha90/shopcore/internal/flow"
)

// catalogImageUploadPath builds a relative catalog image path for uploaded admin photos.
//
// The path is stored in catalog image_url and later rendered by Telegram/Web adapters.
// Runtime owns the file system details; flow only receives the final path.
func catalogImageUploadPath(target flow.CatalogImageInputTarget, now time.Time) (string, error) {
	if target.EntityID <= 0 {
		return "", fmt.Errorf("catalog image target entity id is invalid")
	}

	stem := catalogImageFileStem(target)

	switch target.Kind {
	case flow.CatalogImageTargetProduct:
		return fmt.Sprintf(
			"assets/catalog/products/%s-%d.jpg",
			stem,
			now.Unix(),
		), nil

	case flow.CatalogImageTargetVariant:
		return fmt.Sprintf(
			"assets/catalog/variants/%s-%d.jpg",
			stem,
			now.Unix(),
		), nil

	default:
		return "", fmt.Errorf("unknown catalog image target kind %q", target.Kind)
	}
}

func catalogImageFileStem(target flow.CatalogImageInputTarget) string {
	code := sanitizeCatalogImageCode(target.EntityCode)
	if code == "" {
		return fmt.Sprintf("%d", target.EntityID)
	}

	return fmt.Sprintf("%d-%s", target.EntityID, code)
}

func sanitizeCatalogImageCode(v string) string {
	v = strings.TrimSpace(strings.ToLower(v))
	if v == "" {
		return ""
	}

	var b strings.Builder
	for _, r := range v {
		switch {
		case r >= 'a' && r <= 'z':
			b.WriteRune(r)
		case r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == '-' || r == '_':
			b.WriteRune(r)
		}
	}

	return b.String()
}
