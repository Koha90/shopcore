package flow

import (
	"context"
	"strconv"
	"strings"
)

const (
	adminProductImageSelectProductPrefix = "admin:catalog:product_image:select_product:"
	adminVariantImageSelectVariantPrefix = "admin:catalog:variant_image:select_variant:"
)

func buildAdminProductImageProductSelectView(products []ProductListItem, validation string) ViewModel {
	text := buildAdminSelectText(
		"Изображение товара",
		nil,
		validation,
		"Выберите товар:",
	)

	actions := make([]ActionButton, 0, len(products)+1)
	for _, product := range products {
		actions = append(actions, ActionButton{
			ID:    adminProductImageSelectProductAction(product.ID),
			Label: product.Label,
		})
	}

	actions = append(actions, ActionButton{
		ID: ActionBack, Label: "Назад",
	})

	return ViewModel{
		Text: text,
		Inline: &InlineKeyboardView{
			Sections: []ActionSection{
				{Columns: 1, Actions: actions},
			},
		},
		RemoveReply: true,
	}
}

func buildAdminVariantImageVariantSelectView(variants []VariantListItem, validation string) ViewModel {
	text := buildAdminSelectText(
		"Изображение варианта",
		nil,
		validation,
		"Выберите вариант:",
	)

	actions := make([]ActionButton, 0, len(variants)+1)
	for _, variant := range variants {
		label := variant.Label
		if variant.ProductLabel != "" {
			label = variant.ProductLabel + " · " + variant.Label
		}

		actions = append(actions, ActionButton{
			ID:    adminVariantImageSelectVariantAction(variant.ID),
			Label: label,
		})
	}

	actions = append(actions, ActionButton{ID: ActionBack, Label: "Назад"})

	return ViewModel{
		Text: text,
		Inline: &InlineKeyboardView{
			Sections: []ActionSection{
				{
					Columns: 1,
					Actions: actions,
				},
			},
		},
		RemoveReply: true,
	}
}

func buildAdminProductImageInputView(productName, validation string) ViewModel {
	text := buildAdminTextWithValidation(
		"Изображение товара",
		[]string{
			formatAdminFieldLine("Товар", productName),
		},
		validation,
		"Введите путь или URL изображения сообщением.\nНапример: assets/catalog/products/rose-box.jpg",
	)

	return ViewModel{
		Text: text,
		Inline: &InlineKeyboardView{
			Sections: []ActionSection{
				{
					Columns: 1,
					Actions: []ActionButton{
						{
							ID: ActionBack, Label: "Назад",
						},
					},
				},
			},
		},
		RemoveReply: true,
	}
}

func buildAdminVariantImageInputView(variantName, validation string) ViewModel {
	text := buildAdminTextWithValidation(
		"Изображение варианта",
		[]string{
			formatAdminFieldLine("Вариант", variantName),
		},
		validation,
		"Введите путь или URL изображения сообщением.\nНапример: assets/catalog/variants/rose-box-large.jpg",
	)

	return ViewModel{
		Text: text,
		Inline: &InlineKeyboardView{
			Sections: []ActionSection{
				{
					Columns: 1,
					Actions: []ActionButton{
						{
							ID: ActionBack, Label: "Назад",
						},
					},
				},
			},
		},
		RemoveReply: true,
	}
}

func buildAdminProductImageDoneView() ViewModel {
	return ViewModel{
		Text: "Изображение товара\n\nИзображение сохранено.",
		Inline: &InlineKeyboardView{
			Sections: []ActionSection{
				{
					Columns: 1,
					Actions: []ActionButton{
						{ID: ActionAdminCatalogOpen, Label: "В главное меню"},
						{ID: ActionBack, Label: "Назад"},
					},
				},
			},
		},
		RemoveReply: true,
	}
}

func buildAdminVariantImageDoneView() ViewModel {
	return ViewModel{
		Text: "Изображение варианта\n\nИзображение сохранено.",
		Inline: &InlineKeyboardView{
			Sections: []ActionSection{
				{
					Columns: 1,
					Actions: []ActionButton{
						{ID: ActionAdminCatalogOpen, Label: "В главное меню"},
						{ID: ActionBack, Label: "Назад"},
					},
				},
			},
		},
		RemoveReply: true,
	}
}

func (s *Service) buildAdminProductImageProductSelectScreen(ctx context.Context) ViewModel {
	if s == nil || s.productLister == nil {
		return buildAdminProductImageProductSelectView(nil, "Не удалось загрузить список товаров.")
	}

	products, err := s.productLister.ListProducts(ctx)
	if err != nil {
		return buildAdminProductImageProductSelectView(nil, "Не удалось загрузить список товаров.")
	}
	if len(products) == 0 {
		return buildAdminProductImageProductSelectView(nil, "Нет доступных товаров.")
	}

	return buildAdminProductImageProductSelectView(products, "")
}

