package flow

import (
	"context"
	"strconv"
)

func buildAdminRootView() ViewModel {
	return ViewModel{
		Text: "Админка\n\nВыберите раздел:",
		Inline: &InlineKeyboardView{
			Sections: []ActionSection{
				{
					Columns: 1,
					Actions: []ActionButton{
						{ID: ActionAdminCatalogOpen, Label: "Каталог"},
						{ID: ActionBack, Label: "Назад"},
					},
				},
			},
		},
		RemoveReply: true,
	}
}

func buildAdminCatalogView() ViewModel {
	return ViewModel{
		Text: "Админка · Каталог\n\nВыберите действие:",
		Inline: &InlineKeyboardView{
			Sections: []ActionSection{
				{
					Columns: 1,
					Actions: []ActionButton{
						{ID: ActionAdminDistrictCreateStart, Label: "Создать район"},
						{ID: ActionAdminCityCreateStart, Label: "Создать город"},
						{ID: ActionAdminCategoryCreateStart, Label: "Создать категорию"},
						{ID: ActionAdminProductCreateStart, Label: "Создать товар"},
						{ID: ActionAdminVariantCreateStart, Label: "Создать вариант"},
						{ID: ActionAdminDistrictVariantCreateStart, Label: "Разместить вариант"},
						{ID: ActionAdminDistrictVariantPriceUpdateStart, Label: "Изменить цену варианта"},
						{ID: ActionBack, Label: "Назад"},
					},
				},
			},
		},
		RemoveReply: true,
	}
}

func buildAdminCategoryCreateInputView(validation string) ViewModel {
	text := buildAdminTextWithValidation(
		"Новая категория",
		nil,
		validation,
		"Введите название категории сообщением.",
	)

	return ViewModel{
		Text: text,
		Inline: &InlineKeyboardView{
			Sections: []ActionSection{
				{
					Columns: 1,
					Actions: []ActionButton{
						{ID: ActionBack, Label: "Назад"},
					},
				},
			},
		},
		RemoveReply: true,
	}
}

func buildAdminCategoryCodeInputView(validation, suggested string) ViewModel {
	switch {
	case validation != "":
		text := buildAdminTextWithValidation(
			"Новая категория",
			nil,
			validation,
			"Введите code категории сообщением.",
		)

		return ViewModel{
			Text: text,
			Inline: &InlineKeyboardView{
				Sections: []ActionSection{
					{
						Columns: 1,
						Actions: []ActionButton{
							{ID: ActionBack, Label: "Назад"},
						},
					},
				},
			},
			RemoveReply: true,
		}

	case suggested != "":
		text := buildAdminText(
			"Новая категория",
			[]string{formatAdminAutoCodeLine(suggested)},
			"Введите code категории сообщением.",
		)

		return ViewModel{
			Text: text,
			Inline: &InlineKeyboardView{
				Sections: []ActionSection{
					{
						Columns: 1,
						Actions: []ActionButton{
							{ID: ActionBack, Label: "Назад"},
						},
					},
				},
			},
			RemoveReply: true,
		}

	default:
		text := buildAdminText(
			"Новая категория",
			nil,
			"Введите code категории сообщением.",
		)

		return ViewModel{
			Text: text,
			Inline: &InlineKeyboardView{
				Sections: []ActionSection{
					{
						Columns: 1,
						Actions: []ActionButton{
							{ID: ActionBack, Label: "Назад"},
						},
					},
				},
			},
			RemoveReply: true,
		}
	}
}

func buildAdminCategoryCreateDoneView() ViewModel {
	return ViewModel{
		Text: "Новая категория\n\nКатегория создана.",
		Inline: &InlineKeyboardView{
			Sections: []ActionSection{
				{
					Columns: 1,
					Actions: []ActionButton{
						{ID: ActionBack, Label: "Назад"},
					},
				},
			},
		},
		RemoveReply: true,
	}
}

func buildAdminCityCreateInputView(validation string) ViewModel {
	text := buildAdminTextWithValidation(
		"Новый город",
		nil,
		validation,
		"Введите название города сообщением.",
	)

	return ViewModel{
		Text: text,
		Inline: &InlineKeyboardView{
			Sections: []ActionSection{
				{
					Columns: 1,
					Actions: []ActionButton{
						{ID: ActionBack, Label: "Назад"},
					},
				},
			},
		},
		RemoveReply: true,
	}
}

