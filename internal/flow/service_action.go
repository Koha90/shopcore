package flow

import (
	"context"
	"errors"
	"strconv"
)

// HandleAction resolve the next flow view for an action.
//
// Resolution order:
//   - ActionBack uses session history
//   - ActionCatalogStart opens scenario-aware catalog root
//   - admin actions open stable admin screen or start pending text input
//   - generic catalog selection actions advance inside CatalogSchema
//   - explicit non-catalog action open stable detail screen
//
// Any non-pending action transition clears active pending text input state.
func (s *Service) HandleAction(ctx context.Context, req ActionRequest) (ViewModel, error) {
	catalog, err := s.provider.Catalog(ctx)
	if err != nil {
		return ViewModel{}, err
	}

	session, ok := s.store.Get(req.SessionKey)
	if !ok {
		session = Session{
			Current:  startScreenForScenario(req.StartScenario),
			History:  nil,
			Pending:  PendingInput{},
			CanAdmin: req.CanAdmin,
		}
	} else {
		session = s.syncSessionAccess(req.SessionKey, session, req.CanAdmin, req.StartScenario)
	}
	if isAdminAction(req.ActionID) && !session.CanAdmin {
		return ViewModel{}, ErrUnknownAction
	}

	switch req.ActionID {
	case ActionBack:
		if len(session.History) == 0 {
			if session.Pending.Active() {
				session.Pending = PendingInput{}
				s.store.Put(req.SessionKey, session)
			}
			return s.renderScreen(catalog, session, req.CanAdmin), nil
		}

		prev := session.History[len(session.History)-1]
		session.History = session.History[:len(session.History)-1]
		session.Current = prev
		if !shouldPreservePendingOnBack(prev) {
			session.Pending = PendingInput{}
		}
		s.store.Put(req.SessionKey, session)

		return s.renderScreen(catalog, session, req.CanAdmin), nil

	case ActionCatalogStart:
		next := catalogRootForScenario(req.StartScenario)

		if next != session.Current {
			session.History = append(session.History, session.Current)
			session.Current = next
		}
		session.Pending = PendingInput{}
		s.store.Put(req.SessionKey, session)

		return s.renderScreen(catalog, session, req.CanAdmin), nil

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
		s.store.Put(req.SessionKey, session)

		return s.renderScreen(catalog, session, req.CanAdmin), nil

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
		s.store.Put(req.SessionKey, session)

		return s.renderScreen(catalog, session, req.CanAdmin), nil

	case ActionAdminDistrictCreateStart:
		if s.cityLister == nil {
			return ViewModel{}, errors.New("flow city lister is nil")
		}
		next := ScreenAdminDistrictCitySelect

		if next != session.Current {
			session.History = append(session.History, session.Current)
		}
		session.Current = next
		session.Pending = PendingInput{}
		s.store.Put(req.SessionKey, session)

		return s.renderScreen(catalog, session, req.CanAdmin), nil

	case ActionAdminProductCreateStart:
		if s.categoryLister == nil {
			return ViewModel{}, errors.New("flow category lister is nil")
		}

		next := ScreenAdminProductCategorySelect

		if next != session.Current {
			session.History = append(session.History, session.Current)
		}
		session.Current = next
		session.Pending = PendingInput{}
		s.store.Put(req.SessionKey, session)

		return s.renderScreen(catalog, session, req.CanAdmin), nil

	case ActionAdminVariantCreateStart:
		if s.productLister == nil {
			return ViewModel{}, errors.New("flow product lister is nil")
		}

		next := ScreenAdminVariantProductSelect

		if next != session.Current {
			session.History = append(session.History, session.Current)
		}
		session.Current = next
		session.Pending = PendingInput{}
		s.store.Put(req.SessionKey, session)

		return s.renderScreen(catalog, session, req.CanAdmin), nil

	case ActionAdminDistrictVariantCreateStart:
		if s.districtLister == nil {
			return ViewModel{}, errors.New("flow district lister is nil")
		}

		next := ScreenAdminDistrictVariantDistrictSelect
		if next != session.Current {
			session.History = append(session.History, session.Current)
		}
		session.Current = next
		session.Pending = PendingInput{}
		s.store.Put(req.SessionKey, session)

		return s.renderScreen(catalog, session, req.CanAdmin), nil

	case ActionAdminDistrictVariantPriceUpdateStart:
		if s.districtLister == nil {
			return ViewModel{}, errors.New("flow district lister is nil")
		}

		next := ScreenAdminDistrictVariantPriceUpdateDistrictSelect
		if next != session.Current {
			session.History = append(session.History, session.Current)
		}
		session.Current = next
		session.Pending = PendingInput{}
		s.store.Put(req.SessionKey, session)

		return s.renderScreen(catalog, session, req.CanAdmin), nil
	}

	if cityID, ok := parseAdminDistrictSelectCityAction(req.ActionID); ok {
		if !session.CanAdmin {
			return ViewModel{}, ErrUnknownAction
		}
		if s.cityLister == nil {
			return ViewModel{}, errors.New("flow city lister is nil")
		}

		cities, err := s.cityLister.ListCities(ctx)
		if err != nil {
			return ViewModel{}, err
		}

		var selected *CityListItem
		for i := range cities {
			if cities[i].ID == cityID {
				selected = &cities[i]
				break
			}
		}
		if selected == nil {
			return ViewModel{}, ErrUnknownAction
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
		s.store.Put(req.SessionKey, session)

		return buildAdminDistrictCreateInputView(selected.Label, ""), nil
	}

	if categoryID, ok := parseAdminProductSelectCategoryAction(req.ActionID); ok {
		if !session.CanAdmin {
			return ViewModel{}, ErrUnknownAction
		}
		if s.categoryLister == nil {
			return ViewModel{}, errors.New("flow category lister is nil")
		}

		categories, err := s.categoryLister.ListCategories(ctx)
		if err != nil {
			return ViewModel{}, err
		}

		var selected *CategoryListItem
		for i := range categories {
			if categories[i].ID == categoryID {
				selected = &categories[i]
				break
			}
		}
		if selected == nil {
			return ViewModel{}, ErrUnknownAction
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
		s.store.Put(req.SessionKey, session)

		return buildAdminProductCreateInputView(selected.Label, ""), nil
	}

	if categoryID, ok := parseAdminDistrictVariantPriceUpdateSelectCategoryAction(req.ActionID); ok {
		if !session.CanAdmin {
			return ViewModel{}, ErrUnknownAction
		}
		if session.Current != ScreenAdminDistrictVariantPriceUpdateCategorySelect {
			return ViewModel{}, ErrUnknownAction
		}
		if s.districtPlacements == nil {
			return ViewModel{}, errors.New("flow district placements reader is nil")
		}

		districtID, ok := pendingDistrictID(session.Pending)
		if !ok {
			return ViewModel{}, errors.New("pending district id is invalid")
		}
		districtName := session.Pending.Value(PendingValueDistrictName)
		if districtName == "" {
			return ViewModel{}, errors.New("pending district name is empty")
		}

		categories, err := s.districtPlacements.ListDistrictCategories(ctx, districtID)
		if err != nil {
			return ViewModel{}, err
		}

		var selected *CategoryListItem
		for i := range categories {
			if categories[i].ID == categoryID {
				selected = &categories[i]
				break
			}
		}
		if selected == nil {
			return ViewModel{}, ErrUnknownAction
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
		s.store.Put(req.SessionKey, session)

		return s.buildAdminDistrictVariantPriceUpdateProductSelectScreen(
			districtID,
			districtName,
			selected.ID,
			selected.Label,
		), nil
	}

	if productID, ok := parseAdminVariantSelectProductAction(req.ActionID); ok {
		if !session.CanAdmin {
			return ViewModel{}, ErrUnknownAction
		}
		if s.productLister == nil {
			return ViewModel{}, errors.New("flow product lister is nil")
		}

		products, err := s.productLister.ListProducts(ctx)
		if err != nil {
			return ViewModel{}, err
		}

		var selected *ProductListItem
		for i := range products {
			if products[i].ID == productID {
				selected = &products[i]
				break
			}
		}
		if selected == nil {
			return ViewModel{}, ErrUnknownAction
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
		s.store.Put(req.SessionKey, session)

		return buildAdminVariantCreateInputView(selected.Label, ""), nil
	}

	if productID, ok := parseAdminDistrictVariantPriceUpdateSelectProductAction(req.ActionID); ok {
		if !session.CanAdmin {
			return ViewModel{}, ErrUnknownAction
		}
		if session.Current != ScreenAdminDistrictVariantPriceUpdateProductSelect {
			return ViewModel{}, ErrUnknownAction
		}
		if s.districtPlacements == nil {
			return ViewModel{}, errors.New("flow district placements reader is nil")
		}

		districtID, ok := pendingDistrictID(session.Pending)
		if !ok {
			return ViewModel{}, errors.New("pending district id is invalid")
		}
		districtName := session.Pending.Value(PendingValueDistrictName)
		if districtName == "" {
			return ViewModel{}, errors.New("pending district name is empty")
		}

		categoryID, ok := pendingCategoryID(session.Pending)
		if !ok {
			return ViewModel{}, errors.New("pending category id is invalid")
		}
		categoryName := session.Pending.Value(PendingValueCategoryName)
		if categoryName == "" {
			return ViewModel{}, errors.New("pending category name is empty")
		}

		products, err := s.districtPlacements.ListDistrictProducts(ctx, districtID, categoryID)
		if err != nil {
			return ViewModel{}, err
		}

		var selected *ProductListItem
		for i := range products {
			if products[i].ID == productID {
				selected = &products[i]
				break
			}
		}
		if selected == nil {
			return ViewModel{}, ErrUnknownAction
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
		s.store.Put(req.SessionKey, session)

		return s.buildAdminDistrictVariantPriceUpdateVariantSelectScreen(
			districtID,
			districtName,
			selected.ID,
			selected.Label,
		), nil
	}

	if districtID, ok := parseAdminDistrictVariantSelectDistrictAction(req.ActionID); ok {
		if !session.CanAdmin {
			return ViewModel{}, ErrUnknownAction
		}
		if s.districtLister == nil {
			return ViewModel{}, errors.New("flow distict lister is nil")
		}

		districts, err := s.districtLister.ListDistricts(ctx)
		if err != nil {
			return ViewModel{}, err
		}

		var selected *DistrictListItem
		for i := range districts {
			if districts[i].ID == districtID {
				selected = &districts[i]
				break
			}
		}
		if selected == nil {
			return ViewModel{}, ErrUnknownAction
		}

		switch session.Current {
		case ScreenAdminDistrictVariantDistrictSelect:
			if s.variantLister == nil {
				return ViewModel{}, errors.New("flow variant lister is nil")
			}

			next := ScreenAdminDistrictVariantVariantSelect
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
			s.store.Put(req.SessionKey, session)

			return s.buildAdminDistrictVariantVariantSelectScreen(selected.Label), nil

		case ScreenAdminDistrictVariantPriceUpdateDistrictSelect:
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
			s.store.Put(req.SessionKey, session)

			return s.buildAdminDistrictVariantPriceUpdateCategorySelectScreen(selected.ID, selected.Label), nil

		default:
			return ViewModel{}, ErrUnknownAction
		}
	}
	if variantID, ok := parseAdminDistrictVariantSelectVariantAction(req.ActionID); ok {
		if !session.CanAdmin {
			return ViewModel{}, ErrUnknownAction
		}

		districtID, ok := pendingDistrictID(session.Pending)
		if !ok {
			return ViewModel{}, errors.New("pending district id is invalid")
		}
		districtName := session.Pending.Value(PendingValueDistrictName)
		if districtName == "" {
			return ViewModel{}, errors.New("pending district name is empty")
		}

		switch session.Current {
		case ScreenAdminDistrictVariantVariantSelect:
			if s.variantLister == nil {
				return ViewModel{}, errors.New("flow variant lister is nil")
			}

			variants, err := s.variantLister.ListVariants(ctx)
			if err != nil {
				return ViewModel{}, err
			}

			var selected *VariantListItem
			for i := range variants {
				if variants[i].ID == variantID {
					selected = &variants[i]
					break
				}
			}
			if selected == nil {
				return ViewModel{}, ErrUnknownAction
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
			s.store.Put(req.SessionKey, session)

			return buildAdminDistrictVariantPriceInputView(districtName, variantDisplayLabel, ""), nil

		case ScreenAdminDistrictVariantPriceUpdateVariantSelect:
			if s.districtPlacements == nil {
				return ViewModel{}, errors.New("flow district placements reader is nil")
			}

			productID, ok := pendingProductID(session.Pending)
			if !ok {
				return ViewModel{}, errors.New("pending product id is invalid")
			}

			productName := session.Pending.Value(PendingValueProductName)
			if productName == "" {
				return ViewModel{}, errors.New("pending product name is empty")
			}

			variants, err := s.districtPlacements.ListDistrictVariants(ctx, districtID, productID)
			if err != nil {
				return ViewModel{}, err
			}

			var selected *DistrictPlacementVariantListItem
			for i := range variants {
				if variants[i].ID == variantID {
					selected = &variants[i]
					break
				}
			}
			if selected == nil {
				return ViewModel{}, ErrUnknownAction
			}

			variantDisplayLabel := buildAdminQualifiedVariantOptionLabel(productName, selected.Label)

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
			s.store.Put(req.SessionKey, session)

			return buildAdminDistrictVariantPriceUpdateInputView(
				districtName,
				variantDisplayLabel,
				formatDistrictPlacementVariantPrice(selected.Price, selected.PriceText),
				"",
			), nil

		default:
			return ViewModel{}, ErrUnknownAction
		}
	}

	if next, err := s.resolveCatalogScreen(catalog, session.Current, req.ActionID); err == nil {
		if next != session.Current {
			session.History = append(session.History, session.Current)
			session.Current = next
		}
		session.Pending = PendingInput{}
		s.store.Put(req.SessionKey, session)

		return s.renderScreen(catalog, session, req.CanAdmin), nil
	}

	next, err := resolveNextScreen(req.ActionID)
	if err != nil {
		return ViewModel{}, err
	}

	if next != session.Current {
		session.History = append(session.History, session.Current)
		session.Current = next
	}
	session.Pending = PendingInput{}
	s.store.Put(req.SessionKey, session)

	return s.renderScreen(catalog, session, req.CanAdmin), nil
}

// resolveCatalogScreen resolves one generic catalog selection action
// relative to the current screen state.
//
// It validates:
//   - action payload format
//   - current catalog path
//   - expected next schema level
//   - existence of target node in catalog tree
func (s *Service) resolveCatalogScreen(catalog Catalog, current ScreenID, actionID ActionID) (ScreenID, error) {
	level, id, ok := parseCatalogSelectAction(actionID)
	if !ok {
		return "", ErrUnknownAction
	}

	var currentPath CatalogPath

	switch current {
	case ScreenRootCompact, ScreenRootExtended:
		currentPath = nil

	default:
		path, ok := parseCatalogScreen(current)
		if !ok {
			return "", ErrUnknownAction
		}
		currentPath = path
	}

	expectedLevel, ok := s.expectedNextCatalogLevel(catalog, currentPath)
	if !ok {
		return "", ErrUnknownAction
	}
	if level != expectedLevel {
		return "", ErrUnknownAction
	}

	nextPath := currentPath.Append(level, id)

	if _, ok := catalog.FindNode(nextPath); !ok {
		return "", ErrUnknownAction
	}

	return catalogScreen(nextPath), nil
}

// expectedNextCatalogLevel returns which catalog level may be selected next
// for the provided path.
func (s *Service) expectedNextCatalogLevel(catalog Catalog, path CatalogPath) (CatalogLevel, bool) {
	if len(path) == 0 {
		return catalog.RootLevel()
	}

	last, ok := path.Last()
	if !ok {
		return "", false
	}

	return catalog.Schema.Next(last.Level)
}

func resolveNextScreen(actionID ActionID) (ScreenID, error) {
	switch actionID {
	case ActionRootCompact:
		return ScreenRootCompact, nil

	case ActionRootExtended:
		return ScreenRootExtended, nil

	case ActionCabinetOpen:
		return ScreenCabinet, nil

	case ActionSupportOpen:
		return ScreenSupport, nil

	case ActionReviewsOpen:
		return ScreenReviews, nil

	case ActionBalanceOpen:
		return ScreenBalance, nil

	case ActionBotsMine:
		return ScreenBotsMine, nil

	case ActionOrderLast:
		return ScreenOrderLast, nil

	case ActionAdminOpen:
		return ScreenAdminRoot, nil

	case ActionAdminCatalogOpen:
		return ScreenAdminCatalog, nil

	default:
		return "", ErrUnknownAction
	}
}

func shouldPreservePendingOnBack(screen ScreenID) bool {
	switch screen {
	case ScreenAdminDistrictVariantPriceUpdateCategorySelect,
		ScreenAdminDistrictVariantPriceUpdateProductSelect,
		ScreenAdminDistrictVariantPriceUpdateVariantSelect,
		ScreenAdminDistrictVariantPriceUpdatePrice:
		return true
	default:
		return false
	}
}
