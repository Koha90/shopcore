package flow

import (
	"context"
	"errors"
	"strconv"
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
		return s.renderScreen(catalog, session, req.CanAdmin), nil
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

		return s.renderScreen(catalog, session, req.CanAdmin), nil

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

		return s.renderScreen(catalog, session, req.CanAdmin), nil

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

		return s.renderScreen(catalog, session, req.CanAdmin), nil

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

		return s.renderScreen(catalog, session, req.CanAdmin), nil

	case PendingInputDistrictName:
		name := strings.TrimSpace(req.Text)
		if name == "" {
			return buildAdminDistrictCreateInputView(
				session.Pending.Value(PendingValueCityName),
				"Название района не может быть пустым.",
			), nil
		}

		cityID, ok := pendingCityID(session.Pending)
		if !ok {
			return ViewModel{}, errors.New("pending district city id is invalid")
		}

		cityName := session.Pending.Value(PendingValueCityName)
		session.Pending.SetValue(PendingValueName, name)

		suggestedCode := catalogservice.SuggestCode(name)
		if suggestedCode == "" {
			session.Current = ScreenAdminDistrictCode
			session.Pending.Kind = PendingInputDistrictCode
			s.store.Put(req.SessionKey, session)

			return buildAdminDistrictCodeInputView(
				cityName,
				"Не удалось автоматически подобрать code.",
				""), nil
		}

		session.Pending.SetValue(PendingValueCode, suggestedCode)

		if s.districts == nil {
			return ViewModel{}, errors.New("flow district creator is nil")
		}

		err := s.districts.CreateDistrict(ctx, CreateDistrictParams{
			CityID: cityID,
			Code:   suggestedCode,
			Name:   name,
		})
		if err != nil {
			session.Current = ScreenAdminDistrictCode
			session.Pending.Kind = PendingInputDistrictCode
			s.store.Put(req.SessionKey, session)

			return buildAdminDistrictCodeInputView(
				cityName,
				"Не удалось создать район с автоматическим code.",
				"",
			), nil
		}

		session.Pending = PendingInput{}
		session.Current = ScreenAdminDistrictCreateDone
		s.store.Put(req.SessionKey, session)

		return s.renderScreen(catalog, session, req.CanAdmin), nil

	case PendingInputDistrictCode:
		code := strings.TrimSpace(req.Text)
		if code == "" {
			return buildAdminDistrictCodeInputView(
				session.Pending.Value(PendingValueCityName),
				"Code района не может быть пустым.",
				"",
			), nil
		}

		cityID, ok := pendingCityID(session.Pending)
		if !ok {
			return ViewModel{}, errors.New("pending district city id is invalid")
		}

		name := strings.TrimSpace(session.Pending.Value(PendingValueName))
		if name == "" {
			return ViewModel{}, errors.New("pending district name is empty")
		}

		cityName := session.Pending.Value(PendingValueCityName)
		session.Pending.SetValue(PendingValueCode, code)

		if s.districts == nil {
			return ViewModel{}, errors.New("flow district creator is nil")
		}

		err := s.districts.CreateDistrict(ctx, CreateDistrictParams{
			CityID: cityID,
			Code:   session.Pending.Value(PendingValueCode),
			Name:   name,
		})
		if err != nil {
			return buildAdminDistrictCodeInputView(
				cityName,
				"Не удалось создать район. Попробуйте другой code.",
				code,
			), nil
		}

		session.Pending = PendingInput{}
		session.Current = ScreenAdminDistrictCreateDone
		s.store.Put(req.SessionKey, session)

		return s.renderScreen(catalog, session, req.CanAdmin), nil

	case PendingInputProductName:
		name := strings.TrimSpace(req.Text)
		if name == "" {
			return buildAdminProductCreateInputView(
				session.Pending.Value(PendingValueCategoryName),
				"Название товара не может быть пустым.",
			), nil
		}

		categoryID, ok := pendingCategoryID(session.Pending)
		if !ok {
			return ViewModel{}, errors.New("pending product category id is invalid")
		}

		categoryName := session.Pending.Value(PendingValueCategoryName)
		session.Pending.SetValue(PendingValueName, name)

		suggestedCode := catalogservice.SuggestCode(name)
		if suggestedCode == "" {
			session.Current = ScreenAdminProductCode
			session.Pending.Kind = PendingInputProductCode
			s.store.Put(req.SessionKey, session)

			return buildAdminProductCodeInputView(
				categoryName,
				"Не удалось автоматически подобрать code.",
				"",
			), nil
		}

		session.Pending.SetValue(PendingValueCode, suggestedCode)

		if s.products == nil {
			return ViewModel{}, errors.New("flow product creator is nil")
		}

		err := s.products.CreateProduct(ctx, CreateProductParams{
			CategoryID: categoryID,
			Code:       suggestedCode,
			Name:       name,
		})
		if err != nil {
			session.Current = ScreenAdminProductCode
			session.Pending.Kind = PendingInputProductCode
			s.store.Put(req.SessionKey, session)

			return buildAdminProductCodeInputView(
				categoryName,
				"Не удалось создать товар с автоматическим code.",
				"",
			), nil
		}

		session.Pending = PendingInput{}
		session.Current = ScreenAdminProductCreateDone
		s.store.Put(req.SessionKey, session)

		return s.renderScreen(catalog, session, req.CanAdmin), nil

	case PendingInputProductCode:
		code := strings.TrimSpace(req.Text)
		if code == "" {
			return buildAdminProductCodeInputView(
				session.Pending.Value(PendingValueCategoryName),
				"Code товара не может быть пустым.",
				"",
			), nil
		}

		categoryID, ok := pendingCategoryID(session.Pending)
		if !ok {
			return ViewModel{}, errors.New("pending product category id is invalid")
		}

		name := strings.TrimSpace(session.Pending.Value(PendingValueName))
		if name == "" {
			return ViewModel{}, errors.New("pending product name is empty")
		}

		categoryName := session.Pending.Value(PendingValueCategoryName)
		session.Pending.SetValue(PendingValueCode, code)

		if s.products == nil {
			return ViewModel{}, errors.New("flow product creator is nil")
		}

		err := s.products.CreateProduct(ctx, CreateProductParams{
			CategoryID: categoryID,
			Code:       session.Pending.Value(PendingValueCode),
			Name:       name,
		})
		if err != nil {
			return buildAdminProductCodeInputView(
				categoryName,
				"Не удалось создать товар. Попробуйте другой code.",
				code,
			), nil
		}

		session.Pending = PendingInput{}
		session.Current = ScreenAdminProductCreateDone
		s.store.Put(req.SessionKey, session)

		return s.renderScreen(catalog, session, req.CanAdmin), nil

	case PendingInputVariantName:
		name := strings.TrimSpace(req.Text)
		if name == "" {
			return buildAdminVariantCreateInputView(
				session.Pending.Value(PendingValueProductName),
				"Название варианта не может быть пустым.",
			), nil
		}

		productID, ok := pendingProductID(session.Pending)
		if !ok {
			return ViewModel{}, errors.New("pending variant product id is invalid")
		}

		productName := session.Pending.Value(PendingValueProductName)
		session.Pending.SetValue(PendingValueName, name)

		suggestedCode := catalogservice.SuggestCode(name)
		if suggestedCode == "" {
			session.Current = ScreenAdminVariantCode
			session.Pending.Kind = PendingInputVariantCode
			s.store.Put(req.SessionKey, session)

			return buildAdminVariantCodeInputView(
				productName,
				"Не удалось автоматически подобрать code",
				"",
			), nil

		}

		session.Pending.SetValue(PendingValueCode, suggestedCode)

		if s.variants == nil {
			return ViewModel{}, errors.New("flow variant creator is nil")
		}

		err := s.variants.CreateVariant(ctx, CreateVariantParams{
			ProductID: productID,
			Code:      suggestedCode,
			Name:      name,
		})
		if err != nil {
			session.Current = ScreenAdminVariantCode
			session.Pending.Kind = PendingInputVariantCode
			s.store.Put(req.SessionKey, session)

			return buildAdminVariantCodeInputView(
				productName,
				"Не удалось создать вариант с автоматическим code.",
				"",
			), nil
		}

		session.Pending = PendingInput{}
		session.Current = ScreenAdminVariantCreateDone
		s.store.Put(req.SessionKey, session)

		return s.renderScreen(catalog, session, req.CanAdmin), nil

	case PendingInputVariantCode:
		code := strings.TrimSpace(req.Text)
		if code == "" {
			return buildAdminVariantCodeInputView(
				session.Pending.Value(PendingValueProductName),
				"Code варианта не может быть пустым.",
				"",
			), nil
		}

		productID, ok := pendingProductID(session.Pending)
		if !ok {
			return ViewModel{}, errors.New("pending variant product id is invalid")
		}

		name := strings.TrimSpace(session.Pending.Value(PendingValueName))
		if name == "" {
			return ViewModel{}, errors.New("pending variant name is empty")
		}

		productName := session.Pending.Value(PendingValueProductName)
		session.Pending.SetValue(PendingValueCode, code)

		if s.variants == nil {
			return ViewModel{}, errors.New("flow variant creator is nil")
		}

		err := s.variants.CreateVariant(ctx, CreateVariantParams{
			ProductID: productID,
			Code:      session.Pending.Value(PendingValueCode),
			Name:      name,
		})
		if err != nil {
			return buildAdminVariantCodeInputView(
				productName,
				"Не удалось создать вариант. Попробуйте другой code.",
				code,
			), nil
		}

		session.Pending = PendingInput{}
		session.Current = ScreenAdminVariantCreateDone
		s.store.Put(req.SessionKey, session)

		return s.renderScreen(catalog, session, req.CanAdmin), nil

	case PendingInputDistrictVariantPrice:
		raw := strings.TrimSpace(req.Text)
		if raw == "" {
			return buildAdminDistrictVariantPriceInputView(
				session.Pending.Value(PendingValueDistrictName),
				session.Pending.Value(PendingValueVariantName),
				"Цена не может быть пустой.",
			), nil
		}

		price, err := strconv.Atoi(raw)
		if err != nil || price <= 0 {
			return buildAdminDistrictVariantPriceInputView(
				session.Pending.Value(PendingValueDistrictName),
				session.Pending.Value(PendingValueVariantName),
				"Цена должна быть положительным числом.",
			), nil
		}

		districtID, ok := pendingDistrictID(session.Pending)
		if !ok {
			return ViewModel{}, errors.New("pending district id is invalid")
		}

		variantID, ok := pendingVariantID(session.Pending)
		if !ok {
			return ViewModel{}, errors.New("pending variant id is invalid")
		}

		if s.districtVariants == nil {
			return ViewModel{}, errors.New("flow district variant creator is nil")
		}

		err = s.districtVariants.CreateDistrictVariant(ctx, CreateDistrictVariantParams{
			DistrictID: districtID,
			VariantID:  variantID,
			Price:      price,
		})
		if err != nil {
			validation := "Не удалось разместить вариант в районе."

			if errors.Is(err, catalogservice.ErrDistrictVariantAlreadyExists) {
				validation = "Вариант уже размещён в выбранном районе. Используйте обновление цены."
			}
			return buildAdminDistrictVariantPriceInputView(
				session.Pending.Value(PendingValueDistrictName),
				session.Pending.Value(PendingValueVariantName),
				validation,
			), nil
		}

		cityID := session.Pending.Value(PendingValueCityID)
		cityName := session.Pending.Value(PendingValueCityName)
		districtIDStr := strconv.Itoa(districtID)
		districtName := session.Pending.Value(PendingValueDistrictName)
		categoryID := session.Pending.Value(PendingValueCategoryID)
		categoryName := session.Pending.Value(PendingValueCategoryName)
		productID := session.Pending.Value(PendingValueProductID)
		productName := session.Pending.Value(PendingValueProductName)

		session.Pending = PendingInput{
			Kind: PendingInputNone,
			Payload: PendingInputPayload{
				PendingValueCityID:       cityID,
				PendingValueCityName:     cityName,
				PendingValueDistrictID:   districtIDStr,
				PendingValueDistrictName: districtName,
				PendingValueCategoryID:   categoryID,
				PendingValueCategoryName: categoryName,
				PendingValueProductID:    productID,
				PendingValueProductName:  productName,
			},
		}
		session.Current = ScreenAdminDistrictVariantCreateDone
		session.History = trimHistoryToScreen(session.History, ScreenAdminDistrictVariantVariantSelect)
		s.store.Put(req.SessionKey, session)

		return s.renderScreen(catalog, session, req.CanAdmin), nil

	case PendingInputDistrictVariantPriceUpdate:
		raw := strings.TrimSpace(req.Text)
		if raw == "" {
			return buildAdminDistrictVariantPriceUpdateInputView(
				session.Pending.Value(PendingValueDistrictName),
				session.Pending.Value(PendingValueVariantName),
				currentPlacementPriceTextFromPending(session.Pending),
				"Цена не может быть пустой.",
			), nil
		}

		price, err := strconv.Atoi(raw)
		if err != nil || price <= 0 {
			return buildAdminDistrictVariantPriceUpdateInputView(
				session.Pending.Value(PendingValueDistrictName),
				session.Pending.Value(PendingValueVariantName),
				currentPlacementPriceTextFromPending(session.Pending),
				"Цена должна быть положительным числом.",
			), nil
		}

		districtID, ok := pendingDistrictID(session.Pending)
		if !ok {
			return ViewModel{}, errors.New("pending district id is invalid")
		}

		variantID, ok := pendingVariantID(session.Pending)
		if !ok {
			return ViewModel{}, errors.New("pending variant id is invalid")
		}

		if s.districtVariantPrices == nil {
			return ViewModel{}, errors.New("flow district variant price updater is nil")
		}

		err = s.districtVariantPrices.UpdateDistrictVariantPrice(ctx, UpdateDistrictVariantPriceParams{
			DistrictID: districtID,
			VariantID:  variantID,
			Price:      price,
		})
		if err != nil {
			return buildAdminDistrictVariantPriceUpdateInputView(
				session.Pending.Value(PendingValueDistrictName),
				session.Pending.Value(PendingValueVariantName),
				currentPlacementPriceTextFromPending(session.Pending),
				"Не удалось обновить цену варианта.",
			), nil
		}

		session.Pending = PendingInput{}
		session.Current = ScreenAdminDistrictVariantPriceUpdateDone
		session.History = trimHistoryToScreen(session.History, ScreenAdminCatalog)
		s.store.Put(req.SessionKey, session)

		return s.renderScreen(catalog, session, req.CanAdmin), nil

	default:
		return ViewModel{}, ErrUnknownPendingInput
	}
}

func trimHistoryToScreen(history []ScreenID, targer ScreenID) []ScreenID {
	for i := len(history) - 1; i >= 0; i-- {
		if history[i] == targer {
			return append([]ScreenID(nil), history[:i+1]...)
		}
	}

	return nil
}
