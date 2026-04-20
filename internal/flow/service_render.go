package flow

// renderScreen converts logical screen identifiers into transport-agnostic view models.
//
// Stable root/detail screens are handled directly.
// Dynamic catalog drill-down screens are rendered from CatalogPath.
func (s *Service) renderScreen(catalog Catalog, session Session, canAdmin bool) ViewModel {
	screen := session.Current
	// pending := session.Pending

	switch screen {
	case ScreenReplyWelcome:
		return buildReplyWelcomeStart()

	case ScreenRootCompact:
		return buildCompactRootSelectionView(catalog.RootNodes())

	case ScreenRootExtended:
		return buildExtendedRootSelectionView(catalog.RootNodes(), canAdmin)

	case ScreenCabinet:
		return buildReplyDetailView(
			"Мой кабинет",
			"Здесь будут профиль, история, настройки и персональные данные пользователя.",
		)

	case ScreenSupport:
		return buildReplyDetailView(
			"Поддержка",
			"Здесь будет связь с оператором, FAQ и обработка обращений.",
		)

	case ScreenReviews:
		return buildReplyDetailView(
			"Отзывы",
			"Здесь будут отзывы клиентов, рейтинг и публикация новых отзывов.",
		)

	case ScreenBalance:
		return buildDetailView(
			"Баланс",
			"Здесь будет баланс аккаунта, пополнение и история операций.",
			ActionBack,
		)

	case ScreenBotsMine:
		return buildDetailView(
			"Мои боты",
			"Здесь будет список пользовательских ботов и быстрые действия по ним.",
			ActionBack,
		)

	case ScreenOrderLast:
		return buildDetailView(
			"Последний заказ",
			"Здесь будет карточка последнего заказа и повторное оформление.",
			ActionBack,
		)

	case ScreenAdminRoot:
		return buildAdminRootView()

	case ScreenAdminCatalog:
		return buildAdminCatalogView()

	case ScreenAdminCategoryCreate:
		return buildAdminCategoryCreateInputView("")

	case ScreenAdminCategoryCode:
		return buildAdminCategoryCodeInputView("", "")

	case ScreenAdminCategoryCreateDone:
		return buildAdminCategoryCreateDoneView()

	case ScreenAdminCityCreate:
		return buildAdminCityCreateInputView("")

	case ScreenAdminCityCode:
		return buildAdminCityCodeInputView("", "")

	case ScreenAdminCityCreateDone:
		return buildAdminCityCreateDoneView()

	case ScreenAdminDistrictCitySelect:
		return s.buildAdminDistrictCitySelectScreen()

	case ScreenAdminDistrictCreate:
		return buildAdminDistrictCreateInputView("", "")

	case ScreenAdminDistrictCode:
		return buildAdminDistrictCodeInputView("", "", "")

	case ScreenAdminDistrictCreateDone:
		return buildAdminDistrictCreateDoneView()

	case ScreenAdminProductCategorySelect:
		return s.buildAdminProductCategorySelectScreen()

	case ScreenAdminProductCreate:
		return buildAdminProductCreateInputView("", "")

	case ScreenAdminProductCode:
		return buildAdminProductCodeInputView("", "", "")

	case ScreenAdminProductCreateDone:
		return buildAdminProductCreateDoneView()

	case ScreenAdminVariantProductSelect:
		return s.buildAdminVariantProductSelectScreen()

	case ScreenAdminVariantCreate:
		return buildAdminVariantCreateInputView("", "")

	case ScreenAdminVariantCode:
		return buildAdminVariantCodeInputView("", "", "")

	case ScreenAdminVariantCreateDone:
		return buildAdminVariantCreateDoneView()

	case ScreenAdminDistrictVariantCitySelect:
		return s.buildAdminDistrictVariantCitySelectScreen()

	case ScreenAdminDistrictVariantDistrictSelect:
		cityID, ok := pendingCityID(session.Pending)
		if !ok {
			return buildAdminCatalogView()
		}
		cityName := session.Pending.Value(PendingValueCityName)
		return s.buildAdminDistrictVariantDistrictSelectScreen(cityID, cityName)

	case ScreenAdminDistrictVariantCategorySelect:
		cityName := session.Pending.Value(PendingValueCityName)
		districtName := session.Pending.Value(PendingValueDistrictName)

		return s.buildAdminDistrictVariantCategorySelectScreen(cityName, districtName)

	case ScreenAdminDistrictVariantProductSelect:
		cityName := session.Pending.Value(PendingValueCityName)
		districtName := session.Pending.Value(PendingValueDistrictName)
		categoryName := session.Pending.Value(PendingValueCategoryName)

		return s.buildAdminDistrictVariantProductSelectScreen(
			cityName,
			districtName,
			categoryName,
		)

	case ScreenAdminDistrictVariantVariantSelect:
		productID, ok := pendingProductID(session.Pending)
		if !ok {
			return buildAdminCatalogView()
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
		)

	case ScreenAdminDistrictVariantPrice:
		return buildAdminDistrictVariantPriceInputView("", "", "")

	case ScreenAdminDistrictVariantCreateDone:
		return buildAdminDistrictVariantCreateDoneView()

	case ScreenAdminDistrictVariantPriceUpdateDistrictSelect:
		return s.buildAdminDistrictVariantPriceUpdateDistrictSelectScreen()

	case ScreenAdminDistrictVariantPriceUpdateCategorySelect:
		districtID, ok := pendingDistrictID(session.Pending)
		if !ok {
			return buildAdminCatalogView()
		}
		districtName := session.Pending.Value(PendingValueDistrictName)
		return s.buildAdminDistrictVariantPriceUpdateCategorySelectScreen(districtID, districtName)

	case ScreenAdminDistrictVariantPriceUpdateProductSelect:
		districtID, ok := pendingDistrictID(session.Pending)
		if !ok {
			return buildAdminCatalogView()
		}
		categoryID, ok := pendingCategoryID(session.Pending)
		if !ok {
			return buildAdminCatalogView()
		}
		districtName := session.Pending.Value(PendingValueDistrictName)
		categoryName := session.Pending.Value(PendingValueCategoryName)
		return s.buildAdminDistrictVariantPriceUpdateProductSelectScreen(
			districtID,
			districtName,
			categoryID,
			categoryName,
		)

	case ScreenAdminDistrictVariantPriceUpdateVariantSelect:
		districtID, ok := pendingDistrictID(session.Pending)
		if !ok {
			return buildAdminCatalogView()
		}
		productID, ok := pendingProductID(session.Pending)
		if !ok {
			return buildAdminCatalogView()
		}
		districtName := session.Pending.Value(PendingValueDistrictName)
		productName := session.Pending.Value(PendingValueProductName)
		return s.buildAdminDistrictVariantPriceUpdateVariantSelectScreen(
			districtID,
			districtName,
			productID,
			productName,
		)

	case ScreenAdminDistrictVariantPriceUpdatePrice:
		return buildAdminDistrictVariantPriceUpdateInputView(
			session.Pending.Value(PendingValueDistrictName),
			session.Pending.Value(PendingValueVariantName),
			currentPlacementPriceTextFromPending(session.Pending),
			"",
		)

	case ScreenAdminDistrictVariantPriceUpdateDone:
		return buildAdminDistrictVariantPriceUpdateDoneView()
	}

	path, ok := parseCatalogScreen(screen)
	if !ok {
		return buildReplyWelcomeStart()
	}

	node, found := catalog.FindNode(path)
	if !found {
		return buildReplyWelcomeStart()
	}

	if len(node.Children) > 0 {
		return buildCatalogNodeView(node)
	}

	return buildCatalogLeafView(node)
}

