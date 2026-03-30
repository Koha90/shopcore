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
type Service struct {
	store Store
}

// NewService construct a flow service.
func NewService(store Store) *Service {
	if store == nil {
		store = NewMemoryStore()
	}
	return &Service{store: store}
}

// Start resolves the initial bot view for /start.
func (s *Service) Start(_ context.Context, req StartRequest) (ViewModel, error) {
	screen := startScreenForScenario(req.StartScenario)

	s.store.Put(req.SessionKey, Session{
		Current: screen,
		History: nil,
	})

	return s.renderScreen(screen), nil
}

// HandleAction resolve the next bot view for an action.
//
// ActionBack navigates to the previous screen stored in session history.
func (s *Service) HandleAction(ctx context.Context, req ActionRequest) (ViewModel, error) {
	session, ok := s.store.Get(req.SessionKey)
	if !ok {
		session = Session{
			Current: startScreenForScenario(req.StartScenario),
			History: nil,
		}
	}

	switch req.ActionID {
	case ActionBack:
		if len(session.History) == 0 {
			return s.renderScreen(session.Current), nil
		}

		prev := session.History[len(session.History)-1]
		session.History = session.History[:len(session.History)-1]
		session.Current = prev
		s.store.Put(req.SessionKey, session)

		return s.renderScreen(prev), nil
	}

	next, err := resolveNextScreen(req.ActionID)
	if err != nil {
		return ViewModel{}, err
	}

	if next != session.Current {
		session.History = append(session.History, session.Current)
		session.Current = next
		s.store.Put(req.SessionKey, session)
	}

	return s.renderScreen(next), nil
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

func buildEntityView(title string) ViewModel {
	return ViewModel{
		Text: title + "\n\nЗдесь будет следующий шаг сценария для выбранной сущности.",
		Inline: &InlineKeyboardView{
			Sections: []ActionSection{
				{
					Columns: 1,
					Actions: []ActionButton{
						{ID: ActionBack, Label: "Назад"},
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

func resolveNextScreen(actionID ActionID) (ScreenID, error) {
	switch actionID {
	case ActionCatalogStart:
		return ScreenRootCompact, nil

	case ActionRootCompact:
		return ScreenRootCompact, nil

	case ActionRootExtended:
		return ScreenRootExtended, nil

	case ActionCabinetOpen:
		return ScreenCabinet, nil
	case ActionSupportOpen:
		return ScreenSupport, nil
	case ActionReviewsOpen:
		return ScreenReviews, nil

	case ActionBalanceOpen:
		return ScreenBalance, nil
	case ActionBotsMine:
		return ScreenBotsMine, nil
	case ActionOrderLast:
		return ScreenOrderLast, nil

	case ActionEntity1:
		return ScreenEntity1, nil
	case ActionEntity2:
		return ScreenEntity2, nil
	case ActionEntity3:
		return ScreenEntity3, nil
	case ActionEntity4:
		return ScreenEntity4, nil

	default:
		return "", ErrUnknownAction
	}
}

func (s *Service) renderScreen(screen ScreenID) ViewModel {
	switch screen {
	case ScreenReplyWelcome:
		return buildReplyWelcomeStart()

	case ScreenRootCompact:
		return buildCompactRootSelectionView()

	case ScreenRootExtended:
		return buildExtendedRootSelectionView()

	case ScreenEntity1:
		return buildEntityView("Москва")
	case ScreenEntity2:
		return buildEntityView("СПб")
	case ScreenEntity3:
		return buildEntityView("Казань")
	case ScreenEntity4:
		return buildEntityView("Екатеринбург")

	case ScreenCabinet:
		return buildReplyDetailView(
			"Мой кабинет",
			"Здесь будут профиль, история, настройки и персональные данные пользователя.",
		)

	case ScreenSupport:
		return buildReplyDetailView(
			"Поддержка",
			"Здесь будет связь с оператором, FAQ и обработка обращений.",
		)

	case ScreenReviews:
		return buildReplyDetailView(
			"Отзывы",
			"Здесь будут отзывы клиентов, рейтинг и публикация новых отзывов.",
		)

	case ScreenBalance:
		return buildDetailView(
			"Баланс",
			"Здесь будет баланс аккаунта, пополнение и история операций.",
			ActionBack,
		)

	case ScreenBotsMine:
		return buildDetailView(
			"Мои боты",
			"Здесь будет список пользовательских ботов и быстрые действия по ним.",
			ActionBack,
		)

	case ScreenOrderLast:
		return buildDetailView(
			"Последний заказ",
			"Здесь будет карточка последнего заказа и повторное оформление.",
			ActionBack,
		)

	default:
		return buildReplyWelcomeStart()

	}
}
