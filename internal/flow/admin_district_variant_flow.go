package flow

import (
	"context"
	"errors"
	"strconv"
	"strings"
)

const (
	ActionAdminDistrictVariantCreateStart ActionID = "admin_district_variant_create_start"
)

func (s *Service) handleAdminDistrictVariantAction(
	ctx context.Context,
	session Session,
	req ActionRequest,
) (ViewModel, Session, bool, error) {
	switch req.ActionID {
	case ActionAdminDistrictVariantCreateStart:
		if s.cityLister == nil {
			return ViewModel{}, session, true, errors.New("flow city lister is nil")
		}

		next := ScreenAdminDistrictVariantCitySelect
		if next != session.Current {
			session.History = append(session.History, session.Current)
		}
		session.Current = next
		session.Pending = PendingInput{}

		return s.buildAdminDistrictVariantCitySelectScreen(), session, true, nil
	}

	if cityID, ok := parseAdminDistrictVariantSelectCityAction(req.ActionID); ok {
		if session.Current != ScreenAdminDistrictVariantCitySelect {
			return ViewModel{}, session, false, nil
		}
		if s.cityLister == nil {
			return ViewModel{}, session, true, errors.New("flow city lister is nil")
		}

		cities, err := s.cityLister.ListCities(ctx)
		if err != nil {
			return ViewModel{}, session, true, err
		}

		var selected *CityListItem
		for i := range cities {
			if cities[i].ID == cityID {
				selected = &cities[i]
				break
			}
		}
		if selected == nil {
			return ViewModel{}, session, true, ErrUnknownAction
		}

		next := ScreenAdminDistrictVariantDistrictSelect
		if next != session.Current {
			session.History = append(session.History, session.Current)
		}
		session.Current = next
		session.Pending = PendingInput{
			Kind: PendingInputNone,
			Payload: PendingInputPayload{
				PendingValueCityID:   strconv.Itoa(selected.ID),
				PendingValueCityName: selected.Label,
			},
		}

		return s.buildAdminDistrictVariantDistrictSelectScreen(selected.ID, selected.Label), session, true, nil
	}

	if districtID, ok := parseAdminDistrictVariantSelectDistrictAction(req.ActionID); ok {
		if session.Current != ScreenAdminDistrictVariantDistrictSelect {
			return ViewModel{}, session, false, nil
		}
		if s.districtLister == nil {
			return ViewModel{}, session, true, errors.New("flow district lister is nil")
		}
		if s.categoryLister == nil {
			return ViewModel{}, session, true, errors.New("flow category lister is nil")
		}

		cityID, ok := pendingCityID(session.Pending)
		if !ok {
			return ViewModel{}, session, true, errors.New("pending city id is invalid")
		}
		cityName := session.Pending.Value(PendingValueCityName)
		if cityName == "" {
			return ViewModel{}, session, true, errors.New("pending city name is empty")
		}

		districts, err := s.districtLister.ListDistrictsByCity(ctx, cityID)
		if err != nil {
			return ViewModel{}, session, true, err
		}

		var selected *DistrictListItem
		for i := range districts {
			if districts[i].ID == districtID {
				selected = &districts[i]
				break
			}
		}
		if selected == nil {
			return ViewModel{}, session, true, ErrUnknownAction
		}

		next := ScreenAdminDistrictVariantCategorySelect
		if next != session.Current {
			session.History = append(session.History, session.Current)
		}
		session.Current = next
		session.Pending = PendingInput{
			Kind: PendingInputNone,
			Payload: PendingInputPayload{
				PendingValueCityID:       strconv.Itoa(cityID),
				PendingValueCityName:     cityName,
				PendingValueDistrictID:   strconv.Itoa(selected.ID),
				PendingValueDistrictName: selected.Label,
			},
		}

		return s.buildAdminDistrictVariantCategorySelectScreen(cityName, selected.Label), session, true, nil
	}

	if categoryID, ok := parseAdminDistrictVariantSelectCategoryAction(req.ActionID); ok {
		if session.Current != ScreenAdminDistrictVariantCategorySelect {
			return ViewModel{}, session, false, nil
		}
		if s.categoryLister == nil {
			return ViewModel{}, session, true, errors.New("flow category lister is nil")
		}
		if s.productLister == nil {
			return ViewModel{}, session, true, errors.New("flow product lister is nil")
		}

		categories, err := s.categoryLister.ListCategories(ctx)
		if err != nil {
			return ViewModel{}, session, true, err
		}

		var selected *CategoryListItem
		for i := range categories {
			if categories[i].ID == categoryID {
				selected = &categories[i]
				break
			}
		}
		if selected == nil {
			return ViewModel{}, session, true, ErrUnknownAction
		}

		cityID := session.Pending.Value(PendingValueCityID)
		cityName := session.Pending.Value(PendingValueCityName)
		districtID := session.Pending.Value(PendingValueDistrictID)
		districtName := session.Pending.Value(PendingValueDistrictName)

		next := ScreenAdminDistrictVariantProductSelect
		if next != session.Current {
			session.History = append(session.History, session.Current)
		}
		session.Current = next
		session.Pending = PendingInput{
			Kind: PendingInputNone,
			Payload: PendingInputPayload{
				PendingValueCityID:       cityID,
				PendingValueCityName:     cityName,
				PendingValueDistrictID:   districtID,
				PendingValueDistrictName: districtName,
				PendingValueCategoryID:   strconv.Itoa(selected.ID),
				PendingValueCategoryName: selected.Label,
			},
		}

		return s.buildAdminDistrictVariantProductSelectScreen(
			cityName,
			districtName,
			selected.Label,
		), session, true, nil
	}

	if productID, ok := parseAdminDistrictVariantSelectProductAction(req.ActionID); ok {
		if session.Current != ScreenAdminDistrictVariantProductSelect {
			return ViewModel{}, session, false, nil
		}
		if s.productLister == nil {
			return ViewModel{}, session, true, errors.New("flow product lister is nil")
		}
		if s.variantLister == nil {
			return ViewModel{}, session, true, errors.New("flow variant lister is nil")
		}

		products, err := s.productLister.ListProducts(ctx)
		if err != nil {
			return ViewModel{}, session, true, err
		}

		var selected *ProductListItem
		for i := range products {
			if products[i].ID == productID {
				selected = &products[i]
				break
			}
		}
		if selected == nil {
			return ViewModel{}, session, true, ErrUnknownAction
		}

		cityID := session.Pending.Value(PendingValueCityID)
		cityName := session.Pending.Value(PendingValueCityName)
		districtID := session.Pending.Value(PendingValueDistrictID)
		districtName := session.Pending.Value(PendingValueDistrictName)
		categoryID := session.Pending.Value(PendingValueCategoryID)
		categoryName := session.Pending.Value(PendingValueCategoryName)

		next := ScreenAdminDistrictVariantVariantSelect
		if next != session.Current {
			session.History = append(session.History, session.Current)
		}
		session.Current = next
		session.Pending = PendingInput{
			Kind: PendingInputNone,
			Payload: PendingInputPayload{
				PendingValueCityID:       cityID,
				PendingValueCityName:     cityName,
				PendingValueDistrictID:   districtID,
				PendingValueDistrictName: districtName,
				PendingValueCategoryID:   categoryID,
				PendingValueCategoryName: categoryName,
				PendingValueProductID:    strconv.Itoa(selected.ID),
				PendingValueProductName:  selected.Label,
			},
		}

		return s.buildAdminDistrictVariantVariantSelectScreen(
			cityName,
			districtName,
			categoryName,
			selected.ID,
			selected.Label,
		), session, true, nil
	}

	if variantID, ok := parseAdminDistrictVariantSelectVariantAction(req.ActionID); ok {
		if session.Current != ScreenAdminDistrictVariantVariantSelect {
			return ViewModel{}, session, false, nil
		}
		if s.variantLister == nil {
			return ViewModel{}, session, true, errors.New("flow variant lister is nil")
		}

		districtID, ok := pendingDistrictID(session.Pending)
		if !ok {
			return ViewModel{}, session, true, errors.New("pending district id is invalid")
		}
		districtName := session.Pending.Value(PendingValueDistrictName)
		if districtName == "" {
			return ViewModel{}, session, true, errors.New("pending district name is empty")
		}

		productID, ok := pendingProductID(session.Pending)
		if !ok {
			return ViewModel{}, session, true, errors.New("pending product id is invalid")
		}

		variants, err := s.variantLister.ListVariantsByProduct(ctx, productID)
		if err != nil {
			return ViewModel{}, session, true, err
		}

		var selected *VariantListItem
		for i := range variants {
			if variants[i].ID == variantID {
				selected = &variants[i]
				break
			}
		}
		if selected == nil {
			return ViewModel{}, session, true, ErrUnknownAction
		}

		variantDisplayLabel := buildAdminVariantOptionLabel(*selected)

		next := ScreenAdminDistrictVariantPrice
		if next != session.Current {
			session.History = append(session.History, session.Current)
		}
		session.Current = next
		session.Pending = PendingInput{
			Kind: PendingInputDistrictVariantPrice,
			Payload: PendingInputPayload{
				PendingValueDistrictID:   strconv.Itoa(districtID),
				PendingValueDistrictName: districtName,
				PendingValueVariantID:    strconv.Itoa(selected.ID),
				PendingValueVariantName:  variantDisplayLabel,
			},
		}

		return buildAdminDistrictVariantPriceInputView(
			districtName,
			variantDisplayLabel,
			"",
		), session, true, nil
	}

	return ViewModel{}, session, false, nil
}

func parseAdminDistrictVariantSelectCityAction(actionID ActionID) (int, bool) {
	const prefix = "admin:district_variant:city:"

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

func parseAdminDistrictVariantSelectCategoryAction(actionID ActionID) (int, bool) {
	const prefix = "admin:district_variant:category:"

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

func parseAdminDistrictVariantSelectProductAction(actionID ActionID) (int, bool) {
	const prefix = "admin:district_variant:product:"

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
