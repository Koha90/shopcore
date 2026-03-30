package flow

// buildCatalogNodeView renders one non-leaf catalog node.
//
// Child nodes are rendered as selectable actions.
// Back navigation is always appended as the last action.
func buildCatalogNodeView(node CatalogNode) ViewModel {
	actions := make([]ActionButton, 0, len(node.Children)+1)

	for _, child := range node.Children {
		actions = append(actions, ActionButton{
			ID:    catalogSelectAction(child.Level, child.ID),
			Label: child.Label,
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
	}
}

// buildCatalogLeafView renders one terminal catalog node.
//
// Leaf nodes display final product or variant information
// and provide only back navigation.
func buildCatalogLeafView(node CatalogNode) ViewModel {
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
					Actions: []ActionButton{
						{ID: ActionBack, Label: "Назад"},
					},
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
