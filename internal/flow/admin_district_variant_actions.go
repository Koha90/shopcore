package flow

import (
	"fmt"
)

func adminDistrictVariantSelectCityAction(cityID int) ActionID {
	return ActionID(fmt.Sprintf("admin:district_variant:city:%d", cityID))
}

func adminDistrictVariantSelectDistrictAction(districtID int) ActionID {
	return ActionID(fmt.Sprintf("admin:district_variant:district:%d", districtID))
}

func adminDistrictVariantSelectProductAction(productID int) ActionID {
	return ActionID(fmt.Sprintf("admin:district_variant:product:%d", productID))
}

func adminDistrictVariantSelectCategoryAction(categoryID int) ActionID {
	return ActionID(fmt.Sprintf("admin:district_variant:category:%d", categoryID))
}

func adminDistrictVariantSelectVariantAction(variantID int) ActionID {
	return ActionID(fmt.Sprintf("admin:district_variant:variant:%d", variantID))
}
