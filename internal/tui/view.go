package tui

import (
	"fmt"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"botmanager/internal/manager"
)

// View renders full TUI screen.
func (m Model) View() tea.View {
	if m.layout == LayoutMobile {
		return m.viewMobile()
	}
	return m.viewDesktop()
}

func (m Model) viewDesktop() tea.View {
	title := m.theme.Title.Render("platform admin · telegram runtime")
	status := m.renderStatusBar()

	var body string
	var help string

	switch m.screen {
	case ScreenBotActions:
		body = m.renderBotActions()
		help = m.theme.Help.Render("↑/↓ or j/k move • enter select • esc back • q quit")

	case ScreenBotConfig:
		body = m.renderBotConfig()
		help = m.theme.Help.Render("esc back • q quit")

	case ScreenEditBotConfig:
		body = m.renderEditBotConfig()
		help = m.theme.Help.Render("↑/↓ or j/k move • type name • space toggle • enter select • esc back • q quit")

	case ScreenSelecteDatabaseProfile:
		body = m.renderDatabaseProfileSelect()
		help = m.theme.Help.Render("j/k move • enter select • esc back • q quit")

	default:
		summary := m.renderSummary()
		filterBar := m.renderFilterBar()

		content := lipgloss.JoinHorizontal(
			lipgloss.Top,
			m.renderList(),
			m.renderDetails(),
		)

		body = lipgloss.JoinVertical(
			lipgloss.Left,
			summary,
			filterBar,
			"",
			content,
		)
		help = m.theme.Help.Render(
			"↑/↓ or j/k move • enter options • / filter • mouse click select • wheel scroll • s start • x stop • r restart • q quit",
		)
	}

	v := tea.NewView(
		lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			body,
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

func (m Model) renderMobileSummary() string {
	line1 := strings.Join([]string{
		m.theme.Muted.Render(fmt.Sprintf("%d total", m.summary.Total)),
		m.theme.Running.Render(fmt.Sprintf("%d run", m.summary.Running)),
		m.theme.Stopped.Render(fmt.Sprintf("%d stop", m.summary.Stopped)),
	}, "  •  ")

	line2 := strings.Join([]string{
		m.theme.Failed.Render(fmt.Sprintf("%d fail", m.summary.Failed)),
		m.theme.Starting.Render(fmt.Sprintf("%d start", m.summary.Starting)),
		m.theme.Stopping.Render(fmt.Sprintf("%d stoping", m.summary.Stopping)),
	}, "  •  ")

	filter := m.theme.Help.Render(
		fmt.Sprintf("filter: %s  |  1 all  2 run  3 stop  4 fail", m.statusFilter),
	)

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.theme.Border.GetForeground()).
		Padding(0, 1)

	return box.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			line1,
			line2,
			filter,
		),
	)
}

func (m Model) viewMobile() tea.View {
	title := m.theme.Title.Render("platform admin")
	status := m.renderStatusBar()

	var body string
	var help string

	switch m.screen {
	case ScreenBotActions:
		body = m.renderBotActions()
		help = m.theme.Help.Render("j/k move • enter select • esc back • q quit")

	case ScreenBotConfig:
		body = m.renderBotConfig()
		help = m.theme.Help.Render("esc back • q quit")

	case ScreenEditBotConfig:
		body = m.renderEditBotConfig()
		help = m.theme.Help.Render("j/k move • type name • space toggle • enter select • esc back • q quit")

	case ScreenSelecteDatabaseProfile:
		body = m.renderDatabaseProfileSelect()
		help = m.theme.Help.Render("j/k move • enter select • esc back • q quit")

	default:
		summary := m.renderMobileSummary()
		filterBar := m.renderFilterBar()
		content := m.renderList()

		body = lipgloss.JoinVertical(
			lipgloss.Left,
			summary,
			filterBar,
			"",
			content,
		)
		help = m.theme.Help.Render("j/k move • enter options • / filter • 1/2/3/4 status • q quit")
	}

	v := tea.NewView(
		lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			body,
			"",
			status,
			help,
		),
	)
	v.AltScreen = true
	v.MouseMode = tea.MouseModeCellMotion
	return v
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

	const labelWidth = 25

	if len(visible) == 0 {
		lines = append(lines, m.theme.Muted.Render("no bots match filter"))
	} else {
		for i, bot := range visible {
			absoluteIndex := m.offset + i
			line := renderKeyValue(labelWidth, bot.Name, m.renderStatus(bot.Status))

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
	cfg := m.selectedConfig()

	var lines []string
	lines = append(lines, m.theme.ListHeader.Render("Details"))

	const labelWidth = 12

	if info == nil {
		lines = append(lines, m.theme.Muted.Render("nothing selected"))
	} else {
		lines = append(lines, m.theme.ListHeader.Render("Runtime"))
		lines = append(lines, renderKeyValue(labelWidth, "ID", info.ID))
		lines = append(lines, renderKeyValue(labelWidth, "Name", info.Name))
		lines = append(lines, renderKeyValue(labelWidth, "Status", m.renderStatus(info.Status)))

		if info.LastError == "" {
			lines = append(lines, "Last error: none")
		} else {
			lines = append(lines, m.theme.Error.Render("Last error: "+info.LastError))
		}

		lines = append(lines, "")
		lines = append(lines, m.theme.ListHeader.Render("Config"))

		if cfg == nil {
			lines = append(lines, m.theme.Muted.Render("config unavailable"))
		} else {
			lines = append(lines, renderKeyValue(labelWidth, "Config name", cfg.Name))
			lines = append(lines, renderKeyValue(labelWidth, "Token", cfg.TokenMasked))
			lines = append(lines, renderKeyValue(labelWidth, "Database", cfg.DatabaseName))
			lines = append(lines, renderKeyValue(labelWidth, "Enabled", fmt.Sprintf("%t", cfg.IsEnabled)))
		}

	}

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.theme.Border.GetForeground()).
		Padding(1, 2).
		Width(56)

	return box.Render(strings.Join(lines, "\n"))
}

