package tui

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"charm.land/bubbles/v2/textinput"

	"github.com/koha90/shopcore/internal/botconfig"
)

func newTextInput() textinput.Model {
	ti := textinput.New()
	ti.Placeholder = ""
	ti.CharLimit = 512
	ti.SetWidth(48)
	return ti
}

func (m *Model) resetTextInput() {
	m.textInput = newTextInput()
	m.inputMode = InputModeNone
}

func nextStartScenario(current string) string {
	switch current {
	case botconfig.StartScenarioInlineCatalog:
		return botconfig.StartScenarioReplyWelcome
	default:
		return botconfig.StartScenarioInlineCatalog
	}
}

func formatTelegramAdminUserIDs(ids []int64) string {
	if len(ids) == 0 {
		return ""
	}

	parts := make([]string, 0, len(ids))
	for _, id := range ids {
		parts = append(parts, strconv.FormatInt(id, 10))
	}

	return strings.Join(parts, ", ")
}

func parseTelegramAdminUserIDs(input string) ([]int64, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, nil
	}

	parts := strings.Split(input, ",")
	seen := make(map[int64]struct{}, len(parts))
	result := make([]int64, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		id, err := strconv.ParseInt(part, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid telegram admin user id %q", part)
		}
		if id <= 0 {
			return nil, fmt.Errorf("invalid telegram admin user id %d", id)
		}
		if _, ok := seen[id]; ok {
			continue
		}

		seen[id] = struct{}{}
		result = append(result, id)
	}

	slices.Sort(result)

	if len(result) == 0 {
		return nil, nil
	}

	return result, nil
}
