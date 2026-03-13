package tui

import (
	"context"
	"fmt"

	tea "charm.land/bubbletea/v2"

	"botmanager/internal/manager"
)

// BotManager defines mangager operations required by TUI.
type BotManager interface {
	List() []manager.Info
	Start(ctx context.Context, id string) error
	Stop(id string) error
	Restart(ctx context.Context, id string) error
	Info(id string) (manager.Info, error)
}

type actionResultMsg struct {
	message string
	err     error
}

// Model represents Bubble Tea application model.
type Model struct {
	manager BotManager
	theme   Theme

	bots    []manager.Info
	cursor  int
	width   int
	height  int
	message string
	lastErr error
}

// NewModel creates new TUI model.
func NewModel(m BotManager, theme Theme) Model {
	model := Model{
		manager: m,
		theme:   theme,
	}
	model.refresh()
	return model
}

// Init initializes TUI model.
func (m Model) Init() tea.Cmd {
	return nil
}

func (m *Model) refresh() {
	m.bots = m.manager.List()

	if len(m.bots) == 0 {
		m.cursor = 0
		return
	}

	if m.cursor >= len(m.bots) {
		m.cursor = len(m.bots) - 1
	}
}

func (m Model) selectedID() string {
	if len(m.bots) == 0 {
		return ""
	}
	return m.bots[m.cursor].ID
}

func (m Model) selectedInfo() *manager.Info {
	if len(m.bots) == 0 {
		return nil
	}
	info := m.bots[m.cursor]
	return &info
}

// Update handles TUI messages and user actions.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case actionResultMsg:
		m.message = msg.message
		m.lastErr = msg.err
		m.refresh()
		return m, nil

	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "up", "k", "л":
			if m.cursor > 0 {
				m.cursor--
			}
			return m, nil

		case "down", "j", "о":
			if m.cursor < len(m.bots)-1 {
				m.cursor++
			}
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
		}
	}

	return m, nil
}
