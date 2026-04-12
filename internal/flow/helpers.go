package flow

func startScreenForScenario(startScenario string) ScreenID {
	switch NormalizeStartScenario(startScenario) {
	case StartScenarioInlineCatalog:
		return ScreenRootExtended
	case StartScenarioReplyWelcome:
		fallthrough
	default:
		return ScreenReplyWelcome
	}
}

func catalogRootForScenario(startScenario string) ScreenID {
	switch NormalizeStartScenario(startScenario) {
	case StartScenarioInlineCatalog:
		return ScreenRootExtended
	default:
		return ScreenRootCompact
	}
}

func normalizeColumns(v int) int {
	if v <= 0 {
		return 1
	}
	return v
}

func catalogChildButtonLabel(node CatalogNode) string {
	if node.Level != LevelVariant {
		return node.Label
	}
	if node.PriceText == "" {
		return node.Label
	}
	return node.Label + " • " + node.PriceText
}
