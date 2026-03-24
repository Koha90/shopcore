package tui

import (
	"strings"

	"github.com/koha90/shopcore/internal/manager"
)

func (m *Model) refresh() {
	m.bots = m.manager.List()
	m.summary = buildSummary(m.bots)
	m.applyFilter()
	m.clampCursor()
	m.ensureCursorVisible()
}

func (m *Model) applyFilter() {
	q := strings.ToLower(strings.TrimSpace(m.filter))
	result := make([]manager.Info, 0, len(m.bots))

	for _, bot := range m.bots {
		if q != "" {
			if !strings.Contains(strings.ToLower(bot.ID), q) &&
				!strings.Contains(strings.ToLower(bot.Name), q) &&
				!strings.Contains(strings.ToLower(string(bot.Status)), q) {
				continue
			}
		}

		if m.statusFilter != StatusFilterAll &&
			string(bot.Status) != string(m.statusFilter) {
			continue
		}

		result = append(result, bot)
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

func detectLayout(width int) LayoutMode {
	if width < 100 {
		return LayoutMobile
	}
	return LayoutDesktop
}

func (m Model) listTop() int {
	top := 0

	// title
	top++

	if m.layout == LayoutMobile {
		// mobile summary: 3 string in box.
		top += 5
	} else {
		// desktop summary box
		top += 3
	}

	// filter bar
	top++

	// empty line before content
	top++

	// list box:
	// border top + padding top + header
	top += 3

	return top
}
