package flow

import (
	"strconv"
	"strings"
)

const (
	adminSectionSeparator                  = "\n\n"
	adminFieldLabelSeparator               = ": "
	districtPlacementVariantLabelSeparator = " - "
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
	var b strings.Builder

	b.WriteString(title)
	writeAdminSections(&b, fields)

	if body != "" {
		b.WriteString(adminSectionSeparator)
		b.WriteString(body)
	}

	return b.String()
}

func buildAdminTextWithValidation(title string, fields []string, validation, body string) string {
	var b strings.Builder

	b.WriteString(title)
	writeAdminSections(&b, fields)

	if validation != "" {
		b.WriteString(adminSectionSeparator)
		b.WriteString(validation)
	}

	if body != "" {
		b.WriteString(adminSectionSeparator)
		b.WriteString(body)
	}

	return b.String()
}

func formatAdminAutoCodeLine(code string) string {
	if code == "" {
		return ""
	}

	return "Авто-код" + adminFieldLabelSeparator + code
}

func writeAdminSections(b *strings.Builder, sections []string) {
	for _, section := range sections {
		if section == "" {
			continue
		}

		b.WriteString(adminSectionSeparator)
		b.WriteString(section)
	}
}

func buildAdminSelectText(title string, fields []string, validation, prompt string) string {
	return buildAdminTextWithValidation(title, fields, validation, prompt)
}
