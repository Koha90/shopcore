package flow

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	ActionAdminVariantCreateStart ActionID = "admin_variant_create_start"
)

func adminVariantSelectProductAction(productID int) ActionID {
	return ActionID(fmt.Sprintf("admin:variant:product:%d", productID))
}

func parseAdminVariantSelectProductAction(actionID ActionID) (int, bool) {
	const prefix = "admin:variant:product:"

	raw := strings.TrimSpace(string(actionID))
	if !strings.HasPrefix(raw, prefix) {
		return 0, false
	}

	idPart := strings.TrimPrefix(raw, prefix)
	if idPart == "" {
		return 0, false
	}

	productID, err := strconv.Atoi(idPart)
	if err != nil || productID <= 0 {
		return 0, false
	}

	return productID, true
}
