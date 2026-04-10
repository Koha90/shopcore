package flow

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	ActionAdminProductCreateStart ActionID = "admin_product_create_start"
)

func adminProductSelectCategoryAction(categoryID int) ActionID {
	return ActionID(fmt.Sprintf("admin:product:category:%d", categoryID))
}

func parseAdminProductSelectCategoryAction(actionID ActionID) (int, bool) {
	const prefix = "admin:product:category:"

	raw := strings.TrimSpace(string(actionID))
	if !strings.HasPrefix(raw, prefix) {
		return 0, false
	}

	idPart := strings.TrimPrefix(raw, prefix)
	if idPart == "" {
		return 0, false
	}

	categoryID, err := strconv.Atoi(idPart)
	if err != nil || categoryID <= 0 {
		return 0, false
	}

	return categoryID, true
}
