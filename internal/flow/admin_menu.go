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
						{ID: ActionAdminProductImageUpdateStart, Label: "Изменить изображение товара"},
						{ID: ActionAdminVariantCreateStart, Label: "Создать вариант"},
						{ID: ActionAdminVariantImageUpdateStart, Label: "Изменить изображение варианта"},
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
