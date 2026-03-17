package tui

import (
	"context"

	tea "charm.land/bubbletea/v2"

	"botmanager/internal/botconfig"
	"botmanager/internal/manager"
)

const defaultPageSize = 10

type LayoutMode string

const (
	LayoutDesktop LayoutMode = "desktop"
	LayoutMobile  LayoutMode = "mobile"
)

type StatusFilter string

const (
	StatusFilterAll      StatusFilter = "all"
	StatusFilterRunning  StatusFilter = "running"
	StatusFilterStopped  StatusFilter = "stopped"
	StatusFilterFailed   StatusFilter = "failed"
	StatusFilterStarting StatusFilter = "starting"
	StatusFilterStopping StatusFilter = "stopping"
)

type ScreenMode string

const (
	ScreenList       ScreenMode = "list"
	ScreenBotActions ScreenMode = "bot_actions"
	ScreenBotConfig  ScreenMode = "bot_config"
)

// BotConfigReader defines configuration queries required by TUI.
type BotConfigReader interface {
	BotByID(ctx context.Context, id string) (botconfig.BotView, error)
}

// BotManager defines mangager operations required by TUI.
type BotManager interface {
	List() []manager.Info
	Start(ctx context.Context, id string) error
	Stop(id string) error
	Restart(ctx context.Context, id string) error
	Info(id string) (manager.Info, error)
}

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
	config  BotConfigReader
	theme   Theme

	layout LayoutMode
	screen ScreenMode

	statusFilter StatusFilter

	selectedBotConfig *botconfig.BotView

	bots         []manager.Info
	filteredBots []manager.Info
	summary      Summary

	cursor       int
	actionCursor int

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
func NewModel(m BotManager, cfg BotConfigReader, theme Theme) Model {
	model := Model{
		manager:      m,
		config:       cfg,
		theme:        theme,
		pageSize:     defaultPageSize,
		layout:       LayoutDesktop,
		statusFilter: StatusFilterAll,
		screen:       ScreenList,
	}

	model.refresh()
	return model
}

// Init initializes TUI model.
func (m Model) Init() tea.Cmd {
	return tickCmd()
}

// selectedConfig loads selected bot config on demand.
// TODO: remove I/O from view path and reuse loaded config state/cmd flow.
func (m Model) selectedConfig() *botconfig.BotView {
	id := m.selectedID()
	if id == "" || m.config == nil {
		return nil
	}

	cfg, err := m.config.BotByID(context.Background(), id)
	if err != nil {
		return nil
	}

	return &cfg
}
