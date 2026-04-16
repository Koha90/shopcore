package flow

// Catalog is in-memory demo catalog used by flow navigation.
//
// It intentionally does not represent final storage or commerce model.
// Its purpose is to validate navigation shape and history-back behavior.
type Catalog struct {
	Schema CatalogSchema
	Roots  []CatalogNode
}

// CatalogNodeMedia describes optional media attached to one catalog node.
//
// Media remains catalog-level data. It does not contain transport-specific
// rendering details. Concrete transports decide how to present the media
// source in their own UI.
//
// ImageSource contains the image location prepared by catalog data builders.
// Empty ImageSource means that the node has no image attached.
//
// ImageAlt contains human-readable fallback text for the image.
// When empty, renderers may fall back to the node label.
//
// Media is optional. Its absence means that the node should be rendered
// without attached media.
type CatalogNodeMedia struct {
	ImageSource string
	ImageAlt    string
}

// CatalogNode describes one catalog tree node prepared for flow rendering.
//
// A node may represent either an intermediate section with children or
// a terminal leaf with final product or variant information.
//
// Media is optional. Its absence means that the node should be rendered
type CatalogNode struct {
	Level       CatalogLevel
	ID          string
	Label       string
	Description string
	PriceText   string
	Children    []CatalogNode

	Media *CatalogNodeMedia
}

func buildCatalogLeafMedia(node CatalogNode) *MediaView {
	if node.Media == nil {
		return nil
	}

	alt := node.Media.ImageAlt
	if alt == "" {
		alt = node.Label
	}

	return NewImageMedia(node.Media.ImageSource, alt)
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

// RootLevel returns the first configured navigation level.
func (c Catalog) RootLevel() (CatalogLevel, bool) {
	return c.Schema.First()
}

// RootNodes returns root-level catalog nodes.
func (c Catalog) RootNodes() []CatalogNode {
	return c.Roots
}

// FindNode resolve one node by full catalog path.
//
// Path must start from root level and follow the tree in the same order
// as configured by the catalog schema.
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
