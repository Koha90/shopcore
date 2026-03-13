package tui

import (
	"context"
	"fmt"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"

	"botmanager/internal/manager"
)

const defaultPageSize = 10

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

type tickMsg time.Time

// Summary contains aggregated runtime counters for bots.
type Summary struct {
	Total    int
	Running  int
	Stopped  int
	Failed   int
	Starting int
	Stopping int
}

// Model represents Bubble Tea application model.
type Model struct {
	manager BotManager
	theme   Theme

	bots         []manager.Info
	filteredBots []manager.Info
	summary      Summary

	cursor   int
	offset   int
	pageSize int

	width   int
	height  int
	message string
	lastErr error

	filtering bool
	filter    string
}

// NewModel creates new TUI model.
func NewModel(m BotManager, theme Theme) Model {
	model := Model{
		manager:  m,
		theme:    theme,
		pageSize: defaultPageSize,
	}
	model.refresh()
	return model
}

// Init initializes TUI model.
func (m Model) Init() tea.Cmd {
	return tickCmd()
}

func tickCmd() tea.Cmd {
	return tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m *Model) refresh() {
	m.bots = m.manager.List()
	m.summary = buildSummary(m.bots)
	m.applyFilter()
	m.clampCursor()
	m.ensureCursorVisible()
}

func buildSummary(bots []manager.Info) Summary {
	var s Summary
	s.Total = len(bots)

	for _, bot := range bots {
		switch bot.Status {
		case manager.StatusRunning:
			s.Running++
		case manager.StatusStopped:
			s.Stopped++
		case manager.StatusFailed:
			s.Failed++
		case manager.StatusStarting:
			s.Starting++
		case manager.StatusStopping:
			s.Stopping++
		}
	}

	return s
}

func (m *Model) applyFilter() {
	if strings.TrimSpace(m.filter) == "" {
		m.filteredBots = append([]manager.Info(nil), m.bots...)
		return
	}

	q := strings.ToLower(strings.TrimSpace(m.filter))
	result := make([]manager.Info, 0, len(m.bots))

	for _, bot := range m.bots {
		if strings.Contains(strings.ToLower(bot.ID), q) ||
			strings.Contains(strings.ToLower(bot.Name), q) ||
			strings.Contains(strings.ToLower(string(bot.Status)), q) {
			result = append(result, bot)
		}
	}

	m.filteredBots = result
}

func (m *Model) resetListPosition() {
	m.cursor = 0
	m.offset = 0
}

func (m *Model) clampCursor() {
	if len(m.filteredBots) == 0 {
		m.cursor = 0
		m.offset = 0
		return
	}

	if m.cursor < 0 {
		m.cursor = 0
	}
	if m.cursor >= len(m.filteredBots) {
		m.cursor = len(m.filteredBots) - 1
	}
}

func (m *Model) ensureCursorVisible() {
	if len(m.filteredBots) == 0 {
		m.offset = 0
		return
	}

	if m.cursor < m.offset {
		m.offset = m.cursor
	}

	if m.cursor >= m.offset+m.pageSize {
		m.offset = m.cursor - m.pageSize + 1
	}

	if m.offset < 0 {
		m.offset = 0
	}

	maxOffset := len(m.filteredBots) - m.pageSize
	if maxOffset < 0 {
		maxOffset = 0
	}
	if m.offset > maxOffset {
		m.offset = maxOffset
	}
}

func (m Model) visibleBots() []manager.Info {
	if len(m.filteredBots) == 0 {
		return nil
	}

	start := m.offset
	end := start + m.pageSize
	if end > len(m.filteredBots) {
		end = len(m.filteredBots)
	}

	return m.filteredBots[start:end]
}

func (m Model) selectedID() string {
	if len(m.filteredBots) == 0 {
		return ""
	}
	return m.filteredBots[m.cursor].ID
}

func (m Model) selectedInfo() *manager.Info {
	if len(m.filteredBots) == 0 {
		return nil
	}
	info := m.filteredBots[m.cursor]
	return &info
}

func (m *Model) moveUp() {
	if len(m.filteredBots) == 0 {
		return
	}

	if m.cursor == 0 {
		m.cursor = len(m.filteredBots) - 1
	} else {
		m.cursor--
	}
	m.ensureCursorVisible()
}

func (m *Model) moveDown() {
	if len(m.filteredBots) == 0 {
		return
	}

	if m.cursor >= len(m.filteredBots)-1 {
		m.cursor = 0
	} else {
		m.cursor++
	}
	m.ensureCursorVisible()
}

func (m *Model) scrollUp() {
	if m.offset > 0 {
		m.offset--
	}
	if m.cursor < m.offset {
		m.cursor = m.offset
	}
}

func (m *Model) scrollDown() {
	maxOffset := len(m.filteredBots) - m.pageSize
	if maxOffset < 0 {
		maxOffset = 0
	}
	if m.offset < maxOffset {
		m.offset++
	}
	if m.cursor < m.offset {
		m.cursor = m.offset
	}
	if m.cursor >= m.offset+m.pageSize {
		m.cursor = m.offset + m.pageSize - 1
	}
	if m.cursor >= len(m.filteredBots) {
		m.cursor = len(m.filteredBots) - 1
	}
}

// Update handles TUI messages and user actions.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
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
		mouse := msg.Mouse()
		switch mouse.Button {
		case tea.MouseWheelUp:
			m.scrollUp()
		case tea.MouseWheelDown:
			m.scrollDown()
		}
		return m, nil

	case tea.MouseClickMsg:
		mouse := msg.Mouse()
		if mouse.Button == tea.MouseLeft {
			listTop := 9
			row := mouse.Y - listTop
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
		}
	}

	return m, nil
}
