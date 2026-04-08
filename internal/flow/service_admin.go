package flow

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

func isAdminAction(actionID ActionID) bool {
	switch actionID {
	case ActionAdminOpen, ActionAdminCatalogOpen, ActionAdminCategoryCreateStart, ActionAdminCityCreateStart:
		return true
	default:
		return false
	}
}

func isAdminScreen(screen ScreenID) bool {
	switch screen {
	case ScreenAdminRoot, ScreenAdminCatalog,
		ScreenAdminCategoryCreate, ScreenAdminCategoryCode, ScreenAdminCategoryCreateDone,
		ScreenAdminCityCreate, ScreenAdminCityCode, ScreenAdminCityCreateDone:
		return true
	default:
		return false
	}
}

func isAdminPending(kind PendingInputKind) bool {
	switch kind {
	case PendingInputCategoryName, PendingInputCategoryCode, PendingInputCityCode:
		return true
	default:
		return false
	}
}
