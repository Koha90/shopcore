package tui

import tea "charm.land/bubbletea/v2"

// Run starts TUI application.
func Run(manager BotManager, config BotConfigService) error {
	program := tea.NewProgram(
		NewModel(manager, config, GruvboxTheme()),
	)

	_, err := program.Run()
	return err
}