func buildAdminCityCodeInputView(validation, suggested string) ViewModel {
	switch {
	case validation != "":
		text := buildAdminTextWithValidation(
			"Новый город",
			nil,
			validation,
			"Введите code города сообщением.",
		)

		return ViewModel{
			Text: text,
			Inline: &InlineKeyboardView{
				Sections: []ActionSection{
					{
						Columns: 1,
						Actions: []ActionButton{
							{ID: ActionBack, Label: "Назад"},
						},
					},
				},
			},
			RemoveReply: true,
		}

	case suggested != "":
		text := buildAdminText(
			"Новый город",
			[]string{formatAdminAutoCodeLine(suggested)},
			"Введите code города сообщением.",
		)

		return ViewModel{
			Text: text,
			Inline: &InlineKeyboardView{
				Sections: []ActionSection{
					{
						Columns: 1,
						Actions: []ActionButton{
							{ID: ActionBack, Label: "Назад"},
						},
					},
				},
			},
			RemoveReply: true,
		}

	default:
		text := buildAdminText(
			"Новый город",
			nil,
			"Введите code города сообщением.",
		)

		return ViewModel{
			Text: text,
			Inline: &InlineKeyboardView{
				Sections: []ActionSection{
					{
						Columns: 1,
						Actions: []ActionButton{
							{ID: ActionBack, Label: "Назад"},
						},
					},
				},
			},
			RemoveReply: true,
		}
	}
}

func buildAdminCityCreateDoneView() ViewModel {
	return ViewModel{
		Text: "Новый город\n\nГород создан.",
		Inline: &InlineKeyboardView{
			Sections: []ActionSection{
				{
					Columns: 1,
					Actions: []ActionButton{
						{ID: ActionBack, Label: "Назад"},
					},
				},
			},
		},
		RemoveReply: true,
	}
}

