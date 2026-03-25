package tui

func (m *Model) resetTextInput() {
	m.textInput = newTexInput()
	m.inputMode = InputModeNone
}
