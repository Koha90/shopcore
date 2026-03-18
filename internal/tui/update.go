package tui

import (
	"context"
	"fmt"

	tea "charm.land/bubbletea/v2"
)

// Update handles TUI messages and user actions.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.layout = detectLayout(msg.Width)
		return m, nil

	case tickMsg:
		m.refresh()
		return m, tickCmd()

	case actionResultMsg:
		m.message = msg.message
		m.lastErr = msg.err
		m.refresh()
		return m, nil

	case tea.MouseWheelMsg:
		if m.screen != ScreenList {
			return m, nil
		}

		mouse := msg.Mouse()
		switch mouse.Button {
		case tea.MouseWheelUp:
			m.scrollUp()
		case tea.MouseWheelDown:
			m.scrollDown()
		}
		return m, nil

	case tea.MouseClickMsg:
		if m.screen != ScreenList {
			return m, nil
		}

		mouse := msg.Mouse()
		if mouse.Button == tea.MouseLeft {
			top := m.listTop()
			if top < 0 {
				return m, nil
			}
			row := mouse.Y - top
			if row >= 0 && row < len(m.visibleBots()) {
				m.cursor = m.offset + row
				m.clampCursor()
				m.ensureCursorVisible()
			}
		}
		return m, nil

	case tea.KeyPressMsg:
		if m.filtering {
			switch msg.String() {
			case "esc", "enter":
				m.filtering = false
				m.refresh()
				return m, nil

			case "backspace":
				if len(m.filter) > 0 {
					m.filter = m.filter[:len(m.filter)-1]
					m.resetListPosition()
					m.refresh()
				}
				return m, nil

			default:
				if msg.Text != "" {
					m.filter += msg.Text
					m.resetListPosition()
					m.refresh()
				}
				return m, nil
			}
		}
		if m.screen == ScreenBotActions {
			switch msg.String() {
			case "q", "ctrl+c":
				return m, tea.Quit

			case "esc":
				m.screen = ScreenList
				m.actionCursor = 0
				return m, nil

			case "up", "k", "л":
				actions := m.botActions()
				if len(actions) == 0 {
					return m, nil
				}
				if m.actionCursor == 0 {
					m.actionCursor = len(actions) - 1
				} else {
					m.actionCursor--
				}
				return m, nil

			case "down", "j", "о":
				actions := m.botActions()
				if len(actions) == 0 {
					return m, nil
				}
				if m.actionCursor >= len(actions)-1 {
					m.actionCursor = 0
				} else {
					m.actionCursor++
				}
				return m, nil

			case "enter":
				return m.handleBotAction()

			}
			return m, nil
		}
		if m.screen == ScreenBotConfig {
			switch msg.String() {
			case "q", "ctrl+c":
				return m, tea.Quit

			case "esc":
				m.screen = ScreenBotActions
				m.lastErr = nil
				m.message = ""
				return m, nil
			}
			return m, nil
		}
		if m.screen == ScreenEditBotConfig {
			switch msg.String() {
			case "q", "ctrl+c":
				return m, tea.Quit

			case "esc":
				m.screen = ScreenBotActions
				m.message = ""
				m.lastErr = nil
				return m, nil

			case "up", "k", "л":
				if m.editCursor > EditFieldName {
					m.editCursor--
				} else {
					m.editCursor = EditFieldCancel
				}
				return m, nil

			case "down", "j", "о":
				if m.editCursor < EditFieldCancel {
					m.editCursor++
				} else {
					m.editCursor = EditFieldName
				}
				return m, nil

			case "left", "h", "р":
				return m.handleEditToggleOrAction()

			case "right", "l", "д":
				return m.handleEditToggleOrAction()

			case " ":
				return m.handleEditToggleOrAction()

			case "enter":
				return m.handleEditEnter()

			case "backspace":
				if m.editCursor == EditFieldName && len(m.editForm.Name) > 0 {
					m.editForm.Name = m.editForm.Name[:len(m.editForm.Name)-1]
					m.editDirty = true
				}
				return m, nil

			default:
				if m.editCursor == EditFieldName && msg.Text != "" {
					m.editForm.Name += msg.Text
					m.editDirty = true
					return m, nil
				}
			}

			return m, nil
		}
		if m.screen == ScreenSelecteDatabaseProfile {
			switch msg.String() {
			case "q", "ctrl+c":
				return m, tea.Quit

			case "esc":
				m.screen = ScreenEditBotConfig
				m.message = ""
				m.lastErr = nil
				return m, nil

			case "up", "k", "л":
				if len(m.databaseProfiles) == 0 {
					return m, nil
				}
				if m.databaseCursor > 0 {
					m.databaseCursor--
				} else {
					m.databaseCursor = len(m.databaseProfiles) - 1
				}
				return m, nil

			case "down", "j", "о":
				if len(m.databaseProfiles) == 0 {
					return m, nil
				}
				if m.databaseCursor < len(m.databaseProfiles)-1 {
					m.databaseCursor++
				} else {
					m.databaseCursor = 0
				}
				return m, nil

			case "enter":
				if len(m.databaseProfiles) == 0 {
					return m, nil
				}

				profile := m.databaseProfiles[m.databaseCursor]
				m.editForm.DatabaseID = profile.ID
				m.editDirty = true
				m.screen = ScreenSelecteDatabaseProfile
				m.message = "database profile selected"
				m.lastErr = nil
				return m, nil
			}

			return m, nil
		}

		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "up", "k", "л":
			m.moveUp()
			return m, nil

		case "down", "j", "о":
			m.moveDown()
			return m, nil

		case "/":
			m.filtering = true
			m.filter = ""
			m.lastErr = nil
			m.message = "filter mode"
			m.resetListPosition()
			m.refresh()
			return m, nil

		case "s", "ы":
			id := m.selectedID()
			if id == "" {
				return m, nil
			}
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

		case "x", "ч":
			id := m.selectedID()
			if id == "" {
				return m, nil
			}
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

		case "r", "к":
			id := m.selectedID()
			if id == "" {
				return m, nil
			}
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

		case "1":
			m.statusFilter = StatusFilterAll
			m.resetListPosition()
			m.refresh()
			return m, nil

		case "2":
			m.statusFilter = StatusFilterRunning
			m.resetListPosition()
			m.refresh()
			return m, nil

		case "3":
			m.statusFilter = StatusFilterStopped
			m.resetListPosition()
			m.refresh()
			return m, nil

		case "4":
			m.statusFilter = StatusFilterFailed
			m.resetListPosition()
			m.refresh()
			return m, nil

		case "enter":
			if m.selectedID() == "" {
				return m, nil
			}
			m.screen = ScreenBotActions
			m.actionCursor = 0
			return m, nil
		}

	case botConfigLoadMsg:
		if msg.err != nil {
			m.selectedBotConfig = nil
			m.lastErr = msg.err
			m.message = "config load failed"
			return m, nil
		}

		cfg := msg.config
		m.selectedBotConfig = &cfg
		m.lastErr = nil
		m.message = "config loaded"
		return m, nil

	case databaseProfilesLoadedMsg:
		if msg.err != nil {
			m.databaseProfiles = nil
			m.lastErr = msg.err
			m.message = "database profiles load failed"
			return m, nil
		}

		m.databaseProfiles = msg.profiles
		m.lastErr = nil
		m.message = "database profiles loaded"
		return m, nil
	}

	return m, nil
}
