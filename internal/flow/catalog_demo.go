package flow

// Catalog is in-memory demo catalog used by flow navigation.
//
// It intentionally does not represent final storage or commerce model.
// Its purpose is to validate navigation shape and history-back behavior.
type Catalog struct {
	Schema CatalogSchema
	Roots  []CatalogNode
}

// CatalogNode is one selectable node in catalog tree.
type CatalogNode struct {
	Level       CatalogLevel
	ID          string
	Label       string
	Description string
	PriceText   string
	Children    []CatalogNode
}

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
										Children: []CatalogNode{
											{
												Level:       LevelVariant,
												ID:          "standard",
												Label:       "Standard",
												PriceText:   "3200 ₽",
												Description: "Базовый вариант набора.",
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

// RootLevel returns the first configured navigation level.
func (c Catalog) RootLevel() (CatalogLevel, bool) {
	return c.Schema.First()
}

// RootNodes returns root-level catalog nodes.
func (c Catalog) RootNodes() []CatalogNode {
	return c.Roots
}

// FindNode resolve one node by full catalog path.
func (c Catalog) FindNode(path CatalogPath) (CatalogNode, bool) {
	if len(path) == 0 {
		return CatalogNode{}, false
	}

	nodes := c.Roots
	var found CatalogNode

	for _, sel := range path {
		ok := false
		for _, node := range nodes {
			if node.Level == sel.Level && node.ID == sel.ID {
				found = node
				nodes = node.Children
				ok = true
				break
			}
		}
		if !ok {
			return CatalogNode{}, false
		}
	}

	return found, true
}
