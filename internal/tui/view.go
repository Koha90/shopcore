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
	summary := m.renderSummary()
	filterBar := m.renderFilterBar()

	left := m.renderList()
	right := m.renderDetails()

	content := lipgloss.JoinHorizontal(
		lipgloss.Top,
		left,
		right,
	)

	status := m.renderStatusBar()
	help := m.theme.Help.Render("↑/↓ or j/k move • / filter • mouse click select • wheel scroll • s start • x stop • r restart • q quit")

	v := tea.NewView(
		lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			summary,
			filterBar,
			"",
			content,
			"",
			status,
			help,
		),
	)
	v.AltScreen = true
	v.MouseMode = tea.MouseModeCellMotion
	return v
}

func (m Model) renderSummary() string {
	parts := []string{
		m.theme.Muted.Render(fmt.Sprintf("total: %d", m.summary.Total)),
		m.theme.Running.Render(fmt.Sprintf("running: %d", m.summary.Running)),
		m.theme.Stopped.Render(fmt.Sprintf("stopped: %d", m.summary.Stopped)),
		m.theme.Failed.Render(fmt.Sprintf("failed: %d", m.summary.Failed)),
	}

	if m.summary.Starting > 0 {
		parts = append(parts, m.theme.Starting.Render(fmt.Sprintf("starting: %d", m.summary.Starting)))
	}
	if m.summary.Stopping > 0 {
		parts = append(parts, m.theme.Stopping.Render(fmt.Sprintf("stopping: %d", m.summary.Stopping)))
	}

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.theme.Border.GetForeground()).
		Padding(0, 1)

	return box.Render(strings.Join(parts, "  •  "))
}

func (m Model) renderFilterBar() string {
	if m.filtering {
		return m.theme.StatusBar.Render("filter: " + m.filter)
	}
	if m.filter != "" {
		return m.theme.Help.Render("filter active: " + m.filter)
	}
	return m.theme.Help.Render("filter: off")
}

func (m Model) renderList() string {
	var lines []string

	header := fmt.Sprintf("Bots (%d/%d)", len(m.filteredBots), len(m.bots))
	lines = append(lines, m.theme.ListHeader.Render(header))

	visible := m.visibleBots()

	if len(visible) == 0 {
		lines = append(lines, m.theme.Muted.Render("no bots match filter"))
	} else {
		for i, bot := range visible {
			absoluteIndex := m.offset + i
			line := fmt.Sprintf("%s  %s", bot.Name, m.renderStatus(bot.Status))
			if absoluteIndex == m.cursor {
				lines = append(lines, m.theme.ListSelected.Render(line))
				continue
			}
			lines = append(lines, m.theme.ListItem.Render(line))
		}
	}

	if len(m.filteredBots) > m.pageSize {
		lines = append(lines, "")
		lines = append(lines, m.theme.Muted.Render(
			fmt.Sprintf("showing %d-%d of %d", m.offset+1, min(m.offset+m.pageSize, len(m.filteredBots)), len(m.filteredBots)),
		))
	}

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.theme.Border.GetForeground()).
		Padding(1, 2).
		Width(42)

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
