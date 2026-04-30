package flow

func buildAdminRootView() ViewModel {
	return ViewModel{
		Text: "Админка\n\nВыберите раздел:",
		Inline: &InlineKeyboardView{
			Sections: []ActionSection{
				{
					Columns: 1,
					Actions: []ActionButton{
						{ID: ActionAdminCatalogOpen, Label: "📦 Каталог"},
						backButton(),
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
						{ID: ActionAdminCityCreateStart, Label: "🏙 Создать город"},
						{ID: ActionAdminCategoryCreateStart, Label: "🏷 Создать категорию"},
						{ID: ActionAdminDistrictCreateStart, Label: "📍 Создать район"},

						{ID: ActionAdminProductCreateStart, Label: "🛍 Создать товар"},
						{ID: ActionAdminProductImageUpdateStart, Label: "🖼 Изображение товара"},

						{ID: ActionAdminVariantCreateStart, Label: "🧩 Создать вариант"},
						{ID: ActionAdminVariantImageUpdateStart, Label: "🖼 Изображение варианта"},

						{ID: ActionAdminDistrictVariantCreateStart, Label: "📌 Разместить вариант"},
						{ID: ActionAdminDistrictVariantPriceUpdateStart, Label: "💸 Изменить цену"},

						backButton(),
					},
				},
			},
		},
		RemoveReply: true,
	}
}
