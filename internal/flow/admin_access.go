package flow

func isAdminAction(actionID ActionID) bool {
	if _, ok := parseAdminDistrictSelectCityAction(actionID); ok {
		return true
	}
	if _, ok := parseAdminProductSelectCategoryAction(actionID); ok {
		return true
	}
	if _, ok := parseAdminVariantSelectProductAction(actionID); ok {
		return true
	}
	if _, ok := parseAdminDistrictVariantSelectDistrictAction(actionID); ok {
		return true
	}
	if _, ok := parseAdminDistrictVariantPriceUpdateSelectCategoryAction(actionID); ok {
		return true
	}
	if _, ok := parseAdminDistrictVariantPriceUpdateSelectProductAction(actionID); ok {
		return true
	}
	if _, ok := parseAdminDistrictVariantSelectVariantAction(actionID); ok {
		return true
	}
	if _, ok := parseAdminDistrictVariantSelectCityAction(actionID); ok {
		return true
	}
	if _, ok := parseAdminDistrictVariantSelectCategoryAction(actionID); ok {
		return true
	}
	if _, ok := parseAdminDistrictVariantSelectProductAction(actionID); ok {
		return true
	}
	if _, _, ok := parseAdminCustomerReplyStartAction(actionID); ok {
		return true
	}
	if _, _, ok := parseAdminCustomerPhotoReplyStartAction(actionID); ok {
		return true
	}
	if _, ok := parseAdminProductImageSelectProductAction(actionID); ok {
		return true
	}
	if _, ok := parseAdminVariantImageSelectVariantAction(actionID); ok {
		return true
	}

	switch actionID {
	case ActionAdminOpen,
		ActionAdminCatalogOpen,
		ActionAdminCategoryCreateStart,
		ActionAdminCityCreateStart,
		ActionAdminDistrictCreateStart,
		ActionAdminProductCreateStart,
		ActionAdminVariantCreateStart,
		ActionAdminDistrictVariantCreateStart,
		ActionAdminDistrictVariantPriceUpdateStart,
		ActionAdminProductImageUpdateStart,
		ActionAdminVariantImageUpdateStart:
		return true
	default:
		return false
	}
}

func isAdminScreen(screen ScreenID) bool {
	switch screen {
	case ScreenAdminRoot,
		ScreenAdminCatalog,
		ScreenAdminCategoryCreate,
		ScreenAdminCategoryCode,
		ScreenAdminCategoryCreateDone,
		ScreenAdminCityCreate,
		ScreenAdminCityCode,
		ScreenAdminCityCreateDone,
		ScreenAdminDistrictCitySelect,
		ScreenAdminDistrictCreate,
		ScreenAdminDistrictCode,
		ScreenAdminDistrictCreateDone,
		ScreenAdminProductCategorySelect,
		ScreenAdminProductCreate,
		ScreenAdminProductCode,
		ScreenAdminProductCreateDone,
		ScreenAdminVariantProductSelect,
		ScreenAdminVariantCreate,
		ScreenAdminVariantCode,
		ScreenAdminVariantCreateDone,
		ScreenAdminDistrictVariantCitySelect,
		ScreenAdminDistrictVariantDistrictSelect,
		ScreenAdminDistrictVariantCategorySelect,
		ScreenAdminDistrictVariantProductSelect,
		ScreenAdminDistrictVariantVariantSelect,
		ScreenAdminDistrictVariantPrice,
		ScreenAdminDistrictVariantCreateDone,
		ScreenAdminDistrictVariantPriceUpdateDistrictSelect,
		ScreenAdminDistrictVariantPriceUpdateCategorySelect,
		ScreenAdminDistrictVariantPriceUpdateProductSelect,
		ScreenAdminDistrictVariantPriceUpdateVariantSelect,
		ScreenAdminDistrictVariantPriceUpdatePrice,
		ScreenAdminDistrictVariantPriceUpdateDone,
		ScreenAdminCustomerReply,
		ScreenAdminCustomerPhotoReply,
		ScreenAdminCustomerReplyDone,
		ScreenAdminProductImageProductSelect,
		ScreenAdminProductImageInput,
		ScreenAdminProductImageDone,
		ScreenAdminVariantImageVariantSelect,
		ScreenAdminVariantImageInput,
		ScreenAdminVariantImageDone:
		return true
	default:
		return false
	}
}

func isAdminPending(kind PendingInputKind) bool {
	switch kind {
	case PendingInputCategoryName,
		PendingInputCategoryCode,
		PendingInputCityName,
		PendingInputCityCode,
		PendingInputDistrictName,
		PendingInputDistrictCode,
		PendingInputProductName,
		PendingInputProductCode,
		PendingInputVariantName,
		PendingInputVariantCode,
		PendingInputDistrictVariantPrice,
		PendingInputDistrictVariantPriceUpdate,
		PendingInputAdminCustomerReply,
		PendingInputAdminCustomerPhotoReply,
		PendingInputProductImageURL,
		PendingInputVariantImageURL:
		return true
	default:
		return false
	}
}
