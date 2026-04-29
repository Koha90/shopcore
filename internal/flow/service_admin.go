package flow

import (
	"context"
)

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
						{ID: ActionAdminCatalogOpen, Label: "В главное меню"},
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

func buildAdminDistrictVariantPriceUpdateDoneView() ViewModel {
	return ViewModel{
		Text: "Изменение цены варианта\n\nЦена варианта обновлена.",
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
