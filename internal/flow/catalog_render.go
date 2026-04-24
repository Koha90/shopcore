package flow

// Catalog is in-memory demo catalog used by flow navigation.
//
// It intentionally does not represent final storage or commerce model.
// Its purpose is to validate navigation shape and history-back behavior.
type Catalog struct {
	Schema CatalogSchema
	Roots  []CatalogNode
}

// buildCatalogNodeView renders one non-leaf catalog node.
//
// Child nodes are rendered as selectable actions.
// Back navigation is always appended as the last action.
func buildCatalogNodeView(node CatalogNode) ViewModel {
	actions := make([]ActionButton, 0, len(node.Children)+1)

	for _, child := range node.Children {
		label := child.Label
		if child.Level == LevelVariant && child.PriceText != "" {
			label += " - " + child.PriceText
		}

		actions = append(actions, ActionButton{
			ID:    catalogSelectAction(child.Level, child.ID),
			Label: label,
		})
	}

	actions = append(actions, ActionButton{
		ID:    ActionBack,
		Label: "Назад",
	})

	text := node.Label
	if node.Description != "" {
		text += "\n\n" + node.Description
	}

	if prompt := levelPromptForChildren(node.Children); prompt != "" {
		text += "\n\n" + prompt
	}

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
		Media:       buildCatalogNodeMediaView(node),
	}
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

func buildCatalogNodeMediaView(node CatalogNode) *MediaView {
	if node.Media == nil {
		return nil
	}

	alt := node.Media.ImageAlt
	if alt == "" {
		alt = node.Label
	}

	return NewImageMedia(node.Media.ImageSource, alt)
}

// buildCatalogLeafView renders one terminal catalog node.
//
// Variant leaf may start order flow.
// Other leaf kinds, if they appear in the future, keep regular back navigation.
func buildCatalogLeafView(node CatalogNode) ViewModel {
	actions := make([]ActionButton, 0, 2)

	if node.Level == LevelVariant {
		actions = append(actions, ActionButton{
			ID:    ActionOrderStart,
			Label: "Заказать",
		})
	}

	actions = append(actions, ActionButton{
		ID:    ActionBack,
		Label: "Назад",
	})

	text := node.Label
	if node.PriceText != "" {
		text += "\n\n" + node.PriceText
	}
	if node.Description != "" {
		text += "\n\n" + node.Description
	}

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

// levelPromptForChildren returns human-friendly prompt text
// based on the next child level.
func levelPromptForChildren(children []CatalogNode) string {
	if len(children) == 0 {
		return ""
	}

	switch children[0].Level {
	case LevelCity:
		return "Выберите город:"
	case LevelCategory:
		return "Выберите категорию:"
	case LevelDistrict:
		return "Выберите район:"
	case LevelProduct:
		return "Выберите товар:"
	case LevelVariant:
		return "Выберите вариант:"
	default:
		return "Выберите раздел:"
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
