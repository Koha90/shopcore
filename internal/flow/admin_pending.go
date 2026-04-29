package flow

import "strconv"

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

func pendingCategoryID(p PendingInput) (int, bool) {
	raw := p.Value(PendingValueCategoryID)
	if raw == "" {
		return 0, false
	}

	id, err := strconv.Atoi(raw)
	if err != nil || id <= 0 {
		return 0, false
	}

	return id, true
}

func pendingProductID(p PendingInput) (int, bool) {
	raw := p.Value(PendingValueProductID)
	if raw == "" {
		return 0, false
	}

	id, err := strconv.Atoi(raw)
	if err != nil || id <= 0 {
		return 0, false
	}

	return id, true
}

func pendingDistrictID(p PendingInput) (int, bool) {
	raw := p.Value(PendingValueDistrictID)
	if raw == "" {
		return 0, false
	}

	id, err := strconv.Atoi(raw)
	if err != nil || id <= 0 {
		return 0, false
	}

	return id, true
}

func pendingVariantID(p PendingInput) (int, bool) {
	raw := p.Value(PendingValueVariantID)
	if raw == "" {
		return 0, false
	}

	id, err := strconv.Atoi(raw)
	if err != nil || id <= 0 {
		return 0, false
	}

	return id, true
}

func currentPlacementPriceTextFromPending(p PendingInput) string {
	raw := p.Value(PendingValueCurrentPrice)
	if raw == "" {
		return ""
	}

	v, err := strconv.Atoi(raw)
	if err != nil || v <= 0 {
		return ""
	}

	return strconv.Itoa(v) + " ₽"
}
