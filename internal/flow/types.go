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

	ActionAdminProductImageUpdateStart ActionID = "admin:catalog:product_image:start"
	ActionAdminVariantImageUpdateStart ActionID = "admin:catalog:variant_image:start"

	ActionRootCompact  ActionID = "root:compact"
	ActionRootExtended ActionID = "root:extended"

	// ActionOrderStart starts order flow from selected catalog variant leaf.
	ActionOrderStart ActionID = "order:start"
	// ActionOrderConfirm confirms current order request in flow-only MVP.
	ActionOrderConfirm ActionID = "order:confirm"

	ActionBack ActionID = "nav:back"
)

// ViewModel describes a bot screen independent of concrete transport.
type ViewModel struct {
	Text string

	// Media contains one optional transport-agnostic media attachment.
	Media *MediaView

	// Inline contains inline button sections, if any.
	Inline *InlineKeyboardView

	// Reply contains reply keyboard rows, if any.
	Reply *ReplyKeyboardView

	// RemoveReply tells transport to hide reply keyboard.
	RemoveReply bool

	// Effects contains transport-agnostic side effects requested by flow.
	//
	// Transport may execute supported effects after resolving a view.
	// The field intentionally contains no Telegram-specific types.
	Effects []Effect
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

// PhotoRequest contains data required to resolve a photo message in current flow.
type PhotoRequest struct {
	BotID         string
	BotName       string
	StartScenario string

	// FileToken is an opaque transport-owned photo identifier.
	FileToken string

	// Caption stores optional user-provided photo caption.
	Caption string

	SessionKey SessionKey
	CanAdmin   bool
}

// EffectKind identifies a transport-agnostic side effect by flow.
type EffectKind string

// EffectMediaKind identifies a media kind attached to a flow effect.
type EffectMediaKind string

const (
	// EffectSendText asks transport to send a plain text message to target.
	EffectSendText EffectKind = "send_text"

	// EffectSendPhoto asks transport to send one photo to target.
	EffectSendPhoto EffectKind = "send_photo"
)

const (
	// EffectMediaPhoto identifies one photo attachment.
	EffectMediaPhoto EffectMediaKind = "photo"
)

// EffectTarget describes a recipient for a flow effect.
type EffectTarget struct {
	ChatID int64
	UserID int64
}

// Effect describes one transport-agnostic side effect.
type Effect struct {
	Kind   EffectKind
	Target EffectTarget
	Text   string
	Media  *EffectMedia
}

// EffectMedia describes one opaque media attachment requested by flow.
//
// FileToken is transport-owned. Flow stores and forwards it without trying to
// uderstand whether it is a Telegram file_id, local path or future storage key.
type EffectMedia struct {
	Kind      EffectMediaKind
	FileToken string
}
