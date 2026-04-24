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
	}
	if vm, handled := s.renderAdminCatalogCreateScreen(session); handled {
		return vm
	}

	if vm, handled := s.renderAdminDistrictVariantScreen(session); handled {
		return vm
	}

	if vm, handled := s.renderAdminDistrictVariantPriceUpdateScreen(session); handled {
		return vm
	}

	if vm, handled := s.renderOrderScreen(catalog, session); handled {
		return vm
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

	if session.CanAdmin != canAdmin {
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
