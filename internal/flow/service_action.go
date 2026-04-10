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
			return s.renderScreen(catalog, session.Current, req.CanAdmin), nil
		}

		prev := session.History[len(session.History)-1]
		session.History = session.History[:len(session.History)-1]
		session.Current = prev
		session.Pending = PendingInput{}
		s.store.Put(req.SessionKey, session)

		return s.renderScreen(catalog, prev, req.CanAdmin), nil

	case ActionCatalogStart:
		next := catalogRootForScenario(req.StartScenario)

		if next != session.Current {
			session.History = append(session.History, session.Current)
			session.Current = next
		}
		session.Pending = PendingInput{}
		s.store.Put(req.SessionKey, session)

		return s.renderScreen(catalog, next, req.CanAdmin), nil

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

		return s.renderScreen(catalog, next, req.CanAdmin), nil

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

		return s.renderScreen(catalog, next, req.CanAdmin), nil

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

		return s.renderScreen(catalog, next, req.CanAdmin), nil

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

		return s.renderScreen(catalog, next, req.CanAdmin), nil
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

	if next, err := s.resolveCatalogScreen(catalog, session.Current, req.ActionID); err == nil {
		if next != session.Current {
			session.History = append(session.History, session.Current)
			session.Current = next
		}
		session.Pending = PendingInput{}
		s.store.Put(req.SessionKey, session)

		return s.renderScreen(catalog, next, req.CanAdmin), nil
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

	return s.renderScreen(catalog, next, req.CanAdmin), nil
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
