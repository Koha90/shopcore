package flow

func (s *Service) renderAdminCatalogCreateScreen(session Session) (ViewModel, bool) {
	switch session.Current {
	case ScreenAdminCategoryCreate:
		return buildAdminCategoryCreateInputView(""), true

	case ScreenAdminCategoryCode:
		return buildAdminCategoryCodeInputView("", ""), true

	case ScreenAdminCategoryCreateDone:
		return buildAdminCategoryCreateDoneView(), true

	case ScreenAdminCityCreate:
		return buildAdminCityCreateInputView(""), true

	case ScreenAdminCityCode:
		return buildAdminCityCodeInputView("", ""), true

	case ScreenAdminCityCreateDone:
		return buildAdminCityCreateDoneView(), true

	case ScreenAdminDistrictCitySelect:
		return s.buildAdminDistrictCitySelectScreen(), true

	case ScreenAdminDistrictCreate:
		return buildAdminDistrictCreateInputView("", ""), true

	case ScreenAdminDistrictCode:
		return buildAdminDistrictCodeInputView("", "", ""), true

	case ScreenAdminDistrictCreateDone:
		return buildAdminDistrictCreateDoneView(), true

	case ScreenAdminProductCategorySelect:
		return s.buildAdminProductCategorySelectScreen(), true

	case ScreenAdminProductCreate:
		return buildAdminProductCreateInputView("", ""), true

	case ScreenAdminProductCode:
		return buildAdminProductCodeInputView("", "", ""), true

	case ScreenAdminProductCreateDone:
		return buildAdminProductCreateDoneView(), true

	case ScreenAdminVariantProductSelect:
		return s.buildAdminVariantProductSelectScreen(), true

	case ScreenAdminVariantCreate:
		return buildAdminVariantCreateInputView("", ""), true

	case ScreenAdminVariantCode:
		return buildAdminVariantCodeInputView("", "", ""), true

	case ScreenAdminVariantCreateDone:
		return buildAdminVariantCreateDoneView(), true

	default:
		return ViewModel{}, false
	}
}
