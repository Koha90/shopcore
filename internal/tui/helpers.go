package tui

import (
	"charm.land/bubbles/v2/textinput"
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
	case "inline_catalog":
		return "reply_catalog"
	default:
		return "inline_catalog"
	}
}
