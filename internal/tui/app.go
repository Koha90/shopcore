package tui

import tea "charm.land/bubbletea/v2"

// Run starts TUI application.
func Run(manager BotManager) error {
	program := tea.NewProgram(NewModel(manager, GruvboxTheme()))
	_, err := program.Run()
	return err
}
