package flow

func (s *Service) renderAdminDistrictVariantScreen(session Session) (ViewModel, bool) {
	switch session.Current {
	case ScreenAdminDistrictVariantCitySelect:
		return s.buildAdminDistrictVariantCitySelectScreen(), true

	case ScreenAdminDistrictVariantDistrictSelect:
		cityID, ok := pendingCityID(session.Pending)
		if !ok {
			return buildAdminCatalogView(), true
		}
		cityName := session.Pending.Value(PendingValueCityName)
		return s.buildAdminDistrictVariantDistrictSelectScreen(cityID, cityName), true

	case ScreenAdminDistrictVariantCategorySelect:
		cityName := session.Pending.Value(PendingValueCityName)
		districtName := session.Pending.Value(PendingValueDistrictName)
		return s.buildAdminDistrictVariantCategorySelectScreen(cityName, districtName), true

	case ScreenAdminDistrictVariantProductSelect:
		categoryID, ok := pendingCategoryID(session.Pending)
		if !ok {
			return buildAdminCatalogView(), true
		}

		cityName := session.Pending.Value(PendingValueCityName)
		districtName := session.Pending.Value(PendingValueDistrictName)
		categoryName := session.Pending.Value(PendingValueCategoryName)

		return s.buildAdminDistrictVariantProductSelectScreen(
			cityName,
			districtName,
			categoryID,
			categoryName,
		), true

	case ScreenAdminDistrictVariantVariantSelect:
		productID, ok := pendingProductID(session.Pending)
		if !ok {
			return buildAdminCatalogView(), true
		}

		cityName := session.Pending.Value(PendingValueCityName)
		districtName := session.Pending.Value(PendingValueDistrictName)
		categoryName := session.Pending.Value(PendingValueCategoryName)
		productName := session.Pending.Value(PendingValueProductName)

		return s.buildAdminDistrictVariantVariantSelectScreen(
			cityName,
			districtName,
			categoryName,
			productID,
			productName,
		), true

	case ScreenAdminDistrictVariantPrice:
		return buildAdminDistrictVariantPriceInputView("", "", ""), true

	case ScreenAdminDistrictVariantCreateDone:
		return buildAdminDistrictVariantCreateDoneView(), true

	default:
		return ViewModel{}, false
	}
}

func (s *Service) renderAdminDistrictVariantPriceUpdateScreen(session Session) (ViewModel, bool) {
	switch session.Current {
	case ScreenAdminDistrictVariantPriceUpdateDistrictSelect:
		return s.buildAdminDistrictVariantPriceUpdateDistrictSelectScreen(), true

	case ScreenAdminDistrictVariantPriceUpdateCategorySelect:
		districtID, ok := pendingDistrictID(session.Pending)
		if !ok {
			return buildAdminCatalogView(), true
		}
		districtName := session.Pending.Value(PendingValueDistrictName)
		return s.buildAdminDistrictVariantPriceUpdateCategorySelectScreen(
			districtID,
			districtName,
		), true

	case ScreenAdminDistrictVariantPriceUpdateProductSelect:
		districtID, ok := pendingDistrictID(session.Pending)
		if !ok {
			return buildAdminCatalogView(), true
		}
		categoryID, ok := pendingCategoryID(session.Pending)
		if !ok {
			return buildAdminCatalogView(), true
		}
		districtName := session.Pending.Value(PendingValueDistrictName)
		categoryName := session.Pending.Value(PendingValueCategoryName)

		return s.buildAdminDistrictVariantPriceUpdateProductSelectScreen(
			districtID,
			districtName,
			categoryID,
			categoryName,
		), true

	case ScreenAdminDistrictVariantPriceUpdateVariantSelect:
		districtID, ok := pendingDistrictID(session.Pending)
		if !ok {
			return buildAdminCatalogView(), true
		}
		productID, ok := pendingProductID(session.Pending)
		if !ok {
			return buildAdminCatalogView(), true
		}
		districtName := session.Pending.Value(PendingValueDistrictName)
		productName := session.Pending.Value(PendingValueProductName)

		return s.buildAdminDistrictVariantPriceUpdateVariantSelectScreen(
			districtID,
			districtName,
			productID,
			productName,
		), true

	case ScreenAdminDistrictVariantPriceUpdatePrice:
		return buildAdminDistrictVariantPriceUpdateInputView(
			session.Pending.Value(PendingValueDistrictName),
			session.Pending.Value(PendingValueVariantName),
			currentPlacementPriceTextFromPending(session.Pending),
			"",
		), true

	case ScreenAdminDistrictVariantPriceUpdateDone:
		return buildAdminDistrictVariantPriceUpdateDoneView(), true

	default:
		return ViewModel{}, false
	}
}
