package flow

import (
	"context"
	"errors"
	"strings"
)

// ErrUnknownAction is returned when flow cannot resolve an action.
var ErrUnknownAction error = errors.New("unknown flow action")

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
		return buildInlineCatalogStart(), nil
	case StartScenarioReplyWelcome:
		fallthrough
	default:
		return buildReplyWelcomeStart(), nil
	}
}

// HandleAction resolves the next bot view for an action.
//
// nav:back always returns the original start scenario for the current bot.
func (s *Service) HandleAction(ctx context.Context, req ActionRequest) (ViewModel, error) {
	switch req.ActionID {
	case ActionBack:
		return s.Start(ctx, StartRequest{
			BotID:         req.BotID,
			BotName:       req.BotName,
			StartScenario: req.StartScenario,
		})

	case ActionCatalogStart:
		return buildInlineCatalogStart(), nil

	case ActionCabinetOpen:
		return buildDetailView(
			"Мой кабинет",
			"Здесь будут профиль, история, настройки и персональные данные пользователя.",
		), nil

	case ActionSupportOpen:
		return buildDetailView(
			"Поддержка",
			"Здесь будет связь с оператором, FAQ и обработка обращений.",
		), nil

	case ActionReviewsOpen:
		return buildDetailView(
			"Отзывы",
			"Здесь будут отзывы клиентов, рейтинг и публикация новых отзывов.",
		), nil

	case ActionBalanceOpen:
		return buildDetailView(
			"Баланс",
			"Здесь будет баланс аккаунта, пополнение и история операций.",
		), nil

	case ActionBotsMine:
		return buildDetailView(
			"Мои боты",
			"Здесь будет список пользовательских ботов и быстрые действия по ним.",
		), nil

	case ActionOrderLast:
		return buildDetailView(
			"Последний заказ",
			"Здесь будет карточка последнего заказа и повторное оформление.",
		), nil

	case ActionCategoryPhones:
		return buildCategoryView("Телефоны"), nil

	case ActionCategoryLaptops:
		return buildCategoryView("Ноутбуки"), nil

	case ActionCategoryRouters:
		return buildCategoryView("Роутеры"), nil

	case ActionCategoryAudio:
		return buildCategoryView("Аудио"), nil

	default:
		return ViewModel{}, ErrUnknownAction
	}
}

// ResolveReplyAction maps reply-button text to action identifiers.
//
// This is intentionally narrow and only handles known start-menu buttons.
// More advanced reply routing should later becomee session-aware.
func (s *Service) ResolveReplyAction(text string) (ActionID, bool) {
	switch strings.TrimSpace(text) {
	case "Каталог":
		return ActionCatalogStart, true
	case "Мой кабинет":
		return ActionCabinetOpen, true
	case "Поддержка":
		return ActionSupportOpen, true
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

func buildInlineCatalogStart() ViewModel {
	return ViewModel{
		Text: "Каталог",
		Inline: &InlineKeyboardView{
			Sections: []ActionSection{
				{
					Columns: 2,
					Actions: []ActionButton{
						{ID: ActionCategoryPhones, Label: "Телефоны"},
						{ID: ActionCategoryLaptops, Label: "Ноутбуки"},
						{ID: ActionCategoryRouters, Label: "Роутеры"},
						{ID: ActionCategoryAudio, Label: "Аудио"},
					},
				},
				{
					Columns: 1,
					Actions: []ActionButton{
						{ID: ActionBalanceOpen, Label: "Баланс"},
						{ID: ActionBotsMine, Label: "Мои боты"},
						{ID: ActionOrderLast, Label: "Последний заказ"},
					},
				},
			},
		},
		RemoveReply: true,
	}
}

func buildCategoryView(title string) ViewModel {
	return buildDetailView(
		title,
		"Здесь будет каталог категорий, пагинация, карточка товаров и переход к покупке.",
	)
}

func buildDetailView(title, body string) ViewModel {
	return ViewModel{
		Text: title + "\n\n" + body,
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
