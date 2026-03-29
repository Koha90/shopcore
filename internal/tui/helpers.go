package tui

import (
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
