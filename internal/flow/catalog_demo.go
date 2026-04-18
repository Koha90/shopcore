package flow

// DemoCatalog returns in-memory demo catalog tree for flow development.
func DemoCatalog() Catalog {
	return Catalog{
		Schema: DemoCatalogSchema(),
		Roots: []CatalogNode{
			{
				Level: LevelCity,
				ID:    "moscow",
				Label: "Москва",
				Children: []CatalogNode{
					{
						Level: LevelCategory,
						ID:    "flowers",
						Label: "Цветы",
						Children: []CatalogNode{
							{
								Level: LevelDistrict,
								ID:    "center",
								Label: "Центр",
								Children: []CatalogNode{
									{
										Level:       LevelProduct,
										ID:          "rose-box",
										Label:       "Rose Box",
										Description: "Композиция из роз для центрального района.",
										Children: []CatalogNode{
											{
												Level:       LevelVariant,
												ID:          "small",
												Label:       "S / 9 шт",
												PriceText:   "2500 ₽",
												Description: "Компактная упаковка.",
											},
											{
												Level:       LevelVariant,
												ID:          "large",
												Label:       "L / 25 шт",
												PriceText:   "5900 ₽",
												Description: "Большая упаковка.",
												Media: &CatalogNodeMedia{
													ImageSource: "assets/demo/catalog/variant/rose-box-lagrge.jpg",
													ImageAlt:    "Rose Box L / 25 шт",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			{
				Level: LevelCity,
				ID:    "spb",
				Label: "СПб",
				Children: []CatalogNode{
					{
						Level: LevelCategory,
						ID:    "gifts",
						Label: "Подарки",
						Children: []CatalogNode{
							{
								Level: LevelDistrict,
								ID:    "petrogradka",
								Label: "Петроградка",
								Children: []CatalogNode{
									{
										Level:       LevelProduct,
										ID:          "gift-box",
										Label:       "Gift Box",
										Description: "Подарочный набор.",
										Media: &CatalogNodeMedia{
											ImageSource: "assets/demo/catalog/products/gift-box.jpg",
											ImageAlt:    "Gift Box",
										},
										Children: []CatalogNode{
											{
												Level:       LevelVariant,
												ID:          "standard",
												Label:       "Standard",
												PriceText:   "3200 ₽",
												Description: "Базовый вариант набора.",
												Media: &CatalogNodeMedia{
													ImageSource: "assets/demo/catalog/variant/gift-box.jpg",
													ImageAlt:    "Gift Box Standard",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
