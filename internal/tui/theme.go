package tui

import "charm.land/lipgloss/v2"

// Theme defines all visual styles used by TUI.
type Theme struct {
	App          lipgloss.Style
	Title        lipgloss.Style
	Border       lipgloss.Style
	ListHeader   lipgloss.Style
	ListItem     lipgloss.Style
	ListSelected lipgloss.Style
	StatusBar    lipgloss.Style
	Help         lipgloss.Style
	Error        lipgloss.Style
	Success      lipgloss.Style
	Muted        lipgloss.Style
	Running      lipgloss.Style
	Stopped      lipgloss.Style
	Failed       lipgloss.Style
	Starting     lipgloss.Style
	Stopping     lipgloss.Style

	FormItem     lipgloss.Style
	FormSelected lipgloss.Style
}

// GruvboxTheme returns default gruvbox-inspired theme.
func GruvboxTheme() Theme {
	bg := lipgloss.Color("#282828")
	fg := lipgloss.Color("#ebdbb2")
	muted := lipgloss.Color("#a89984")
	border := lipgloss.Color("#504945")

	yellow := lipgloss.Color("#d79921")
	green := lipgloss.Color("#98971a")
	red := lipgloss.Color("#cc241d")
	blue := lipgloss.Color("#458588")
	aqua := lipgloss.Color("#689d6a")
	orange := lipgloss.Color("#d65d0e")

	rowBase := lipgloss.NewStyle()

	return Theme{
		App: lipgloss.NewStyle().
			Background(bg).
			Foreground(fg),

		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(yellow),

		Border: lipgloss.NewStyle().
			Foreground(border),

		ListHeader: lipgloss.NewStyle().
			Bold(true).
			Foreground(aqua),

		ListItem: rowBase.
			Foreground(aqua),

		ListSelected: rowBase.
			Bold(true).
			Foreground(bg).
			Background(blue),

		FormItem: rowBase.
			Foreground(aqua),

		FormSelected: rowBase.
			Bold(true).
			Foreground(bg).
			Background(blue),

		StatusBar: lipgloss.NewStyle().
			Foreground(bg).
			Background(yellow).
			Padding(0, 1),

		Help: lipgloss.NewStyle().
			Foreground(muted),

		Error: lipgloss.NewStyle().
			Foreground(red).
			Bold(true),

		Success: lipgloss.NewStyle().
			Foreground(green).
			Bold(true),

		Muted: lipgloss.NewStyle().
			Foreground(muted),

		Running: lipgloss.NewStyle().
			Foreground(green).
			Bold(true),

		Stopped: lipgloss.NewStyle().
			Foreground(muted),

		Failed: lipgloss.NewStyle().
			Foreground(red).
			Bold(true),

		Starting: lipgloss.NewStyle().
			Foreground(yellow).
			Bold(true),

		Stopping: lipgloss.NewStyle().
			Foreground(orange).
			Bold(true),
	}
}
