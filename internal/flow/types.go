package flow

import "strings"

// StartScenario identifies the initial bot representation shown on /start.
type StartScenario string

const (
	// StartScenarioReplyWelcome renders welcome text with reply keyboard.
	StartScenarioReplyWelcome StartScenario = "reply_welcome"

	// StartScenarioInlineCatalog renders inline catalog immediately.
	StartScenarioInlineCatalog StartScenario = "inline_catalog"
)

// NormalizeStartScenario converts arbitrary input into a supported scenario.
//
// Empty or unknown values fall back to StartScenarioReplyWelcome.
func NormalizeStartScenario(v string) StartScenario {
	switch StartScenario(strings.TrimSpace(v)) {
	case StartScenarioInlineCatalog:
		return StartScenarioInlineCatalog
	case StartScenarioReplyWelcome:
		fallthrough
	default:
		return StartScenarioReplyWelcome
	}
}

// ActionID is a transport-agnostic action identifier.
//
// It is used by both inline callbacks and reply-button routing.
type ActionID string

const (
	ActionCatalogStart ActionID = "catalog:start"
	ActionCabinetOpen  ActionID = "cabinet:open"
	ActionSupportOpen  ActionID = "support:open"
	ActionReviewsOpen  ActionID = "reviews:open"

	ActionBalanceOpen ActionID = "balance:open"
	ActionBotsMine    ActionID = "bots:mine"
	ActionOrderLast   ActionID = "order:last"

	ActionAdminOpen                ActionID = "admin:open"
	ActionAdminCatalogOpen         ActionID = "admin:catalog:open"
	ActionAdminCategoryCreateStart ActionID = "admin:category:create:start"
	ActionAdminCityCreateStart     ActionID = "admin:city:create:start"

	ActionAdminDistrictVariantPriceUpdateStart ActionID = "admin:district_variant:price:update:start"

	ActionRootCompact  ActionID = "root:compact"
	ActionRootExtended ActionID = "root:extended"

	ActionBack ActionID = "nav:back"
)

// ViewModel describes a bot screen independent of concrete transport.
type ViewModel struct {
	Text string

	// Inline contains inline button sections, if any.
	Inline *InlineKeyboardView

	// Reply contains reply keyboard rows, if any.
	Reply *ReplyKeyboardView

	// RemoveReply tells transport to hide reply keyboard.
	RemoveReply bool
}

// InlineKeyboardView describes an inline keyboard grouped into sections.
type InlineKeyboardView struct {
	Sections []ActionSection
}

// ActionSection describes one inline keyboard section.
//
// Columns controls how many buttons should be rendered in a single row for
// this section.
type ActionSection struct {
	Columns int
	Actions []ActionButton
}

// ActionButton describes one transport-agnostic clickable action.
type ActionButton struct {
	ID    ActionID
	Label string
}

// ReplyKeyboardView describes a reply keyboard as rows of buttons.
type ReplyKeyboardView struct {
	Rows [][]ReplyButton
}

// ReplyButton describes one reply-keyboard button.
type ReplyButton struct {
	ID    ActionID
	Label string
}

// StartRequest contains data required to resolve /start presentation.
type StartRequest struct {
	BotID         string
	BotName       string
	StartScenario string
	SessionKey    SessionKey
	CanAdmin      bool
}

// ActionRequest contains data required to resolve the next screen after action.
type ActionRequest struct {
	BotID         string
	BotName       string
	StartScenario string
	ActionID      ActionID
	SessionKey    SessionKey
	CanAdmin      bool
}

// TextRequest contains data required to resolve a text message in current flow.
type TextRequest struct {
	BotID         string
	BotName       string
	StartScenario string
	Text          string
	SessionKey    SessionKey
	CanAdmin      bool
}
