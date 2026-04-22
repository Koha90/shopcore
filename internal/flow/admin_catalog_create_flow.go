package flow

import (
	"context"
	"errors"
	"strconv"
)

func (s *Service) handleAdminCatalogCreateAction(
	ctx context.Context,
	session Session,
	req ActionRequest,
) (ViewModel, Session, bool, error) {
	switch req.ActionID {
	case ActionAdminCategoryCreateStart:
		next := ScreenAdminCategoryCreate

		if next != session.Current {
			session.History = append(session.History, session.Current)
		}
		session.Current = next
		session.Pending = PendingInput{
			Kind:    PendingInputCategoryName,
			Payload: nil,
		}

		return buildAdminCategoryCreateInputView(""), session, true, nil

	case ActionAdminCityCreateStart:
		next := ScreenAdminCityCreate

		if next != session.Current {
			session.History = append(session.History, session.Current)
		}
		session.Current = next
		session.Pending = PendingInput{
			Kind:    PendingInputCityName,
			Payload: nil,
		}

		return buildAdminCityCreateInputView(""), session, true, nil

	case ActionAdminDistrictCreateStart:
		if s.cityLister == nil {
			return ViewModel{}, session, true, errors.New("flow city lister is nil")
		}

		next := ScreenAdminDistrictCitySelect
		if next != session.Current {
			session.History = append(session.History, session.Current)
		}
		session.Current = next
		session.Pending = PendingInput{}

		return s.buildAdminDistrictCitySelectScreen(), session, true, nil

	case ActionAdminProductCreateStart:
		if s.categoryLister == nil {
			return ViewModel{}, session, true, errors.New("flow category lister is nil")
		}

		next := ScreenAdminProductCategorySelect
		if next != session.Current {
			session.History = append(session.History, session.Current)
		}
		session.Current = next
		session.Pending = PendingInput{}

		return s.buildAdminProductCategorySelectScreen(), session, true, nil

	case ActionAdminVariantCreateStart:
		if s.productLister == nil {
			return ViewModel{}, session, true, errors.New("flow product lister is nil")
		}

		next := ScreenAdminVariantProductSelect
		if next != session.Current {
			session.History = append(session.History, session.Current)
		}
		session.Current = next
		session.Pending = PendingInput{}

		return s.buildAdminVariantProductSelectScreen(), session, true, nil
	}

	if cityID, ok := parseAdminDistrictSelectCityAction(req.ActionID); ok {
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

		next := ScreenAdminDistrictCreate
		if next != session.Current {
			session.History = append(session.History, session.Current)
		}
		session.Current = next
		session.Pending = PendingInput{
			Kind: PendingInputDistrictName,
			Payload: PendingInputPayload{
				PendingValueCityID:   strconv.Itoa(selected.ID),
				PendingValueCityName: selected.Label,
			},
		}

		return buildAdminDistrictCreateInputView(selected.Label, ""), session, true, nil
	}

	if categoryID, ok := parseAdminProductSelectCategoryAction(req.ActionID); ok {
		if s.categoryLister == nil {
			return ViewModel{}, session, true, errors.New("flow category lister is nil")
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

		next := ScreenAdminProductCreate
		if next != session.Current {
			session.History = append(session.History, session.Current)
		}
		session.Current = next
		session.Pending = PendingInput{
			Kind: PendingInputProductName,
			Payload: PendingInputPayload{
				PendingValueCategoryID:   strconv.Itoa(selected.ID),
				PendingValueCategoryName: selected.Label,
			},
		}

		return buildAdminProductCreateInputView(selected.Label, ""), session, true, nil
	}

	if productID, ok := parseAdminVariantSelectProductAction(req.ActionID); ok {
		if s.productLister == nil {
			return ViewModel{}, session, true, errors.New("flow product lister is nil")
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

		next := ScreenAdminVariantCreate
		if next != session.Current {
			session.History = append(session.History, session.Current)
		}
		session.Current = next
		session.Pending = PendingInput{
			Kind: PendingInputVariantName,
			Payload: PendingInputPayload{
				PendingValueProductID:   strconv.Itoa(selected.ID),
				PendingValueProductName: selected.Label,
			},
		}

		return buildAdminVariantCreateInputView(selected.Label, ""), session, true, nil
	}

	return ViewModel{}, session, false, nil
}
