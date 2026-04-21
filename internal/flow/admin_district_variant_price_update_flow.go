package flow

import (
	"context"
	"errors"
	"strconv"
)

func (s *Service) handleAdminDistrictVariantPriceUpdateAction(
	ctx context.Context,
	session Session,
	req ActionRequest,
) (ViewModel, Session, bool, error) {
	switch req.ActionID {
	case ActionAdminDistrictVariantPriceUpdateStart:
		if s.districtLister == nil {
			return ViewModel{}, session, true, errors.New("flow district lister is nil")
		}

		next := ScreenAdminDistrictVariantPriceUpdateDistrictSelect
		if next != session.Current {
			session.History = append(session.History, session.Current)
		}
		session.Current = next
		session.Pending = PendingInput{}

		return s.buildAdminDistrictVariantPriceUpdateDistrictSelectScreen(), session, true, nil
	}

	if categoryID, ok := parseAdminDistrictVariantPriceUpdateSelectCategoryAction(req.ActionID); ok {
		if session.Current != ScreenAdminDistrictVariantPriceUpdateCategorySelect {
			return ViewModel{}, session, false, nil
		}
		if s.districtPlacements == nil {
			return ViewModel{}, session, true, errors.New("flow district placements reader is nil")
		}

		districtID, ok := pendingDistrictID(session.Pending)
		if !ok {
			return ViewModel{}, session, true, errors.New("pending district id is invalid")
		}
		districtName := session.Pending.Value(PendingValueDistrictName)
		if districtName == "" {
			return ViewModel{}, session, true, errors.New("pending district name is empty")
		}

		categories, err := s.districtPlacements.ListDistrictCategories(ctx, districtID)
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

		next := ScreenAdminDistrictVariantPriceUpdateProductSelect
		if next != session.Current {
			session.History = append(session.History, session.Current)
		}
		session.Current = next
		session.Pending = PendingInput{
			Kind: PendingInputNone,
			Payload: PendingInputPayload{
				PendingValueDistrictID:   strconv.Itoa(districtID),
				PendingValueDistrictName: districtName,
				PendingValueCategoryID:   strconv.Itoa(selected.ID),
				PendingValueCategoryName: selected.Label,
			},
		}

		return s.buildAdminDistrictVariantPriceUpdateProductSelectScreen(
			districtID,
			districtName,
			selected.ID,
			selected.Label,
		), session, true, nil
	}

	if productID, ok := parseAdminDistrictVariantPriceUpdateSelectProductAction(req.ActionID); ok {
		if session.Current != ScreenAdminDistrictVariantPriceUpdateProductSelect {
			return ViewModel{}, session, false, nil
		}
		if s.districtPlacements == nil {
			return ViewModel{}, session, true, errors.New("flow district placements reader is nil")
		}

		districtID, ok := pendingDistrictID(session.Pending)
		if !ok {
			return ViewModel{}, session, true, errors.New("pending district id is invalid")
		}
		districtName := session.Pending.Value(PendingValueDistrictName)
		if districtName == "" {
			return ViewModel{}, session, true, errors.New("pending district name is empty")
		}

		categoryID, ok := pendingCategoryID(session.Pending)
		if !ok {
			return ViewModel{}, session, true, errors.New("pending category id is invalid")
		}
		categoryName := session.Pending.Value(PendingValueCategoryName)
		if categoryName == "" {
			return ViewModel{}, session, true, errors.New("pending category name is empty")
		}

		products, err := s.districtPlacements.ListDistrictProducts(ctx, districtID, categoryID)
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

		next := ScreenAdminDistrictVariantPriceUpdateVariantSelect
		if next != session.Current {
			session.History = append(session.History, session.Current)
		}
		session.Current = next
		session.Pending = PendingInput{
			Kind: PendingInputNone,
			Payload: PendingInputPayload{
				PendingValueDistrictID:   strconv.Itoa(districtID),
				PendingValueDistrictName: districtName,
				PendingValueCategoryID:   strconv.Itoa(categoryID),
				PendingValueCategoryName: categoryName,
				PendingValueProductID:    strconv.Itoa(selected.ID),
				PendingValueProductName:  selected.Label,
			},
		}

		return s.buildAdminDistrictVariantPriceUpdateVariantSelectScreen(
			districtID,
			districtName,
			selected.ID,
			selected.Label,
		), session, true, nil
	}

	if districtID, ok := parseAdminDistrictVariantSelectDistrictAction(req.ActionID); ok {
		if session.Current != ScreenAdminDistrictVariantPriceUpdateDistrictSelect {
			return ViewModel{}, session, false, nil
		}
		if s.districtLister == nil {
			return ViewModel{}, session, true, errors.New("flow district lister is nil")
		}

		districts, err := s.districtLister.ListDistricts(ctx)
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

		next := ScreenAdminDistrictVariantPriceUpdateCategorySelect
		if next != session.Current {
			session.History = append(session.History, session.Current)
		}
		session.Current = next
		session.Pending = PendingInput{
			Kind: PendingInputNone,
			Payload: PendingInputPayload{
				PendingValueDistrictID:   strconv.Itoa(selected.ID),
				PendingValueDistrictName: selected.Label,
			},
		}

		return s.buildAdminDistrictVariantPriceUpdateCategorySelectScreen(
			selected.ID,
			selected.Label,
		), session, true, nil
	}

	if variantID, ok := parseAdminDistrictVariantSelectVariantAction(req.ActionID); ok {
		if session.Current != ScreenAdminDistrictVariantPriceUpdateVariantSelect {
			return ViewModel{}, session, false, nil
		}
		if s.districtPlacements == nil {
			return ViewModel{}, session, true, errors.New("flow district placements reader is nil")
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

		productName := session.Pending.Value(PendingValueProductName)
		if productName == "" {
			return ViewModel{}, session, true, errors.New("pending product name is empty")
		}

		variants, err := s.districtPlacements.ListDistrictVariants(ctx, districtID, productID)
		if err != nil {
			return ViewModel{}, session, true, err
		}

		var selected *DistrictPlacementVariantListItem
		for i := range variants {
			if variants[i].ID == variantID {
				selected = &variants[i]
				break
			}
		}
		if selected == nil {
			return ViewModel{}, session, true, ErrUnknownAction
		}

		variantDisplayLabel := buildAdminQualifiedVariantLabel(productName, selected.Label)

		next := ScreenAdminDistrictVariantPriceUpdatePrice
		if next != session.Current {
			session.History = append(session.History, session.Current)
		}
		session.Current = next
		session.Pending = PendingInput{
			Kind: PendingInputDistrictVariantPriceUpdate,
			Payload: PendingInputPayload{
				PendingValueDistrictID:   strconv.Itoa(districtID),
				PendingValueDistrictName: districtName,
				PendingValueProductID:    strconv.Itoa(productID),
				PendingValueProductName:  productName,
				PendingValueVariantID:    strconv.Itoa(selected.ID),
				PendingValueVariantName:  variantDisplayLabel,
				PendingValueCurrentPrice: strconv.Itoa(selected.Price),
			},
		}

		return buildAdminDistrictVariantPriceUpdateInputView(
			districtName,
			variantDisplayLabel,
			formatDistrictPlacementVariantPrice(selected.Price, selected.PriceText),
			"",
		), session, true, nil
	}

	return ViewModel{}, session, false, nil
}
