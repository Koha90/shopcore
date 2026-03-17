package tui

import (
	"context"
	"fmt"

	tea "charm.land/bubbletea/v2"
)

func (m Model) botActions() []string {
	return []string{
		"start",
		"stop",
		"restart",
		"view config",
		"edit config",
		"back",
	}
}

func (m Model) handleBotAction() (tea.Model, tea.Cmd) {
	id := m.selectedID()
	if id == "" {
		m.screen = ScreenList
		m.actionCursor = 0
		return m, nil
	}

	action := m.botActions()[m.actionCursor]

	switch action {
	case "start":
		return m, func() tea.Msg {
			err := m.manager.Start(context.Background(), id)
			if err != nil {
				return actionResultMsg{
					message: fmt.Sprintf("start %s failed", id),
					err:     err,
				}
			}
			return actionResultMsg{
				message: fmt.Sprintf("started %s", id),
			}
		}

	case "stop":
		return m, func() tea.Msg {
			err := m.manager.Stop(id)
			if err != nil {
				return actionResultMsg{
					message: fmt.Sprintf("stop %s failed", id),
					err:     err,
				}
			}
			return actionResultMsg{
				message: fmt.Sprintf("stopped %s", id),
			}
		}

	case "restart":
		return m, func() tea.Msg {
			err := m.manager.Restart(context.Background(), id)
			if err != nil {
				return actionResultMsg{
					message: fmt.Sprintf("restart %s failed", id),
					err:     err,
				}
			}
			return actionResultMsg{
				message: fmt.Sprintf("restarted %s", id),
			}
		}

	case "view config":
		id := m.selectedID()
		if id == "" {
			return m, nil
		}
		m.screen = ScreenBotConfig
		m.selectedBotConfig = nil
		m.lastErr = nil
		m.message = "loading config..."
		return m, loadBotConfigCmd(m.config, id)

	case "edit config":
		m.message = "edit config: next step"
		return m, nil

	case "back":
		m.screen = ScreenList
		m.actionCursor = 0
		return m, nil
	}

	return m, nil
}
