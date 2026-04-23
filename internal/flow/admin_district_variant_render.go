package flow

import "context"

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

func buildAdminDistrictVariantCategorySelectView(
	cityName, districtName string,
	categories []CategoryListItem,
	validation string,
) ViewModel {
	text := buildAdminSelectText(
		"Размещение варианта",
		[]string{
			formatAdminFieldLine("Город", cityName),
			formatAdminFieldLine("Район", districtName),
		},
		validation,
		"Выберите категорию:",
	)

	actions := make([]ActionButton, 0, len(categories)+1)
	for _, category := range categories {
		actions = append(actions, ActionButton{
			ID:    adminDistrictVariantSelectCategoryAction(category.ID),
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

func buildAdminDistrictVariantProductSelectView(
	cityName, districtName, categoryName string,
	products []ProductListItem,
	validation string,
) ViewModel {
	text := buildAdminSelectText(
		"Размещение варианта",
		[]string{
			formatAdminFieldLine("Город", cityName),
			formatAdminFieldLine("Район", districtName),
			formatAdminFieldLine("Категория", categoryName),
		},
		validation,
		"Выберите товар:",
	)

	actions := make([]ActionButton, 0, len(products)+1)
	for _, product := range products {
		actions = append(actions, ActionButton{
			ID:    adminDistrictVariantSelectProductAction(product.ID),
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

func buildAdminDistrictVariantVariantSelectView(
	cityName, districtName, categoryName, productName string,
	variants []VariantListItem,
	validation string,
) ViewModel {
	text := buildAdminSelectText(
		"Размещение варианта",
		[]string{
			formatAdminFieldLine("Город", cityName),
			formatAdminFieldLine("Район", districtName),
			formatAdminFieldLine("Категория", categoryName),
			formatAdminFieldLine("Товар", productName),
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

func buildAdminDistrictVariantPriceInputView(
	districtName, variantName, validation string,
) ViewModel {
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
						{ID: ActionAdminCatalogOpen, Label: "В главное меню"},
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

func (s *Service) buildAdminDistrictVariantCategorySelectScreen(
	cityName, districtName string,
) ViewModel {
	if s == nil || s.categoryLister == nil {
		return buildAdminDistrictVariantCategorySelectView(
			cityName,
			districtName,
			nil,
			"Не удалось загрузить список категорий.",
		)
	}

	categories, err := s.categoryLister.ListCategories(context.Background())
	if err != nil {
		return buildAdminDistrictVariantCategorySelectView(
			cityName,
			districtName,
			nil,
			"Не удалось загрузить список категорий.",
		)
	}
	if len(categories) == 0 {
		return buildAdminDistrictVariantCategorySelectView(
			cityName,
			districtName,
			nil,
			"Нет доступных категорий.",
		)
	}

	return buildAdminDistrictVariantCategorySelectView(cityName, districtName, categories, "")
}

func (s *Service) buildAdminDistrictVariantProductSelectScreen(
	cityName, districtName string,
	categoryID int,
	categoryName string,
) ViewModel {
	if s == nil || s.productLister == nil {
		return buildAdminDistrictVariantProductSelectView(
			cityName,
			districtName,
			categoryName,
			nil,
			"Не удалось загрузить список товаров.",
		)
	}

	products, err := s.productLister.ListProductsByCategory(context.Background(), categoryID)
	if err != nil {
		return buildAdminDistrictVariantProductSelectView(
			cityName,
			districtName,
			categoryName,
			nil,
			"Не удалось загрузить список товаров.",
		)
	}
	if len(products) == 0 {
		return buildAdminDistrictVariantProductSelectView(
			cityName,
			districtName,
			categoryName,
			nil,
			"Нет доступных товаров.",
		)
	}

	return buildAdminDistrictVariantProductSelectView(cityName, districtName, categoryName, products, "")
}

func (s *Service) buildAdminDistrictVariantVariantSelectScreen(
	cityName, districtName, categoryName string,
	districtID, productID int,
	productName string,
) ViewModel {
	if s == nil || s.districtPlacements == nil {
		return buildAdminDistrictVariantVariantSelectView(
			cityName,
			districtName,
			categoryName,
			productName,
			nil,
			"Не удалось загрузить список вариантов.",
		)
	}

	variants, err := s.districtPlacements.ListAvailableVariantsForDistrictProduct(
		context.Background(),
		districtID,
		productID,
	)
	if err != nil {
		return buildAdminDistrictVariantVariantSelectView(
			cityName,
			districtName,
			categoryName,
			productName,
			nil,
			"Не удалось загрузить список вариантов.",
		)
	}
	if len(variants) == 0 {
		return buildAdminDistrictVariantVariantSelectView(
			cityName,
			districtName,
			categoryName,
			productName,
			nil,
			"Не удалось загрузить список вариантов.",
		)
	}

	return buildAdminDistrictVariantVariantSelectView(
		cityName,
		districtName,
		categoryName,
		productName,
		variants,
		"",
	)
}
