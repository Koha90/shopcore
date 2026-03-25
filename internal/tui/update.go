package tui

import (
	"context"
	"errors"
	"fmt"
	"strings"

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

		// Shared text input flow or name/token editing.
		if m.inputMode != InputModeNone {
			var cmd tea.Cmd
			m.textInput, cmd = m.textInput.Update(msg)

			switch msg.String() {
			case "esc":
				m.resetTextInput()
				m.message = ""
				m.lastErr = nil
				return m, nil

			case "enter":
				value := strings.TrimSpace(m.textInput.Value())

				switch m.inputMode {
				case InputModeEditName:
					if value == "" {
						m.message = ""
						m.lastErr = errors.New("name is empty")
						return m, nil
					}

					m.editForm.Name = value
					m.editDirty = true
					m.resetTextInput()
					m.message = "name updated"
					m.lastErr = nil
					return m, nil

				case InputModeEditToken:
					if value == "" {
						m.message = ""
						m.lastErr = errors.New("token is empty")
						return m, nil
					}

					m.message = "saving token..."
					m.lastErr = nil
					token := value
					m.resetTextInput()
					return m, updateBotTokenCmd(m.config, m.selectedID(), token)
				}
			}

			return m, cmd
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
				if m.editDirty {
					m.screen = ScreenConfirmDiscardEdit
					m.confirmCursor = 0
					m.message = "unsaved changes"
					m.lastErr = nil
					return m, nil
				}

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
				if m.editCursor == EditFieldName {
					m.textInput = newTexInput()
					m.inputMode = InputModeEditName
					m.textInput.SetValue(m.editForm.Name)
					m.textInput.Focus()
					m.message = "editing name (enter to apply, esc to cancel)"
					m.lastErr = nil
					return m, nil
				}
				return m.handleEditEnter()
			}

			return m, nil
		}

		if m.screen == ScreenEditBotToken {
			switch msg.String() {
			case "q", "ctrl+c":
				return m, tea.Quit

			case "esc":
				m.resetTextInput()
				m.screen = ScreenBotActions
				m.message = ""
				m.lastErr = nil
				return m, nil

			case "enter":
				if m.inputMode != InputModeEditToken {
					m.textInput = newTexInput()
					m.inputMode = InputModeEditToken
					m.textInput.SetValue("")
					m.textInput.Focus()
					m.message = "editing toke (enter to save, esc to cancel)"
					m.lastErr = nil
					return m, nil
				}
			}

			return m, nil
		}

		if m.screen == ScreenConfirmDiscardEdit {
			switch msg.String() {
			case "q", "ctrl+c":
				return m, tea.Quit

			case "esc":
				m.screen = ScreenEditBotConfig
				m.message = ""
				m.lastErr = nil
				return m, nil

			case "up", "k", "л", "left", "h", "р":
				if m.confirmCursor > 0 {
					m.confirmCursor--
				} else {
					m.confirmCursor = 1
				}
				return m, nil

			case "down", "j", "о", "right", "l", "д":
				if m.confirmCursor < 1 {
					m.confirmCursor++
				} else {
					m.confirmCursor = 0
				}
				return m, nil

			case "enter":
				if m.confirmCursor == 0 {
					m.editDirty = false
					m.screen = ScreenBotActions
					m.message = "changes discarded"
					m.lastErr = nil
					m.editForm = BotConfigEditForm{}
					m.editCursor = EditFieldName
					return m, nil
				}

				m.screen = ScreenEditBotConfig
				m.message = ""
				m.lastErr = nil
				return m, nil
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
				m.screen = ScreenEditBotConfig
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
		m.selectedBotConfigLoading = false

		if msg.err != nil {
			m.selectedBotConfig = nil
			m.lastErr = msg.err
			m.message = "config load failed"
			return m, nil
		}

		cfg := msg.config
		m.selectedBotConfig = &cfg
		m.selectedBotConfigID = cfg.ID
		m.lastErr = nil
		m.message = "config loaded"
		return m, nil

	case editBotConfigLoadedMsg:
		if msg.err != nil {
			m.lastErr = msg.err
			m.message = "cannot open edit config"
			return m, nil
		}

		cfg := msg.config
		m.selectedBotConfig = &cfg
		m.selectedBotConfigID = cfg.ID
		m.selectedBotConfigLoading = false

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

	case botConfigSavedMsg:
		if msg.err != nil {
			m.lastErr = msg.err
			m.message = "config save failed"
			return m, nil
		}

		if err := m.manager.Rename(msg.id, msg.name); err != nil {
			m.lastErr = err
			m.message = "config saved, runtime name sync failed"
			return m, loadBotConfigCmd(m.config, msg.id)
		}

		m.editDirty = false
		m.lastErr = nil
		m.message = "config saved"
		m.screen = ScreenBotConfig
		return m, loadBotConfigCmd(m.config, msg.id)

	case botTokenSavedMsg:
		m.resetTextInput()
		m.lastErr = nil
		m.message = "token saved"
		m.screen = ScreenBotConfig

		return m, loadBotConfigCmd(m.config, msg.id)

	case botTokenSaveFailedMsg:
		m.lastErr = msg.err
		m.message = "token save failed"
		return m, nil
	}

	return m, nil
}
