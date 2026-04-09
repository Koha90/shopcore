package flow

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	ActionAdminDistrictCreateStart ActionID = "admin_district_create_start"
)

func adminDistrictSelectCityAction(cityID int) ActionID {
	return ActionID(fmt.Sprintf("admin:district:city:%d", cityID))
}

func parseAdminDistrictSelectCityAction(actionID ActionID) (int, bool) {
	const prefix = "admin:district:city:"

	raw := strings.TrimSpace(string(actionID))
	if !strings.HasPrefix(raw, prefix) {
		return 0, false
	}

	idPart := strings.TrimPrefix(raw, prefix)
	if idPart == "" {
		return 0, false
	}

	cityID, err := strconv.Atoi(idPart)
	if err != nil || cityID <= 0 {
		return 0, false
	}

	return cityID, true
}
