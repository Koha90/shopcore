package flow

import (
	"context"
	"errors"
	"strings"
)

// ErrUnknownAction is returned when flow cannot resolve an action.
var ErrUnknownAction = errors.New("unknown flow action")

const (
	// DefaultCompactRootColumns defines the default column count for inline root
	// view opened from reply-based welcome scenario.
	DefaultCompactRootColumns = 1

	// DefaultExtendedRootColumns defines the default column count for inline root
	// view opened directly by inline start scenario.
	DefaultExtendedRootColumns = 2
)

// RootVariant controls how the root inline selection view sould be rendered.
type RootVariant string

const (
	// RootVariantCompact renders only the main selectable entities.
	RootVariantCompact RootVariant = "compact"

	// CatalogVariantExtended renders the main selectable entities plus utility
	// actions below.
	RootVariantExtended RootVariant = "extended"
)

// Service builds initial and next-step views for bot flows.
//
// The service is transport-agnostic and contains no Telegram-specific code.
type Service struct{}

// NewService construct a flow service.
func NewService() *Service {
	return &Service{}
}

// Start resolves the initial bot view for /start.
func (s *Service) Start(_ context.Context, req StartRequest) (ViewModel, error) {
	switch NormalizeStartScenario(req.StartScenario) {
	case StartScenarioInlineCatalog:
		return buildExtendedRootSelectionView(), nil

	case StartScenarioReplyWelcome:
		fallthrough
	default:
		return buildReplyWelcomeStart(), nil
	}
}

// HandleAction resolves the next bot view for an action.
//
// For nav:back always returns the root inline screen that matches the current
// start scenario. Later it should become session/history-aware.
func (s *Service) HandleAction(ctx context.Context, req ActionRequest) (ViewModel, error) {
	switch req.ActionID {
	case ActionRootCompact:
		return buildCompactRootSelectionView(), nil

	case ActionRootExtended:
		return buildExtendedRootSelectionView(), nil

	case ActionCatalogStart:
		return buildCompactRootSelectionView(), nil

	case ActionCabinetOpen:
		return buildReplyDetailView(
			"Мой кабинет",
			"Здесь будут профиль, история, настройки и персональные данные пользователя.",
		), nil

	case ActionSupportOpen:
		return buildReplyDetailView(
			"Поддержка",
			"Здесь будет связь с оператором, FAQ и обработка обращений.",
		), nil

	case ActionReviewsOpen:
		return buildReplyDetailView(
			"Отзывы",
			"Здесь будут отзывы клиентов, рейтинг и публикация новых отзывов.",
		), nil

	case ActionBalanceOpen:
		return buildDetailView(
			"Баланс",
			"Здесь будет баланс аккаунта, пополнение и история операций.",
			ActionRootExtended,
		), nil

	case ActionBotsMine:
		return buildDetailView(
			"Мои боты",
			"Здесь будет список пользовательских ботов и быстрые действия по ним.",
			ActionRootExtended,
		), nil

	case ActionOrderLast:
		return buildDetailView(
			"Последний заказ",
			"Здесь будет карточка последнего заказа и повторное оформление.",
			ActionRootExtended,
		), nil

	case ActionEntity1:
		return buildEntityView("Москва", backActionForScenario(req.StartScenario)), nil

	case ActionEntity2:
		return buildEntityView("СПб", backActionForScenario(req.StartScenario)), nil

	case ActionEntity3:
		return buildEntityView("Казань", backActionForScenario(req.StartScenario)), nil

	case ActionEntity4:
		return buildEntityView("Екатеринбург", backActionForScenario(req.StartScenario)), nil

	default:
		return ViewModel{}, ErrUnknownAction
	}
}

