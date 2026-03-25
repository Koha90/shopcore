package tui

import (
	"context"

	"charm.land/bubbles/v2/textinput"
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
	StatusFilterDisabled StatusFilter = "disabled"
)

const StatusDisabled = "disabled"

type ScreenMode string

const (
	ScreenList                   ScreenMode = "list"
	ScreenBotActions             ScreenMode = "bot_actions"
	ScreenBotConfig              ScreenMode = "bot_config"
	ScreenEditBotConfig          ScreenMode = "edit_bot_config"
	ScreenEditBotToken           ScreenMode = "edit_bot_token"
	ScreenSelecteDatabaseProfile ScreenMode = "select_database_profile"
	ScreenConfirmDiscardEdit     ScreenMode = "confirm_discard_edit"
)

type InputMode int

const (
	InputModeNone InputMode = iota
	InputModeEditName
	InputModeEditToken
)

type BotRow struct {
	ID           string
	Name         string
	DatabaseID   string
	DatabaseName string
	IsEnabled    bool
	TokenMasked  string

	Status    string
	LastError string
}

// BotConfigService defines configuration operations required by TUI.
type BotConfigService interface {
	BotByID(ctx context.Context, id string) (botconfig.BotView, error)
	ListBots(ctx context.Context) ([]botconfig.BotView, error)
	ListDatabaseProfiles(ctx context.Context) ([]botconfig.DatabaseProfileView, error)
	UpdateBot(ctx context.Context, params botconfig.UpdateBotParams) error

	BotToken(ctx context.Context, id string) (string, error)
	UpdateBotToken(ctx context.Context, id string, token string) error
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
	Disabled int
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

	textInput textinput.Model
	inputMode InputMode

	editForm   BotConfigEditForm
	editCursor EditField
	editDirty  bool

	confirmCursor int

	bots         []BotRow
	filteredBots []BotRow
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
func NewModel(mgr BotManager, cfg BotConfigService, theme Theme) Model {
	model := Model{
		manager:      mgr,
		config:       cfg,
		theme:        theme,
		pageSize:     defaultPageSize,
		layout:       LayoutDesktop,
		statusFilter: StatusFilterAll,
		screen:       ScreenList,
		textInput:    newTexInput(),
		inputMode:    InputModeNone,
	}

	model.refresh()
	return model
}

// Init initializes TUI model.
func (m Model) Init() tea.Cmd {
	return tickCmd()
}

func newTexInput() textinput.Model {
	ti := textinput.New()
	ti.Placeholder = ""
	ti.CharLimit = 512
	ti.SetWidth(48)
	return ti
}
