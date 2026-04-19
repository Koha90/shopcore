package tui

import tea "charm.land/bubbletea/v2"

// Run starts TUI application.
func Run(manager BotManager, config BotConfigService, logs RuntimeLogReader) error {
	program := tea.NewProgram(
		NewModel(manager, config, logs, GruvboxTheme()),
	)

	_, err := program.Run()
	return err
}