// ResolveReplyAction maps reply-button text to action identifiers.
//
// For now this resolver is intentionally narrow and stateless.
// Later reply routing should become session-aware.
func (s *Service) ResolveReplyAction(text string) (ActionID, bool) {
	switch strings.TrimSpace(text) {
	case "Каталог", "♻️ Каталог":
		return ActionCatalogStart, true

	case "Мой кабинет", "⚙️ Мой кабинет":
		return ActionCabinetOpen, true

	case "Поддержка", "🤷‍♂️ Поддержка":
		return ActionSupportOpen, true

	case "Отзывы", "📨 Отзывы":
		return ActionReviewsOpen, true

	default:
		return "", false
	}
}

func buildReplyWelcomeStart() ViewModel {
	return ViewModel{
		Text: "Добро пожаловать 👋\nВыберите раздел:",
		Reply: &ReplyKeyboardView{
			Rows: [][]ReplyButton{
				{
					{ID: ActionCatalogStart, Label: "♻️ Каталог"},
					{ID: ActionCabinetOpen, Label: "⚙️ Мой кабинет"},
				},
				{
					{ID: ActionSupportOpen, Label: "🤷‍♂️ Поддержка"},
					{ID: ActionReviewsOpen, Label: "📨 Отзывы"},
				},
			},
		},
	}
}

func buildCompactRootSelectionView() ViewModel {
	return buildRootSelectionView(DefaultCompactRootColumns, RootVariantCompact)
}

func buildExtendedRootSelectionView() ViewModel {
	return buildRootSelectionView(DefaultExtendedRootColumns, RootVariantExtended)
}

// buildRootSelectionView renders the root inline selection screen.
//
// The compact variant renders only the main selectable entities.
// The extended variant renders the same entities plus utility action below.
func buildRootSelectionView(columns int, variant RootVariant) ViewModel {
	cols := normalizeColumns(columns)

	sections := []ActionSection{
		{
			Columns: cols,
			Actions: []ActionButton{
				{ID: ActionEntity1, Label: "Москва"},
				{ID: ActionEntity2, Label: "СПб"},
				{ID: ActionEntity3, Label: "Казань"},
				{ID: ActionEntity4, Label: "Екатеринбург"},
			},
		},
	}

	if variant == RootVariantExtended {
		sections = append(sections, ActionSection{
			Columns: 1,
			Actions: []ActionButton{
				{ID: ActionBalanceOpen, Label: "Баланс"},
				{ID: ActionBotsMine, Label: "Мои боты"},
				{ID: ActionOrderLast, Label: "Последний заказ"},
			},
		})
	}

	return ViewModel{
		Text: "Каталог\n\nВыберите раздел:",
		Inline: &InlineKeyboardView{
			Sections: sections,
		},
		RemoveReply: true,
	}
}

func buildEntityView(title string, backAction ActionID) ViewModel {
	return ViewModel{
		Text: title + "\n\nЗдесь будет следующий шаг сценартя для выбранной сущности.",
		Inline: &InlineKeyboardView{
			Sections: []ActionSection{
				{
					Columns: 1,
					Actions: []ActionButton{
						{ID: backAction, Label: "Назад"},
					},
				},
			},
		},
		RemoveReply: true,
	}
}

func buildDetailView(title, body string, backAction ActionID) ViewModel {
	return ViewModel{
		Text: title + "\n\n" + body,
		Inline: &InlineKeyboardView{
			Sections: []ActionSection{
				{
					Columns: 1,
					Actions: []ActionButton{
						{ID: backAction, Label: "Назад"},
					},
				},
			},
		},
		RemoveReply: true,
	}
}

// buildReplyDetailView renders a detail screen opened from reply-menu actions.
//
// It intentionally does not include inline back navigation, because those
// screens are not part of inline catalog flow.
func buildReplyDetailView(title, body string) ViewModel {
	return ViewModel{
		Text:        title + "\n\n" + body,
		RemoveReply: false,
	}
}

func backActionForScenario(startScenario string) ActionID {
	switch NormalizeStartScenario(startScenario) {
	case StartScenarioInlineCatalog:
		return ActionRootExtended
	case StartScenarioReplyWelcome:
		fallthrough
	default:
		return ActionRootCompact
	}
}

func normalizeColumns(v int) int {
	if v <= 0 {
		return 1
	}
	return v
}
