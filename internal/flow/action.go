package flow

// ScreenID identifies current logical screen in flow.
//
// Root and detail screen use stable named identiriers.
// Catalog drill-down screen are encoded dynamically from CatalogPath.
type ScreenID string

const (
	ScreenReplyWelcome ScreenID = "reply_welcome"
	ScreenRootCompact  ScreenID = "root_compact"
	ScreenRootExtended ScreenID = "root_extended"

	ScreenCabinet   ScreenID = "cabinet"
	ScreenSupport   ScreenID = "support"
	ScreenReviews   ScreenID = "reviews"
	ScreenBalance   ScreenID = "balance"
	ScreenBotsMine  ScreenID = "bots_mine"
	ScreenOrderLast ScreenID = "order_last"

	ScreenAdminRoot               ScreenID = "admin_root"
	ScreenAdminCatalog            ScreenID = "admin_catalog"
	ScreenAdminCategoryCreate     ScreenID = "admin_category_create"
	ScreenAdminCategoryCode       ScreenID = "admin_category_code"
	ScreenAdminCategoryCreateDone ScreenID = "admin_category_create_done"

	ScreenAdminCityCreate     ScreenID = "admin_city_create"
	ScreenAdminCityCode       ScreenID = "admin_city_code"
	ScreenAdminCityCreateDone ScreenID = "admin_city_create_done"

	ScreenAdminDistrictCitySelect ScreenID = "admin_district_city_select"
	ScreenAdminDistrictCreate     ScreenID = "admin_district_create"
	ScreenAdminDistrictCode       ScreenID = "admin_district_code"
	ScreenAdminDistrictCreateDone ScreenID = "admin_district_create_done"

	ScreenAdminProductCategorySelect ScreenID = "admin_product_category_select"
	ScreenAdminProductCreate         ScreenID = "admin_product_create"
	ScreenAdminProductCode           ScreenID = "admin_product_code"
	ScreenAdminProductCreateDone     ScreenID = "admin_product_create_done"

	ScreenAdminVariantProductSelect ScreenID = "admin_variant_product_select"
	ScreenAdminVariantCreate        ScreenID = "admin_variant_create"
	ScreenAdminVariantCode          ScreenID = "admin_variant_code"
	ScreenAdminVariantCreateDone    ScreenID = "admin_variant_create_done"
)

// PendingInputKind identifies which text input flow currently expects.
type PendingInputKind string

const (
	PendingInputNone PendingInputKind = ""

	PendingInputCategoryName PendingInputKind = "category_name"
	PendingInputCategoryCode PendingInputKind = "category_code"

	PendingInputCityName PendingInputKind = "city_name"
	PendingInputCityCode PendingInputKind = "city_code"

	PendingInputDistrictName PendingInputKind = "district_name"
	PendingInputDistrictCode PendingInputKind = "district_code"

	PendingInputProductName PendingInputKind = "product_name"
	PendingInputProductCode PendingInputKind = "product_code"

	PendingInputVariantName PendingInputKind = "variant_name"
	PendingInputVariantCode PendingInputKind = "variant_code"
)

const (
	// PendingValueName stores one entered name value inside pending input payload.
	PendingValueName = "name"

	// PendingValueCode stores one entered code value inside pending input payload.
	PendingValueCode = "code"

	// PendingValueCityID stores selected city id inside pending input payload.
	PendingValueCityID = "city_id"

	// PendingValueCityName = "city_name"
	PendingValueCityName = "city_name"

	// PendingValueCategoryID stores selected category id inside pending input payload.
	PendingValueCategoryID = "category_id"

	// PendingValueCategoryName stores selected category name label inside pending input payload.
	PendingValueCategoryName = "category_name"

	// PendingValueProductID stores selected product id inside pending input payload.
	PendingValueProductID = "product_id"

	// PendingValueProductName stores selected product label inside pending input payload.
	PendingValueProductName = "product_name"
)