func buildReplyWelcomeStart() ViewModel {
	return ViewModel{
		Text: "Добро пожаловать 👋\nВыберите раздел:",
		Reply: &ReplyKeyboardView{
			Rows: [][]ReplyButton{
				{
					{ID: ActionCatalogStart, Label: "♻️ Каталог"},
					{ID: ActionCabinetOpen, Label: "⚙️ Мой кабинет"},
				},
				{
					{ID: ActionSupportOpen, Label: "🤷‍♂️ Поддержка"},
					{ID: ActionReviewsOpen, Label: "📨 Отзывы"},
				},
			},
		},
	}
}

func buildDetailView(title, body string, backAction ActionID) ViewModel {
	return ViewModel{
		Text: title + "\n\n" + body,
		Inline: &InlineKeyboardView{
			Sections: []ActionSection{
				{
					Columns: 1,
					Actions: []ActionButton{
						{ID: backAction, Label: "Назад"},
					},
				},
			},
		},
		RemoveReply: true,
	}
}

// buildReplyDetailView renders a detail screen opened from reply-menu actions.
//
// It intentionally does not include inline back navigation, because those
// screens are not part of inline catalog flow.
func buildReplyDetailView(title, body string) ViewModel {
	return ViewModel{
		Text:        title + "\n\n" + body,
		RemoveReply: false,
	}
}

func (s *Service) syncSessionAccess(
	key SessionKey,
	session Session,
	canAdmin bool,
	startScenario string,
) Session {
	changed := false

	if session.CanAdmin == canAdmin {
		session.CanAdmin = canAdmin
		changed = true
	}

	if !session.CanAdmin && (isAdminScreen(session.Current) || isAdminPending(session.Pending.Kind)) {
		session.Current = catalogRootForScenario(startScenario)
		session.History = nil
		session.Pending = PendingInput{}
		changed = true
	}

	if changed {
		s.store.Put(key, session)
	}

	return session
}

// buildRootSelectionView renders the root inline selection screen.
//
// The compact variant renders only the main selectable entities.
// The extended variant renders the same entities plus utility action below.
func buildRootSelectionView(columns int, variant RootVariant, roots []CatalogNode, canAdmin bool) ViewModel {
	cols := normalizeColumns(columns)

	actions := make([]ActionButton, 0, len(roots))
	for _, node := range roots {
		actions = append(actions, ActionButton{
			ID:    catalogSelectAction(node.Level, node.ID),
			Label: node.Label,
		})
	}

	sections := []ActionSection{
		{
			Columns: cols,
			Actions: actions,
		},
	}

	if variant == RootVariantExtended {
		utilityActions := []ActionButton{
			{ID: ActionBalanceOpen, Label: "Баланс"},
			{ID: ActionBotsMine, Label: "Мои боты"},
			{ID: ActionOrderLast, Label: "Последний заказ"},
		}
		if canAdmin {
			utilityActions = append(utilityActions, ActionButton{
				ID:    ActionAdminOpen,
				Label: "Админка",
			})
		}
		sections = append(sections, ActionSection{
			Columns: 1,
			Actions: utilityActions,
		})
	}

	return ViewModel{
		Text: "Каталог\n\nВыберите раздел:",
		Inline: &InlineKeyboardView{
			Sections: sections,
		},
		RemoveReply: true,
	}
}