func (m Model) renderBotActions() string {
	info := m.selectedInfo()

	var lines []string
	lines = append(lines, m.theme.ListHeader.Render("Bot options"))

	const labelWidth = 12

	if info == nil {
		lines = append(lines, m.theme.Muted.Render("nothing selected"))
	} else {
		lines = append(lines, renderKeyValue(labelWidth, "Name", info.Name))
		lines = append(lines, renderKeyValue(labelWidth, "Status", m.renderStatus(info.Status)))
		lines = append(lines, "")
		lines = append(lines, m.theme.ListHeader.Render("Actions"))

		for i, action := range m.botActions() {
			line := action
			if i == m.actionCursor {
				lines = append(lines, m.theme.ListSelected.Render("> "+line))
			} else {
				lines = append(lines, m.theme.ListItem.Render("  "+line))
			}
		}
	}

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.theme.Border.GetForeground()).
		Padding(1, 2)

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

func (m Model) renderBotConfig() string {
	var lines []string
	lines = append(lines, m.theme.ListHeader.Render("Bot config"))

	const labelWidth int = 12

	if m.selectedBotConfig == nil {
		lines = append(lines, m.theme.Muted.Render("loading or unavailable"))
	} else {
		cfg := m.selectedBotConfig
		lines = append(lines, renderKeyValue(labelWidth, "ID", cfg.ID)) // fmt.Sprintf("ID: %s", cfg.ID))
		lines = append(lines, renderKeyValue(labelWidth, "Name", cfg.Name))
		lines = append(lines, renderKeyValue(labelWidth, "Token", cfg.TokenMasked))
		lines = append(lines, renderKeyValue(labelWidth, "Database ID", cfg.DatabaseID))
		lines = append(lines, renderKeyValue(labelWidth, "Database", cfg.DatabaseName))
		lines = append(lines, renderKeyValue(labelWidth, "Enabled", fmt.Sprintf("%t", cfg.IsEnabled)))
		lines = append(lines, renderKeyValue(labelWidth, "Updated", cfg.UpdatedAt.Format(time.DateTime)))
	}

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.theme.Border.GetForeground()).
		Padding(1, 2)

	return box.Render(strings.Join(lines, "\n"))
}

func (m Model) renderEditBotConfig() string {
	var lines []string
	lines = append(lines, m.theme.ListHeader.Render("Edit bot config"))
	lines = append(lines, "")

	nameValue := m.editForm.Name
	if nameValue == "" {
		nameValue = "—"
	}

	rows := []struct {
		field EditField
		label string
		value string
	}{
		{EditFieldName, "Name", nameValue},
		{EditFieldEnabled, "Enabled", fmt.Sprintf("%t", m.editForm.IsEnabled)},
		{EditFieldDatabase, "Database ID", m.editForm.DatabaseID},
		{EditFieldSave, "Save", ""},
		{EditFieldCancel, "Cancel", ""},
	}

	const labelWidth = 12

	for _, row := range rows {
		line := renderKeyValue(labelWidth, row.label, row.value)

		cursor := " "
		if row.field == m.editCursor {
			cursor = ">"
		}

		line = cursor + " " + line

		lines = append(lines, m.renderFormRow(row.field == m.editCursor, line))
	}

	if m.editDirty {
		lines = append(lines, "")
		lines = append(lines, m.theme.Failed.Render("unsaved changes"))
	}

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.theme.Border.GetForeground()).
		Padding(1, 2)

	return box.Render(strings.Join(lines, "\n"))
}

func (m Model) renderDatabaseProfileSelect() string {
	var lines []string
	lines = append(lines, m.theme.ListHeader.Render("Select database profile"))
	lines = append(lines, "")

	if len(m.databaseProfiles) == 0 {
		lines = append(lines, m.theme.Muted.Render("loading or unavailable"))
	} else {
		const labelWidth = 12

		for i, profile := range m.databaseProfiles {
			name := profile.Name
			if name == "" {
				name = "—"
			}

			value := renderKeyValue(labelWidth, name, profile.Driver)

			cursor := " "
			if i == m.databaseCursor {
				cursor = ">"
			}

			line := cursor + " " + value
			lines = append(lines, m.renderFormRow(i == m.databaseCursor, line))
		}
	}

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.theme.Border.GetForeground()).
		Padding(1, 2)

	return box.Render(strings.Join(lines, "\n"))
}

func renderKeyValue(labelWidth int, label, value string) string {
	if value == "" {
		return label
	}

	return fmt.Sprintf("%-*s: %s", labelWidth, label, value)
}

func (m Model) renderFormRow(selected bool, line string) string {
	if selected {
		return m.theme.FormSelected.Render(line)
	}
	return m.theme.FormItem.Render(line)
}
