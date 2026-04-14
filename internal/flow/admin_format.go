package flow

import "strconv"

const (
	adminSectionSeparator                  = "\n\n"
	adminFieldLabelSeparator               = ": "
	districtPlacementVariantLabelSeparator = " / "
	adminCurrencySuffix                    = " ₽"
)

func formatAdminFieldLine(label, value string) string {
	if label == "" || value == "" {
		return ""
	}

	return label + adminFieldLabelSeparator + value
}

func formatAdminPriceValue(price int) string {
	if price <= 0 {
		return ""
	}

	return strconv.Itoa(price) + adminCurrencySuffix
}

func formatDistrictPlacementVariantPrice(price int, priceText string) string {
	if priceText != "" {
		return priceText
	}

	return formatAdminPriceValue(price)
}

func formatDistrictPlacementVariantActionLabel(label string, price int, priceText string) string {
	resolvedPriceText := formatDistrictPlacementVariantPrice(price, priceText)

	switch {
	case label == "":
		return resolvedPriceText
	case resolvedPriceText == "":
		return label
	default:
		return label + districtPlacementVariantLabelSeparator + resolvedPriceText
	}
}

func buildAdminText(title string, fields []string, body string) string {
	text := title

	for _, field := range fields {
		if field == "" {
			continue
		}
		text += adminSectionSeparator + field
	}

	if body != "" {
		text += adminSectionSeparator + body
	}

	return text
}

func buildAdminTextWithValidation(title string, fields []string, validation, body string) string {
	if validation == "" {
		return buildAdminText(title, fields, body)
	}

	text := title

	for _, field := range fields {
		if field == "" {
			continue
		}
		text += adminSectionSeparator + field
	}

	text += adminSectionSeparator + validation

	if body != "" {
		text += adminSectionSeparator + body
	}

	return text
}

func formatAdminAutoCodeLine(code string) string {
	if code == "" {
		return ""
	}

	return "Авто-код" + adminFieldLabelSeparator + code
}
