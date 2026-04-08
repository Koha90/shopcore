package flow

import (
	"context"
	"errors"
	"strings"

	catalogservice "github.com/koha90/shopcore/internal/catalog/service"
)

// HandleText resolves a plain text message relative to current session state.
//
// If no pending input exists, the current screen is rendered again.
// If pending input exists, text is handled as a continuation of that flow step.
//
// Current behavior supports admin category creation with automatic code
// suggestion and manual code fallback.
func (s *Service) HandleText(ctx context.Context, req TextRequest) (ViewModel, error) {
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

	if !session.CanAdmin && isAdminPending(session.Pending.Kind) {
		return ViewModel{}, ErrUnknownAction
	}

	if !session.Pending.Active() {
		return s.renderScreen(catalog, session.Current, req.CanAdmin), nil
	}

	switch session.Pending.Kind {
	case PendingInputCategoryName:
		name := strings.TrimSpace(req.Text)
		if name == "" {
			return buildAdminCategoryCreateInputView("Название категории не может быть пустым."), nil
		}

		session.Pending.SetValue(PendingValueName, name)

		suggestedCode := catalogservice.SuggestCode(name)
		if suggestedCode == "" {
			session.Current = ScreenAdminCategoryCode
			session.Pending.Kind = PendingInputCategoryCode
			s.store.Put(req.SessionKey, session)

			return buildAdminCategoryCodeInputView(
				"Не удалось автоматически подобрать code.",
				"",
			), nil
		}

		session.Pending.SetValue(PendingValueCode, suggestedCode)

		if s.categories == nil {
			return ViewModel{}, errors.New("flow category creator is nil")
		}

		err := s.categories.CreateCategory(ctx, CreateCategoryParams{
			Code: suggestedCode,
			Name: name,
		})
		if err != nil {
			session.Current = ScreenAdminCategoryCode
			session.Pending.Kind = PendingInputCategoryCode
			s.store.Put(req.SessionKey, session)

			return buildAdminCategoryCodeInputView(
				"Не удалось создать категорию с автоматическим code.",
				"",
			), nil
		}

		session.Pending = PendingInput{}
		session.Current = ScreenAdminCategoryCreateDone
		s.store.Put(req.SessionKey, session)

		return s.renderScreen(catalog, session.Current, req.CanAdmin), nil

	case PendingInputCategoryCode:
		code := strings.TrimSpace(req.Text)
		if code == "" {
			return buildAdminCategoryCodeInputView("Code категории не может быть пустым.", ""), nil
		}

		name := strings.TrimSpace(session.Pending.Value(PendingValueName))
		if name == "" {
			return ViewModel{}, errors.New("pending category name is empty")
		}

		session.Pending.SetValue(PendingValueCode, code)

		if s.categories == nil {
			return ViewModel{}, errors.New("flow category creator is nil")
		}

		err := s.categories.CreateCategory(ctx, CreateCategoryParams{
			Code: session.Pending.Value(PendingValueCode),
			Name: name,
		})
		if err != nil {
			return buildAdminCategoryCodeInputView("Не удалось создать категорию. Попробуйте другой code.", code), nil
		}

		session.Pending = PendingInput{}
		session.Current = ScreenAdminCategoryCreateDone
		s.store.Put(req.SessionKey, session)

		return s.renderScreen(catalog, session.Current, req.CanAdmin), nil

	case PendingInputCityName:
		name := strings.TrimSpace(req.Text)
		if name == "" {
			return buildAdminCityCreateInputView("Название города не может быть пустым."), nil
		}

		session.Pending.SetValue(PendingValueName, name)

		suggestedCode := catalogservice.SuggestCode(name)
		if suggestedCode == "" {
			session.Current = ScreenAdminCityCode
			session.Pending.Kind = PendingInputCityCode
			s.store.Put(req.SessionKey, session)

			return buildAdminCityCodeInputView(
				"Не удалось автоматически подобрать code.",
				"",
			), nil
		}

		session.Pending.SetValue(PendingValueCode, suggestedCode)

		if s.cities == nil {
			return ViewModel{}, errors.New("flow city creator is nil")
		}

		err := s.cities.CreateCity(ctx, CreateCityParams{
			Code: suggestedCode,
			Name: name,
		})
		if err != nil {
			session.Current = ScreenAdminCityCode
			session.Pending.Kind = PendingInputCityCode
			s.store.Put(req.SessionKey, session)

			return buildAdminCityCodeInputView(
				"Не удалось создать город с автоматическим code.",
				"",
			), nil
		}

		session.Pending = PendingInput{}
		session.Current = ScreenAdminCityCreateDone
		s.store.Put(req.SessionKey, session)

		return s.renderScreen(catalog, session.Current, req.CanAdmin), nil

	case PendingInputCityCode:
		code := strings.TrimSpace(req.Text)
		if code == "" {
			return buildAdminCityCodeInputView("Code города не может быть пустым.", ""), nil
		}

		name := strings.TrimSpace(session.Pending.Value(PendingValueName))
		if name == "" {
			return ViewModel{}, errors.New("pending city name is empty")
		}

		session.Pending.SetValue(PendingValueCode, code)

		if s.cities == nil {
			return ViewModel{}, errors.New("flow city creator is nil")
		}

		err := s.cities.CreateCity(ctx, CreateCityParams{
			Code: session.Pending.Value(PendingValueCode),
			Name: name,
		})
		if err != nil {
			return buildAdminCityCodeInputView("Не удалось создать город. Попробуйте другой code.", code), nil
		}

		session.Pending = PendingInput{}
		session.Current = ScreenAdminCityCreateDone
		s.store.Put(req.SessionKey, session)

		return s.renderScreen(catalog, session.Current, req.CanAdmin), nil

	default:
		return ViewModel{}, ErrUnknownPendingInput
	}
}
