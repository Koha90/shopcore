package tui

import (
	"context"

	tea "charm.land/bubbletea/v2"

	"github.com/koha90/shopcore/internal/botconfig"
	"github.com/koha90/shopcore/internal/manager"
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
	ScreenList                   ScreenMode = "list"
	ScreenBotActions             ScreenMode = "bot_actions"
	ScreenBotConfig              ScreenMode = "bot_config"
	ScreenEditBotConfig          ScreenMode = "edit_bot_config"
	ScreenSelecteDatabaseProfile ScreenMode = "select_database_profile"
	ScreenConfirmDiscardEdit     ScreenMode = "confirm_discard_edit"
)

// BotConfigService defines configuration operations required by TUI.
type BotConfigService interface {
	BotByID(ctx context.Context, id string) (botconfig.BotView, error)
	ListDatabaseProfiles(ctx context.Context) ([]botconfig.DatabaseProfileView, error)
	UpdateBot(ctx context.Context, params botconfig.UpdateBotParams) error
}

// BotManager defines mangager operations required by TUI.
type BotManager interface {
	List() []manager.Info
	Start(ctx context.Context, id string) error
	Stop(id string) error
	Restart(ctx context.Context, id string) error
	Info(id string) (manager.Info, error)
	Rename(id string, name string) error
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

// BotConfigEditForm represents editable bot configuration fields in TUI.
type BotConfigEditForm struct {
	Name       string
	IsEnabled  bool
	DatabaseID string
}

// EditField identifies currently selected field in bot config edit screen.
type EditField int

const (
	EditFieldName EditField = iota
	EditFieldEnabled
	EditFieldDatabase
	EditFieldSave
	EditFieldCancel
)

// Model represents Bubble Tea application model.
type Model struct {
	manager BotManager
	config  BotConfigService
	theme   Theme

	layout LayoutMode
	screen ScreenMode

	statusFilter StatusFilter

	selectedBotConfig        *botconfig.BotView
	selectedBotConfigID      string
	selectedBotConfigLoading bool

	editForm   BotConfigEditForm
	editCursor EditField
	editDirty  bool
	editTyping bool
	editBuffer string

	confirmCursor int

	bots         []manager.Info
	filteredBots []manager.Info
	summary      Summary

	databaseProfiles []botconfig.DatabaseProfileView
	databaseCursor   int

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
func NewModel(m BotManager, cfg BotConfigService, theme Theme) Model {
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