func buildAdminDistrictCitySelectView(cities []CityListItem, validation string) ViewModel {
	text := buildAdminSelectText(
		"Новый район",
		nil,
		validation,
		"Выберите город:",
	)

	actions := make([]ActionButton, 0, len(cities)+1)
	for _, city := range cities {
		actions = append(actions, ActionButton{
			ID:    adminDistrictSelectCityAction(city.ID),
			Label: city.Label,
		})
	}
	actions = append(actions, ActionButton{
		ID:    ActionBack,
		Label: "Назад",
	})

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

func buildAdminDistrictCreateInputView(cityName, validation string) ViewModel {
	text := buildAdminTextWithValidation(
		"Новый район",
		[]string{
			formatAdminFieldLine("Город", cityName),
		},
		validation,
		"Введите название района сообщением.",
	)

	return ViewModel{
		Text: text,
		Inline: &InlineKeyboardView{
			Sections: []ActionSection{
				{
					Columns: 1,
					Actions: []ActionButton{
						{ID: ActionBack, Label: "Назад"},
					},
				},
			},
		},
		RemoveReply: true,
	}
}

func buildAdminDistrictCodeInputView(cityName, validation, suggested string) ViewModel {
	fields := []string{
		formatAdminFieldLine("Город", cityName),
	}

	switch {
	case validation != "":
		text := buildAdminTextWithValidation(
			"Новый район",
			fields,
			validation,
			"Введите code района сообщением.",
		)

		return ViewModel{
			Text: text,
			Inline: &InlineKeyboardView{
				Sections: []ActionSection{
					{
						Columns: 1,
						Actions: []ActionButton{
							{ID: ActionBack, Label: "Назад"},
						},
					},
				},
			},
			RemoveReply: true,
		}

	case suggested != "":
		text := buildAdminText(
			"Новый район",
			append(fields, formatAdminAutoCodeLine(suggested)),
			"Введите code района сообщением.",
		)

		return ViewModel{
			Text: text,
			Inline: &InlineKeyboardView{
				Sections: []ActionSection{
					{
						Columns: 1,
						Actions: []ActionButton{
							{ID: ActionBack, Label: "Назад"},
						},
					},
				},
			},
			RemoveReply: true,
		}

	default:
		text := buildAdminText(
			"Новый район",
			fields,
			"Введите code района сообщением.",
		)

		return ViewModel{
			Text: text,
			Inline: &InlineKeyboardView{
				Sections: []ActionSection{
					{
						Columns: 1,
						Actions: []ActionButton{
							{ID: ActionBack, Label: "Назад"},
						},
					},
				},
			},
			RemoveReply: true,
		}
	}
}

func buildAdminDistrictCreateDoneView() ViewModel {
	return ViewModel{
		Text: "Новый район\n\nРайон создан.",
		Inline: &InlineKeyboardView{
			Sections: []ActionSection{
				{
					Columns: 1,
					Actions: []ActionButton{
						{ID: ActionBack, Label: "Назад"},
					},
				},
			},
		},
		RemoveReply: true,
	}
}

func (s *Service) buildAdminDistrictCitySelectScreen() ViewModel {
	if s == nil || s.cityLister == nil {
		return buildAdminDistrictCitySelectView(nil, "Не удалось загрузить список городов.")
	}

	cities, err := s.cityLister.ListCities(context.Background())
	if err != nil {
		return buildAdminDistrictCitySelectView(nil, "Не удалось загрузить список городов.")
	}
	if len(cities) == 0 {
		return buildAdminDistrictCitySelectView(nil, "Нет доступных городов.")
	}

	return buildAdminDistrictCitySelectView(cities, "")
}

func pendingCityID(p PendingInput) (int, bool) {
	raw := p.Value(PendingValueCityID)
	if raw == "" {
		return 0, false
	}

	id, err := strconv.Atoi(raw)
	if err != nil || id <= 0 {
		return 0, false
	}

	return id, true
}

func buildAdminProductCategorySelectView(categories []CategoryListItem, validation string) ViewModel {
	text := buildAdminSelectText(
		"Новый товар",
		nil,
		validation,
		"Выберите категорию:",
	)

	actions := make([]ActionButton, 0, len(categories)+1)
	for _, category := range categories {
		actions = append(actions, ActionButton{
			ID:    adminProductSelectCategoryAction(category.ID),
			Label: category.Label,
		})
	}
	actions = append(actions, ActionButton{
		ID:    ActionBack,
		Label: "Назад",
	})

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

func buildAdminProductCreateInputView(categoryName, validation string) ViewModel {
	text := buildAdminTextWithValidation(
		"Новый товар",
		[]string{
			formatAdminFieldLine("Категория", categoryName),
		},
		validation,
		"Введите название товара сообщением.",
	)

	return ViewModel{
		Text: text,
		Inline: &InlineKeyboardView{
			Sections: []ActionSection{
				{
					Columns: 1,
					Actions: []ActionButton{
						{ID: ActionBack, Label: "Назад"},
					},
				},
			},
		},
		RemoveReply: true,
	}
}

func buildAdminProductCodeInputView(categoryName, validation, suggested string) ViewModel {
	fields := []string{
		formatAdminFieldLine("Категория", categoryName),
	}

	switch {
	case validation != "":
		text := buildAdminTextWithValidation(
			"Новый товар",
			fields,
			validation,
			"Введите code товара сообщением.",
		)

		return ViewModel{
			Text: text,
			Inline: &InlineKeyboardView{
				Sections: []ActionSection{
					{
						Columns: 1,
						Actions: []ActionButton{
							{ID: ActionBack, Label: "Назад"},
						},
					},
				},
			},
			RemoveReply: true,
		}

	case suggested != "":
		text := buildAdminText(
			"Новый товар",
			append(fields, formatAdminAutoCodeLine(suggested)),
			"Введите code товара сообщением.",
		)

		return ViewModel{
			Text: text,
			Inline: &InlineKeyboardView{
				Sections: []ActionSection{
					{
						Columns: 1,
						Actions: []ActionButton{
							{ID: ActionBack, Label: "Назад"},
						},
					},
				},
			},
			RemoveReply: true,
		}

	default:
		text := buildAdminText(
			"Новый товар",
			fields,
			"Введите code товара сообщением.",
		)

		return ViewModel{
			Text: text,
			Inline: &InlineKeyboardView{
				Sections: []ActionSection{
					{
						Columns: 1,
						Actions: []ActionButton{
							{ID: ActionBack, Label: "Назад"},
						},
					},
				},
			},
			RemoveReply: true,
		}
	}
}

func buildAdminProductCreateDoneView() ViewModel {
	return ViewModel{
		Text: "Новый товар\n\nТовар создан.",
		Inline: &InlineKeyboardView{
			Sections: []ActionSection{
				{
					Columns: 1,
					Actions: []ActionButton{
						{ID: ActionBack, Label: "Назад"},
					},
				},
			},
		},
		RemoveReply: true,
	}
}

func (s *Service) buildAdminProductCategorySelectScreen() ViewModel {
	if s == nil || s.categoryLister == nil {
		return buildAdminProductCategorySelectView(nil, "Не удалось загрузить список категорий.")
	}

	categories, err := s.categoryLister.ListCategories(context.Background())
	if err != nil {
		return buildAdminProductCategorySelectView(nil, "Не удалось загрузить список категорий.")
	}
	if len(categories) == 0 {
		return buildAdminProductCategorySelectView(nil, "Нет доступных категорий.")
	}

	return buildAdminProductCategorySelectView(categories, "")
}

func pendingCategoryID(p PendingInput) (int, bool) {
	raw := p.Value(PendingValueCategoryID)
	if raw == "" {
		return 0, false
	}

	id, err := strconv.Atoi(raw)
	if err != nil || id <= 0 {
		return 0, false
	}

	return id, true
}

func buildAdminVariantProductSelectView(products []ProductListItem, validation string) ViewModel {
	text := buildAdminSelectText(
		"Новый вариант",
		nil,
		validation,
		"Выберите товар:",
	)

	actions := make([]ActionButton, 0, len(products)+1)
	for _, product := range products {
		actions = append(actions, ActionButton{
			ID:    adminVariantSelectProductAction(product.ID),
			Label: product.Label,
		})
	}
	actions = append(actions, ActionButton{
		ID:    ActionBack,
		Label: "Назад",
	})

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

func buildAdminVariantCreateInputView(productName, validation string) ViewModel {
	text := buildAdminTextWithValidation(
		"Новый вариант",
		[]string{
			formatAdminFieldLine("Товар", productName),
		},
		validation,
		"Введите название варианта сообщением.",
	)

	return ViewModel{
		Text: text,
		Inline: &InlineKeyboardView{
			Sections: []ActionSection{
				{
					Columns: 1,
					Actions: []ActionButton{
						{ID: ActionBack, Label: "Назад"},
					},
				},
			},
		},
		RemoveReply: true,
	}
}

func buildAdminVariantCodeInputView(productName, validation, suggested string) ViewModel {
	fields := []string{
		formatAdminFieldLine("Товар", productName),
	}

	switch {
	case validation != "":
		text := buildAdminTextWithValidation(
			"Новый вариант",
			fields,
			validation,
			"Введите code варианта сообщением.",
		)

		return ViewModel{
			Text: text,
			Inline: &InlineKeyboardView{
				Sections: []ActionSection{
					{
						Columns: 1,
						Actions: []ActionButton{
							{ID: ActionBack, Label: "Назад"},
						},
					},
				},
			},
			RemoveReply: true,
		}

	case suggested != "":
		text := buildAdminText(
			"Новый вариант",
			append(fields, formatAdminAutoCodeLine(suggested)),
			"Введите code варианта сообщением.",
		)

		return ViewModel{
			Text: text,
			Inline: &InlineKeyboardView{
				Sections: []ActionSection{
					{
						Columns: 1,
						Actions: []ActionButton{
							{ID: ActionBack, Label: "Назад"},
						},
					},
				},
			},
			RemoveReply: true,
		}

	default:
		text := buildAdminText(
			"Новый вариант",
			fields,
			"Введите code варианта сообщением.",
		)

		return ViewModel{
			Text: text,
			Inline: &InlineKeyboardView{
				Sections: []ActionSection{
					{
						Columns: 1,
						Actions: []ActionButton{
							{ID: ActionBack, Label: "Назад"},
						},
					},
				},
			},
			RemoveReply: true,
		}
	}
}

func buildAdminVariantCreateDoneView() ViewModel {
	return ViewModel{
		Text: "Новый вариант\n\nВариант создан.",
		Inline: &InlineKeyboardView{
			Sections: []ActionSection{
				{
					Columns: 1,
					Actions: []ActionButton{
						{ID: ActionBack, Label: "Назад"},
					},
				},
			},
		},
		RemoveReply: true,
	}
}

func (s *Service) buildAdminVariantProductSelectScreen() ViewModel {
	if s == nil || s.productLister == nil {
		return buildAdminVariantProductSelectView(nil, "Не удалось загрузить список товаров.")
	}

	products, err := s.productLister.ListProducts(context.Background())
	if err != nil {
		return buildAdminVariantProductSelectView(nil, "Не удалось загрузить список товаров.")
	}
	if len(products) == 0 {
		return buildAdminVariantProductSelectView(nil, "Нет доступных товаров.")
	}

	return buildAdminVariantProductSelectView(products, "")
}

func pendingProductID(p PendingInput) (int, bool) {
	raw := p.Value(PendingValueProductID)
	if raw == "" {
		return 0, false
	}

	id, err := strconv.Atoi(raw)
	if err != nil || id <= 0 {
		return 0, false
	}

	return id, true
}

func buildAdminDistrictVariantCitySelectView(cities []CityListItem, validation string) ViewModel {
	text := buildAdminSelectText(
		"Размещение варианта",
		nil,
		validation,
		"Выберите город:",
	)

	actions := make([]ActionButton, 0, len(cities)+1)
	for _, city := range cities {
		actions = append(actions, ActionButton{
			ID:    adminDistrictVariantSelectCityAction(city.ID),
			Label: city.Label,
		})
	}
	actions = append(actions, ActionButton{
		ID:    ActionBack,
		Label: "Назад",
	})

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

func buildAdminDistrictVariantDistrictSelectView(
	cityName string,
	districts []DistrictListItem,
	validation string,
) ViewModel {
	text := buildAdminSelectText(
		"Размещение варианта",
		[]string{
			formatAdminFieldLine("Город", cityName),
		},
		validation,
		"Выберите район:",
	)

	actions := make([]ActionButton, 0, len(districts)+1)
	for _, district := range districts {
		actions = append(actions, ActionButton{
			ID:    adminDistrictVariantSelectDistrictAction(district.ID),
			Label: district.Label,
		})
	}
	actions = append(actions, ActionButton{
		ID:    ActionBack,
		Label: "Назад",
	})

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

func buildAdminDistrictVariantVariantSelectView(
	districtName string,
	variants []VariantListItem,
	validation string,
) ViewModel {
	text := buildAdminSelectText(
		"Размещение варианта",
		[]string{
			formatAdminFieldLine("Район", districtName),
		},
		validation,
		"Выберите вариант:",
	)

	actions := make([]ActionButton, 0, len(variants)+1)
	for _, variant := range variants {

		variantDisplayLabel := buildAdminVariantOptionLabel(variant)

		actions = append(actions, ActionButton{
			ID:    adminDistrictVariantSelectVariantAction(variant.ID),
			Label: variantDisplayLabel,
		})
	}
	actions = append(actions, ActionButton{
		ID:    ActionBack,
		Label: "Назад",
	})

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

func buildAdminDistrictVariantPriceInputView(districtName, variantName, validation string) ViewModel {
	text := buildAdminTextWithValidation(
		"Размещение варианта",
		[]string{
			formatAdminFieldLine("Район", districtName),
			formatAdminFieldLine("Вариант", variantName),
		},
		validation,
		"Введите цену сообщением.",
	)

	return ViewModel{
		Text: text,
		Inline: &InlineKeyboardView{
			Sections: []ActionSection{
				{
					Columns: 1,
					Actions: []ActionButton{
						{ID: ActionBack, Label: "Назад"},
					},
				},
			},
		},
		RemoveReply: true,
	}
}

func buildAdminDistrictVariantCreateDoneView() ViewModel {
	return ViewModel{
		Text: "Размещение варианта\n\nВариант размещён в районе.",
		Inline: &InlineKeyboardView{
			Sections: []ActionSection{
				{
					Columns: 1,
					Actions: []ActionButton{
						{ID: ActionBack, Label: "Назад"},
					},
				},
			},
		},
		RemoveReply: true,
	}
}

func (s *Service) buildAdminDistrictVariantCitySelectScreen() ViewModel {
	if s == nil || s.cityLister == nil {
		return buildAdminDistrictVariantCitySelectView(nil, "Не удалось загрузить список городов.")
	}

	cities, err := s.cityLister.ListCities(context.Background())
	if err != nil {
		return buildAdminDistrictVariantCitySelectView(nil, "Не удалось загрузить список городов.")
	}
	if len(cities) == 0 {
		return buildAdminDistrictVariantCitySelectView(nil, "Нет доступных городов.")
	}

	return buildAdminDistrictVariantCitySelectView(cities, "")
}

func (s *Service) buildAdminDistrictVariantDistrictSelectScreen(cityID int, cityName string) ViewModel {
	if s == nil || s.districtLister == nil {
		return buildAdminDistrictVariantDistrictSelectView(cityName, nil, "Не удалось загрузить список районов.")
	}
	if cityID <= 0 {
		return buildAdminDistrictVariantDistrictSelectView(cityName, nil, "Город не выбран.")
	}

	districts, err := s.districtLister.ListDistrictsByCity(context.Background(), cityID)
	if err != nil {
		return buildAdminDistrictVariantDistrictSelectView(cityName, nil, "Не удалось загрузить список районов.")
	}
	if len(districts) == 0 {
		return buildAdminDistrictVariantDistrictSelectView(cityName, nil, "Нет доступных районов.")
	}

	return buildAdminDistrictVariantDistrictSelectView(cityName, districts, "")
}

func (s *Service) buildAdminDistrictVariantVariantSelectScreen(districtName string) ViewModel {
	if s == nil || s.variantLister == nil {
		return buildAdminDistrictVariantVariantSelectView(districtName, nil, "Не удалось загрузить список вариантов.")
	}

	variants, err := s.variantLister.ListVariants(context.Background())
	if err != nil {
		return buildAdminDistrictVariantVariantSelectView(districtName, nil, "Не удалось загрузить список вариантов.")
	}
	if len(variants) == 0 {
		return buildAdminDistrictVariantVariantSelectView(districtName, nil, "Нет доступных вариантов.")
	}

	return buildAdminDistrictVariantVariantSelectView(districtName, variants, "")
}

func pendingDistrictID(p PendingInput) (int, bool) {
	raw := p.Value(PendingValueDistrictID)
	if raw == "" {
		return 0, false
	}

	id, err := strconv.Atoi(raw)
	if err != nil || id <= 0 {
		return 0, false
	}

	return id, true
}

func pendingVariantID(p PendingInput) (int, bool) {
	raw := p.Value(PendingValueVariantID)
	if raw == "" {
		return 0, false
	}

	id, err := strconv.Atoi(raw)
	if err != nil || id <= 0 {
		return 0, false
	}

	return id, true
}

func buildAdminDistrictVariantPriceUpdateDistrictSelectView(districts []DistrictListItem, validation string) ViewModel {
	text := buildAdminSelectText(
		"Изменение цены варианта",
		nil,
		validation,
		"Выберите район:",
	)

	actions := make([]ActionButton, 0, len(districts)+1)
	for _, district := range districts {
		actions = append(actions, ActionButton{
			ID:    adminDistrictVariantSelectDistrictAction(district.ID),
			Label: district.Label,
		})
	}
	actions = append(actions, ActionButton{
		ID:    ActionBack,
		Label: "Назад",
	})

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

func buildAdminDistrictVariantPriceUpdateVariantSelectView(
	districtName, productName string,
	variants []DistrictPlacementVariantListItem,
	validation string,
) ViewModel {
	text := buildAdminSelectText(
		"Изменение цены варианта",
		[]string{
			formatAdminFieldLine("Район", districtName),
			formatAdminFieldLine("Товар", productName),
		},
		validation,
		"Выберите вариант:",
	)

	actions := make([]ActionButton, 0, len(variants)+1)
	for _, variant := range variants {
		variantDisplayLabel := buildAdminQualifiedVariantLabel(productName, variant.Label)

		actions = append(actions, ActionButton{
			ID: adminDistrictVariantSelectVariantAction(variant.ID),
			Label: formatDistrictPlacementVariantActionLabel(
				variantDisplayLabel,
				variant.Price,
				variant.PriceText,
			),
		})
	}
	actions = append(actions, ActionButton{
		ID:    ActionBack,
		Label: "Назад",
	})

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

func buildAdminDistrictVariantPriceUpdateInputView(
	districtName, variantName, currentPriceText, validation string,
) ViewModel {
	text := buildAdminTextWithValidation(
		"Изменение цены варианта",
		[]string{
			formatAdminFieldLine("Район", districtName),
			formatAdminFieldLine("Вариант", variantName),
			formatAdminFieldLine("Текущая цена", currentPriceText),
		},
		validation,
		"Введите новую цену сообщением.",
	)

	return ViewModel{
		Text: text,
		Inline: &InlineKeyboardView{
			Sections: []ActionSection{
				{
					Columns: 1,
					Actions: []ActionButton{
						{ID: ActionBack, Label: "Назад"},
					},
				},
			},
		},
		RemoveReply: true,
	}
}

func currentPlacementPriceTextFromPending(p PendingInput) string {
	raw := p.Value(PendingValueCurrentPrice)
	if raw == "" {
		return ""
	}

	v, err := strconv.Atoi(raw)
	if err != nil || v <= 0 {
		return ""
	}

	return strconv.Itoa(v) + " ₽"
}

func buildAdminDistrictVariantPriceUpdateDoneView() ViewModel {
	return ViewModel{
		Text: "Изменение цены варианта\n\nЦена варианта обновлена.",
		Inline: &InlineKeyboardView{
			Sections: []ActionSection{
				{
					Columns: 1,
					Actions: []ActionButton{
						{ID: ActionBack, Label: "Назад"},
					},
				},
			},
		},
		RemoveReply: true,
	}
}

func buildAdminDistrictVariantPriceUpdateCategorySelectView(
	districtName string,
	categories []CategoryListItem,
	validation string,
) ViewModel {
	text := buildAdminSelectText(
		"Изменение цены варианта",
		[]string{
			formatAdminFieldLine("Район", districtName),
		},
		validation,
		"Выберите категорию:",
	)

	actions := make([]ActionButton, 0, len(categories)+1)
	for _, category := range categories {
		actions = append(actions, ActionButton{
			ID:    adminDistrictVariantPriceUpdateSelectCategoryAction(category.ID),
			Label: category.Label,
		})
	}
	actions = append(actions, ActionButton{
		ID:    ActionBack,
		Label: "Назад",
	})

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

func buildAdminDistrictVariantPriceUpdateProductSelectView(
	districtName, categoryName string,
	products []ProductListItem,
	validation string,
) ViewModel {
	text := buildAdminSelectText(
		"Изменение цены варианта",
		[]string{
			formatAdminFieldLine("Район", districtName),
			formatAdminFieldLine("Категория", categoryName),
		},
		validation,
		"Выберите товар:",
	)

	actions := make([]ActionButton, 0, len(products)+1)
	for _, product := range products {
		actions = append(actions, ActionButton{
			ID:    adminDistrictVariantPriceUpdateSelectProductAction(product.ID),
			Label: product.Label,
		})
	}
	actions = append(actions, ActionButton{
		ID:    ActionBack,
		Label: "Назад",
	})

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

func (s *Service) buildAdminDistrictVariantPriceUpdateDistrictSelectScreen() ViewModel {
	if s == nil || s.districtLister == nil {
		return buildAdminDistrictVariantPriceUpdateDistrictSelectView(nil, "Не удалось загрузить список районов.")
	}

	districts, err := s.districtLister.ListDistricts(context.Background())
	if err != nil {
		return buildAdminDistrictVariantPriceUpdateDistrictSelectView(nil, "Не удалось загрузить список районов.")
	}
	if len(districts) == 0 {
		return buildAdminDistrictVariantPriceUpdateDistrictSelectView(nil, "Нет доступных районов.")
	}

	return buildAdminDistrictVariantPriceUpdateDistrictSelectView(districts, "")
}

func (s *Service) buildAdminDistrictVariantPriceUpdateCategorySelectScreen(districtID int, districtName string) ViewModel {
	if s == nil || s.districtPlacements == nil {
		return buildAdminDistrictVariantPriceUpdateCategorySelectView(
			districtName,
			nil,
			"Не удалось загрузить список категорий.",
		)
	}

	categories, err := s.districtPlacements.ListDistrictCategories(context.Background(), districtID)
	if err != nil {
		return buildAdminDistrictVariantPriceUpdateCategorySelectView(
			districtName,
			nil,
			"Не удалось загрузить список категорий.",
		)
	}
	if len(categories) == 0 {
		return buildAdminDistrictVariantPriceUpdateCategorySelectView(
			districtName,
			nil,
			"В этом районе нет доступных категорий.",
		)
	}

	return buildAdminDistrictVariantPriceUpdateCategorySelectView(districtName, categories, "")
}

func (s *Service) buildAdminDistrictVariantPriceUpdateProductSelectScreen(
	districtID int,
	districtName string,
	categoryID int,
	categoryName string,
) ViewModel {
	if s == nil || s.districtPlacements == nil {
		return buildAdminDistrictVariantPriceUpdateProductSelectView(
			districtName,
			categoryName,
			nil,
			"Не удалось загрузить список товаров.",
		)
	}

	products, err := s.districtPlacements.ListDistrictProducts(context.Background(), districtID, categoryID)
	if err != nil {
		return buildAdminDistrictVariantPriceUpdateProductSelectView(
			districtName,
			categoryName,
			nil,
			"Не удалось загрузить список товаров.",
		)
	}
	if len(products) == 0 {
		return buildAdminDistrictVariantPriceUpdateProductSelectView(
			districtName,
			categoryName,
			nil,
			"В этой категории нет доступных товаров.",
		)
	}

	return buildAdminDistrictVariantPriceUpdateProductSelectView(districtName, categoryName, products, "")
}

func (s *Service) buildAdminDistrictVariantPriceUpdateVariantSelectScreen(
	districtID int,
	districtName string,
	productID int,
	productName string,
) ViewModel {
	if s == nil || s.districtPlacements == nil {
		return buildAdminDistrictVariantPriceUpdateVariantSelectView(
			districtName,
			productName,
			nil,
			"Не удалось загрузить список вариантов.",
		)
	}

	variants, err := s.districtPlacements.ListDistrictVariants(context.Background(), districtID, productID)
	if err != nil {
		return buildAdminDistrictVariantPriceUpdateVariantSelectView(
			districtName,
			productName,
			nil,
			"Не удалось загрузить список вариантов.",
		)
	}
	if len(variants) == 0 {
		return buildAdminDistrictVariantPriceUpdateVariantSelectView(
			districtName,
			productName,
			nil,
			"У этого товара нет доступных вариантов.",
		)
	}

	return buildAdminDistrictVariantPriceUpdateVariantSelectView(districtName, productName, variants, "")
}

func isAdminAction(actionID ActionID) bool {
	if _, ok := parseAdminDistrictSelectCityAction(actionID); ok {
		return true
	}
	if _, ok := parseAdminProductSelectCategoryAction(actionID); ok {
		return true
	}
	if _, ok := parseAdminVariantSelectProductAction(actionID); ok {
		return true
	}
	if _, ok := parseAdminDistrictVariantSelectDistrictAction(actionID); ok {
		return true
	}
	if _, ok := parseAdminDistrictVariantPriceUpdateSelectCategoryAction(actionID); ok {
		return true
	}
	if _, ok := parseAdminDistrictVariantPriceUpdateSelectProductAction(actionID); ok {
		return true
	}
	if _, ok := parseAdminDistrictVariantSelectVariantAction(actionID); ok {
		return true
	}
	if _, ok := parseAdminDistrictVariantSelectCityAction(actionID); ok {
		return true
	}

	switch actionID {
	case ActionAdminOpen,
		ActionAdminCatalogOpen,
		ActionAdminCategoryCreateStart,
		ActionAdminCityCreateStart,
		ActionAdminDistrictCreateStart,
		ActionAdminProductCreateStart,
		ActionAdminVariantCreateStart,
		ActionAdminDistrictVariantCreateStart,
		ActionAdminDistrictVariantPriceUpdateStart:
		return true
	default:
		return false
	}
}

func isAdminScreen(screen ScreenID) bool {
	switch screen {
	case ScreenAdminRoot,
		ScreenAdminCatalog,
		ScreenAdminCategoryCreate,
		ScreenAdminCategoryCode,
		ScreenAdminCategoryCreateDone,
		ScreenAdminCityCreate,
		ScreenAdminCityCode,
		ScreenAdminCityCreateDone,
		ScreenAdminDistrictCitySelect,
		ScreenAdminDistrictCreate,
		ScreenAdminDistrictCode,
		ScreenAdminDistrictCreateDone,
		ScreenAdminProductCategorySelect,
		ScreenAdminProductCreate,
		ScreenAdminProductCode,
		ScreenAdminProductCreateDone,
		ScreenAdminVariantProductSelect,
		ScreenAdminVariantCreate,
		ScreenAdminVariantCode,
		ScreenAdminVariantCreateDone,
		ScreenAdminDistrictVariantCitySelect,
		ScreenAdminDistrictVariantDistrictSelect,
		ScreenAdminDistrictVariantVariantSelect,
		ScreenAdminDistrictVariantPrice,
		ScreenAdminDistrictVariantCreateDone,
		ScreenAdminDistrictVariantPriceUpdateDistrictSelect,
		ScreenAdminDistrictVariantPriceUpdateCategorySelect,
		ScreenAdminDistrictVariantPriceUpdateProductSelect,
		ScreenAdminDistrictVariantPriceUpdateVariantSelect,
		ScreenAdminDistrictVariantPriceUpdatePrice,
		ScreenAdminDistrictVariantPriceUpdateDone:
		return true
	default:
		return false
	}
}

func isAdminPending(kind PendingInputKind) bool {
	switch kind {
	case PendingInputCategoryName,
		PendingInputCategoryCode,
		PendingInputCityName,
		PendingInputCityCode,
		PendingInputDistrictName,
		PendingInputDistrictCode,
		PendingInputProductName,
		PendingInputProductCode,
		PendingInputVariantName,
		PendingInputVariantCode,
		PendingInputDistrictVariantPrice,
		PendingInputDistrictVariantPriceUpdate:
		return true
	default:
		return false
	}
}
