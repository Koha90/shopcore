package flow

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	ActionAdminDistrictVariantCreateStart ActionID = "admin_district_variant_create_start"
)

func adminDistrictVariantSelectDistrictAction(districtID int) ActionID {
	return ActionID(fmt.Sprintf("admin:district_variant:district:%d", districtID))
}

func parseAdminDistrictVariantSelectDistrictAction(actionID ActionID) (int, bool) {
	const prefix = "admin:district_variant:district:"

	raw := strings.TrimSpace(string(actionID))
	if !strings.HasPrefix(raw, prefix) {
		return 0, false
	}

	idPart := strings.TrimPrefix(raw, prefix)
	if idPart == "" {
		return 0, false
	}

	id, err := strconv.Atoi(idPart)
	if err != nil || id <= 0 {
		return 0, false
	}

	return id, true
}

func adminDistrictVariantSelectVariantAction(variantID int) ActionID {
	return ActionID(fmt.Sprintf("admin:district_variant:variant:%d", variantID))
}

func parseAdminDistrictVariantSelectVariantAction(actionID ActionID) (int, bool) {
	const prefix = "admin:district_variant:variant:"

	raw := strings.TrimSpace(string(actionID))
	if !strings.HasPrefix(raw, prefix) {
		return 0, false
	}

	idPart := strings.TrimPrefix(raw, prefix)
	if idPart == "" {
		return 0, false
	}

	id, err := strconv.Atoi(idPart)
	if err != nil || id <= 0 {
		return 0, false
	}

	return id, true
}

func buildAdminQualifiedVariantOptionLabel(productLabel, variantLabel string) string {
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
	return buildAdminQualifiedVariantOptionLabel(item.ProductLabel, item.Label)
}
