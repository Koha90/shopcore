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
		cfg := m.selectedBotConfig
		if cfg == nil {
			cfg = m.selectedConfig()
		}
		if cfg == nil {
			m.lastErr = fmt.Errorf("config unavailable")
			m.message = "cannot open edit config"
			return m, nil
		}

		m.editForm = BotConfigEditForm{
			Name:       cfg.Name,
			IsEnabled:  cfg.IsEnabled,
			DatabaseID: cfg.DatabaseID,
		}
		m.editCursor = EditFieldName
		m.editDirty = false
		m.screen = ScreenEditBotConfig
		m.lastErr = nil
		m.message = "edit config"
		return m, nil

	case "back":
		m.screen = ScreenList
		m.actionCursor = 0
		return m, nil
	}

	return m, nil
}

// handleEditToggleOrAction handles toggle-like edit actions for current field.
func (m Model) handleEditToggleOrAction() (tea.Model, tea.Cmd) {
	switch m.editCursor {
	case EditFieldEnabled:
		m.editForm.IsEnabled = !m.editForm.IsEnabled
		m.editDirty = true
		return m, nil

	default:
		return m, nil
	}
}

// handleEditEnter handles Enter key in bot config edit screen.
func (m Model) handleEditEnter() (tea.Model, tea.Cmd) {
	switch m.editCursor {
	case EditFieldName:
		return m, nil

	case EditFieldEnabled:
		m.editForm.IsEnabled = !m.editForm.IsEnabled
		m.editDirty = true
		return m, nil

	case EditFieldDatabase:
		m.screen = ScreenSelecteDatabaseProfile
		m.databaseProfiles = nil
		m.databaseCursor = 0
		m.message = "loading database profiles..."
		m.lastErr = nil
		return m, loadDatabaseProfilesCmd(m.config)

	case EditFieldSave:
		m.message = "save config: next step"
		return m, nil

	case EditFieldCancel:
		m.screen = ScreenBotActions
		m.message = "edit cancelled"
		m.lastErr = nil
		return m, nil

	default:
		return m, nil
	}
}
