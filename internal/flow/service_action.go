package flow

import (
	"context"
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
	}

	if vm, nextSession, handled, err := s.handleAdminDistrictVariantAction(ctx, session, req); handled {
		if err != nil {
			return ViewModel{}, err
		}

		s.store.Put(req.SessionKey, nextSession)
		return vm, nil
	}

	if vm, nextSession, handled, err := s.handleAdminDistrictVariantPriceUpdateAction(ctx, session, req); handled {
		if err != nil {
			return ViewModel{}, err
		}

		s.store.Put(req.SessionKey, nextSession)
		return vm, nil
	}

	if vm, nextSession, handled, err := s.handleAdminCatalogCreateAction(ctx, session, req); handled {
		if err != nil {
			return ViewModel{}, err
		}

		s.store.Put(req.SessionKey, nextSession)
		return vm, nil
	}

	if vm, nextSession, handled, err := s.handleOrderAction(catalog, session, req); handled {
		if err != nil {
			return ViewModel{}, err
		}

		s.store.Put(req.SessionKey, nextSession)
		return vm, nil
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
		ScreenAdminDistrictVariantPriceUpdatePrice,
		ScreenAdminDistrictVariantCitySelect,
		ScreenAdminDistrictVariantDistrictSelect,
		ScreenAdminDistrictVariantCategorySelect,
		ScreenAdminDistrictVariantProductSelect,
		ScreenAdminDistrictVariantVariantSelect,
		ScreenAdminDistrictVariantPrice:
		return true
	default:
		return false
	}
}
