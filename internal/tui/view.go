package tui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"botmanager/internal/manager"
)

// View renders full TUI screen.
func (m Model) View() tea.View {
	title := m.theme.Title.Render("botmanager · runtime panel")

	left := m.renderList()
	right := m.renderDetails()

	content := lipgloss.JoinHorizontal(
		lipgloss.Top,
		left,
		right,
	)

	status := m.renderStatusBar()
	help := m.theme.Help.Render("↑/↓ or j/k move • s start • x stop • r restart • q quit")

	v := tea.NewView(
		lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			"",
			content,
			"",
			status,
			help,
		),
	)
	v.AltScreen = true
	return v
}

func (m Model) renderList() string {
	var lines []string

	lines = append(lines, m.theme.ListHeader.Render("Bots"))

	if len(m.bots) == 0 {
		lines = append(lines, m.theme.Muted.Render("no bots registered"))
	} else {
		for i, bot := range m.bots {
			line := fmt.Sprintf("%s  %s", bot.Name, m.renderStatus(bot.Status))
			if i == m.cursor {
				lines = append(lines, m.theme.ListSelected.Render(line))
				continue
			}
			lines = append(lines, m.theme.ListItem.Render(line))
		}
	}

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.theme.Border.GetBackground()).
		Padding(1, 2).
		Width(36)

	return box.Render(strings.Join(lines, "\n"))
}

func (m Model) renderDetails() string {
	info := m.selectedInfo()

	var lines []string
	lines = append(lines, m.theme.ListHeader.Render("Details"))

	if info == nil {
		lines = append(lines, m.theme.Muted.Render("nothing selected"))
	} else {
		lines = append(lines, fmt.Sprintf("ID: %s", info.ID))
		lines = append(lines, fmt.Sprintf("Name: %s", info.Name))
		lines = append(lines, fmt.Sprintf("Status: %s", m.renderStatus(info.Status)))

		if info.LastError == "" {
			lines = append(lines, "Last error: none")
		} else {
			lines = append(lines, m.theme.Error.Render("Last error: "+info.LastError))
		}
	}

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.theme.Border.GetForeground()).
		Padding(1, 2).
		Width(52)

	return box.Render(strings.Join(lines, "\n"))
}

func (m Model) renderStatusBar() string {
	if m.lastErr != nil {
		return m.theme.Error.Render("error: " + m.lastErr.Error())
	}
	if m.message != "" {
		return m.theme.Success.Render(m.message)
	}
	return m.theme.StatusBar.Render("ready")
}

func (m Model) renderStatus(status manager.Status) string {
	switch status {
	case manager.StatusRunning:
		return m.theme.Running.Render(string(status))
	case manager.StatusStopped:
		return m.theme.Stopped.Render(string(status))
	case manager.StatusFailed:
		return m.theme.Failed.Render(string(status))
	case manager.StatusStarting:
		return m.theme.Starting.Render(string(status))
	case manager.StatusStopping:
		return m.theme.Stopping.Render(string(status))
	default:
		return m.theme.Muted.Render(string(status))
	}
}
