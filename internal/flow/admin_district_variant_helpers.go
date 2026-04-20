package flow

import "strings"

func buildAdminQualifiedVariantLabel(productLabel, variantLabel string) string {
	productLabel = strings.TrimSpace(productLabel)
	variantLabel = strings.TrimSpace(variantLabel)

	switch {
	case productLabel == "":
		return variantLabel
	case variantLabel == "":
		return productLabel
	default:
		return productLabel + " - " + variantLabel
	}
}

func buildAdminVariantOptionLabel(item VariantListItem) string {
	return buildAdminQualifiedVariantLabel(item.ProductLabel, item.Label)
}