func (s *Service) buildAdminVariantImageVariantSelectScreen(ctx context.Context) ViewModel {
	if s == nil || s.variantLister == nil {
		return buildAdminVariantImageVariantSelectView(nil, "Не удалось загрузить список вариантов.")
	}

	variants, err := s.variantLister.ListVariants(ctx)
	if err != nil {
		return buildAdminVariantImageVariantSelectView(nil, "Не удалось загрузить список вариантов.")
	}
	if len(variants) == 0 {
		return buildAdminVariantImageVariantSelectView(nil, "Нет доступных вариантов.")
	}

	return buildAdminVariantImageVariantSelectView(variants, "")
}

func (s *Service) handleAdminImageAction(
	ctx context.Context,
	session Session,
	req ActionRequest,
) (ViewModel, Session, bool) {
	switch req.ActionID {
	case ActionAdminProductImageUpdateStart:
		next := ScreenAdminProductImageProductSelect
		if next != session.Current {
			session.History = append(session.History, session.Current)
			session.Current = next
		}
		session.Pending = PendingInput{}

		return s.buildAdminProductImageProductSelectScreen(ctx), session, true

	case ActionAdminVariantImageUpdateStart:
		next := ScreenAdminVariantImageVariantSelect
		if next != session.Current {
			session.History = append(session.History, session.Current)
			session.Current = next
		}
		session.Pending = PendingInput{}

		return s.buildAdminVariantImageVariantSelectScreen(ctx), session, true
	}

	if productID, ok := parseAdminProductImageSelectProductAction(req.ActionID); ok {
		productName := "товар #" + strconv.Itoa(productID)
		if name, found := findProductLabel(ctx, s.productLister, productID); found {
			productName = name
		}

		session.History = append(session.History, session.Current)
		session.Current = ScreenAdminProductImageInput
		session.Pending = PendingInput{
			Kind: PendingInputProductImageURL,
			Payload: PendingInputPayload{
				PendingValueProductID:   strconv.Itoa(productID),
				PendingValueProductName: productName,
			},
		}

		return buildAdminProductImageInputView(productName, ""), session, true
	}

	if variantID, ok := parseAdminVariantImageSelectVariantAction(req.ActionID); ok {
		variantName := "вариант #" + strconv.Itoa(variantID)
		if name, found := findVariantLabel(ctx, s.variantLister, variantID); found {
			variantName = name
		}

		session.History = append(session.History, session.Current)
		session.Current = ScreenAdminVariantImageInput
		session.Pending = PendingInput{
			Kind: PendingInputVariantImageURL,
			Payload: PendingInputPayload{
				PendingValueVariantID:   strconv.Itoa(variantID),
				PendingValueVariantName: variantName,
			},
		}

		return buildAdminVariantImageInputView(variantName, ""), session, true
	}

	return ViewModel{}, session, false
}

func adminProductImageSelectProductAction(productID int) ActionID {
	return ActionID(adminProductImageSelectProductPrefix + strconv.Itoa(productID))
}

func parseAdminProductImageSelectProductAction(actionID ActionID) (int, bool) {
	return parsePositiveIDAction(actionID, adminProductImageSelectProductPrefix)
}

func adminVariantImageSelectVariantAction(variantID int) ActionID {
	return ActionID(adminVariantImageSelectVariantPrefix + strconv.Itoa(variantID))
}

func parseAdminVariantImageSelectVariantAction(actionID ActionID) (int, bool) {
	return parsePositiveIDAction(actionID, adminVariantImageSelectVariantPrefix)
}

func parsePositiveIDAction(actionID ActionID, prefix string) (int, bool) {
	raw := string(actionID)
	if !strings.HasPrefix(raw, prefix) {
		return 0, false
	}

	id, err := strconv.Atoi(strings.TrimPrefix(raw, prefix))
	if err != nil || id <= 0 {
		return 0, false
	}

	return id, true
}

func findProductLabel(ctx context.Context, lister ProductLister, productID int) (string, bool) {
	if lister == nil {
		return "", false
	}

	products, err := lister.ListProducts(ctx)
	if err != nil {
		return "", false
	}

	for _, product := range products {
		if product.ID == productID {
			return product.Label, true
		}
	}

	return "", false
}

func findVariantLabel(ctx context.Context, lister VariantLister, variantID int) (string, bool) {
	if lister == nil {
		return "", false
	}

	variants, err := lister.ListVariants(ctx)
	if err != nil {
		return "", false
	}

	for _, variant := range variants {
		if variant.ID == variantID {
			if variant.ProductLabel != "" {
				return variant.ProductLabel + " · " + variant.Label, true
			}

			return variant.Label, true
		}
	}

	return "", false
}
