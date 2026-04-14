package flow

import (
	"strconv"
	"strings"
)

const (
	adminDistrictVariantPriceUpdateCategoryPrefix = "admin:district_variant:price_update:category"
	adminDistrictVariantPriceUpdateProductPrefix  = "admin:district_variant:price_update:product"
)

func adminDistrictVariantPriceUpdateSelectCategoryAction(categoryID int) ActionID {
	return ActionID(adminDistrictVariantPriceUpdateCategoryPrefix + strconv.Itoa(categoryID))
}

func parseAdminDistrictVariantPriceUpdateSelectCategoryAction(id ActionID) (int, bool) {
	raw := strings.TrimPrefix(string(id), adminDistrictVariantPriceUpdateCategoryPrefix)
	if raw == string(id) || raw == "" {
		return 0, false
	}

	v, err := strconv.Atoi(raw)
	if err != nil || v <= 0 {
		return 0, false
	}

	return v, true
}

func adminDistrictVariantPriceUpdateSelectProductAction(productID int) ActionID {
	return ActionID(adminDistrictVariantPriceUpdateProductPrefix + strconv.Itoa(productID))
}

func parseAdminDistrictVariantPriceUpdateSelectProductAction(id ActionID) (int, bool) {
	raw := strings.TrimPrefix(string(id), adminDistrictVariantPriceUpdateProductPrefix)
	if raw == string(id) || raw == "" {
		return 0, false
	}

	v, err := strconv.Atoi(raw)
	if err != nil || v <= 0 {
		return 0, false
	}

	return v, true
}
