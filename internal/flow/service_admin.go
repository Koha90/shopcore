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
						{ID: ActionBack, Label: "Назад"},
					},
				},
			},
		},
		RemoveReply: true,
	}
}

func buildAdminCategoryCreateInputView(validation string) ViewModel {
	text := "Новая категория\n\nВведите название категории сообщением."
	if validation != "" {
		text = "Новая категория\n\n" + validation + "\n\nВведите название категории сообщением."
	}

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
	text := "Новая категория\n\nВведите code категории сообщением."
	if suggested != "" {
		text = "Новая категория\n\nАвто-код: " + suggested + "\n\nВведите code категории сообщением."
	}
	if validation != "" {
		text = "Новая категория\n\n" + validation + "\n\nВведите code категории сообщением."
	}

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
	text := "Новый город\n\nВведите название города сообщением."
	if validation != "" {
		text = "Новый город\n\n" + validation + "\n\nВведите название города сообщением."
	}

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
	text := "Новый город\n\nВведите code города сообщением."
	if suggested != "" {
		text = "Новый город\n\nАвто-код: " + suggested + "\n\nВведите code города сообщением."
	}
	if validation != "" {
		text = "Новый город\n\n" + validation + "\n\nВведите code города сообщением."
	}

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
	text := "Новый район\n\nВыберите город:"
	if validation != "" {
		text = "Новый район\n\n" + validation + "\n\nВыберите город:"
	}

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
	text := "Новый район"
	if cityName != "" {
		text += "\n\nГород: " + cityName
	}
	text += "\n\nВведите название района сообщением."
	if validation != "" {
		text = "Новый район"
		if cityName != "" {
			text += "\n\nГород: " + cityName
		}
		text += "\n\n" + validation + "\n\nВведите название района сообщением."
	}

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
	text := "Новый район"
	if cityName != "" {
		text += "\n\nГород: " + cityName
	}
	text += "\n\nВведите code района сообщением."

	if suggested != "" {
		text = "Новый район"
		if cityName != "" {
			text += "\n\nГород: " + cityName
		}
		text += "\n\nАвто-код: " + suggested + "\n\nВведите code района сообщением."
	}

	if validation != "" {
		text = "Новый район"
		if cityName != "" {
			text += "\n\nГород: " + cityName
		}
		text += "\n\n" + validation + "\n\nВведите code района сообщением."
	}

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

func isAdminAction(actionID ActionID) bool {
	if _, ok := parseAdminDistrictSelectCityAction(actionID); ok {
		return true
	}

	switch actionID {
	case ActionAdminOpen,
		ActionAdminCatalogOpen,
		ActionAdminCategoryCreateStart,
		ActionAdminCityCreateStart,
		ActionAdminDistrictCreateStart:
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
		ScreenAdminDistrictCreateDone:
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
		PendingInputDistrictCode:
		return true
	default:
		return false
	}
}
