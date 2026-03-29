package tui

import (
	"context"
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"

	"github.com/koha90/shopcore/internal/botconfig"
)

func (m Model) botActions() []string {
	return []string{
		"start",
		"stop",
		"restart",
		"view config",
		"edit config",
		"edit token",
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
		id = m.selectedID()
		if id == "" {
			return m, nil
		}
		m.screen = ScreenBotConfig
		m.selectedBotConfig = nil
		m.selectedBotConfigID = id
		m.selectedBotConfigLoading = true
		m.lastErr = nil
		m.message = "loading config..."
		return m, loadBotConfigCmd(m.config, id)

	case "edit config":
		if id == "" {
			return m, nil
		}

		m.lastErr = nil
		m.message = "loading config for edit..."
		return m, loadEditBotConfigCmd(m.config, id)

	case "edit token":
		m.screen = ScreenEditBotToken
		m.message = ""
		m.lastErr = nil
		return m, loadBotConfigCmd(m.config, m.selectedID())

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
		m.message = "enabled update"
		m.lastErr = nil
		return m, nil

	case EditFieldDatabase:
		m.screen = ScreenSelecteDatabaseProfile
		m.databaseCursor = 0
		m.message = ""
		m.lastErr = nil
		return m, loadDatabaseProfilesCmd(m.config)

	case EditFieldStartScenario:
		m.editForm.StartScenario = nextStartScenario(m.editForm.StartScenario)
		m.editDirty = true
		m.message = "scenario updated"
		m.lastErr = nil
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
		m.message = "enabled update"
		m.lastErr = nil
		return m, nil

	case EditFieldDatabase:
		m.screen = ScreenSelecteDatabaseProfile
		m.databaseProfiles = nil
		m.databaseCursor = 0
		m.message = "loading database profiles..."
		m.lastErr = nil
		return m, loadDatabaseProfilesCmd(m.config)

	case EditFieldStartScenario:
		m.editForm.StartScenario = nextStartScenario(m.editForm.StartScenario)
		m.editDirty = true
		m.message = "scenario updated"
		m.lastErr = nil
		return m, nil

	case EditFieldSave:
		id := m.selectedID()
		if id == "" {
			m.lastErr = fmt.Errorf("bot id is empty")
			m.message = "save failed"
			return m, nil
		}

		if strings.TrimSpace(m.editForm.Name) == "" {
			m.lastErr = fmt.Errorf("name cannot be empty")
			m.message = "validation error"
			return m, nil
		}

		if strings.TrimSpace(m.editForm.DatabaseID) == "" {
			m.lastErr = fmt.Errorf("database id cannot be empty")
			m.message = "validation error"
			return m, nil
		}

		if strings.TrimSpace(m.editForm.StartScenario) == "" {
			m.lastErr = fmt.Errorf("start scenario cannot be empty")
			m.message = "validation error"
			return m, nil
		}

		if !botconfig.IsValidStartScenario(m.editForm.StartScenario) {
			m.lastErr = fmt.Errorf("start scenario is invalid")
			m.message = "validation error"
			return m, nil
		}

		m.message = "saving config..."
		m.lastErr = nil
		return m, saveBotConfigCmd(m.config, id, m.editForm)

	case EditFieldCancel:
		m.screen = ScreenBotActions
		m.message = "edit cancelled"
		m.lastErr = nil
		return m, nil

	default:
		return m, nil
	}
}
