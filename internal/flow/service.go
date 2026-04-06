package flow

import (
	"context"
	"errors"
	"strings"

	catalogservice "github.com/koha90/shopcore/internal/catalog/service"
)

var (
	// ErrUnknownAction is returned when flow cannot resolve an action.
	ErrUnknownAction = errors.New("unknown flow action")

	// ErrUnknownPendingInput is returned when flow cannot resolve active pending input state.
	ErrUnknownPendingInput = errors.New("unknown pending input")
)

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
// It resolves navigation screens, session-aware text input and admin flow actions.
type Service struct {
	store      Store
	provider   CatalogProvider
	categories CategoryCreator
}

// NewService constructs transport-agnostic flow service.
//
// If store is nil, in-memory session storage is used.
// Demo catalog is wired by default until persistent catalog source appears.
// Admin category creator is disabled in this constructor.
func NewService(store Store) *Service {
	return NewServiceWithCatalogProvider(store, NewStaticCatalogProvider(DemoCatalog()))
}

// NewServiceWithCatalogProvider constructs flow service with explicit catalog provider.
//
// This constructor is intended for tests and runtime wiring that only need
// catalog navigation. Admin category creator remains disabled.
func NewServiceWithCatalogProvider(store Store, provider CatalogProvider) *Service {
	if store == nil {
		store = NewMemoryStore()
	}
	if provider == nil {
		provider = NewStaticCatalogProvider(DemoCatalog())
	}

	return &Service{
		store:    store,
		provider: provider,
	}
}

// NewServiceWithDeps constructs flow service with explicit dependencies.
//
// It allows wiring a custom catalog provider and optional admin category creator.
// This constructor is intended for application wiring and tests.
func NewServiceWithDeps(store Store, provider CatalogProvider, categories CategoryCreator) *Service {
	if store == nil {
		store = NewMemoryStore()
	}
	if provider == nil {
		provider = NewStaticCatalogProvider(DemoCatalog())
	}

	return &Service{
		store:      store,
		provider:   provider,
		categories: categories,
	}
}

// Start resolves the initial bot view for /start.
//
// StartScenario controls whether the user sees reply welcome
// or enters inline catalog immediately.
func (s *Service) Start(ctx context.Context, req StartRequest) (ViewModel, error) {
	catalog, err := s.provider.Catalog(ctx)
	if err != nil {
		return ViewModel{}, err
	}

	screen := startScreenForScenario(req.StartScenario)

	s.store.Put(req.SessionKey, Session{
		Current:  screen,
		History:  nil,
		Pending:  PendingInput{},
		CanAdmin: req.CanAdmin,
	})

	return s.renderScreen(catalog, screen), nil
}

// HandleAction resolve the next flow view for an action.
//
// Resolution order:
//   - ActionBack uses session history
//   - ActionCatalogStart opens scenario-aware catalog root
//   - admin actions open stable admin screen or start pending text input
//   - generic catalog selection actions advance inside CatalogSchema
//   - explicit non-catalog action open stable detail screen
//
// Any non-pending action transition clears active pending text input state.
func (s *Service) HandleAction(ctx context.Context, req ActionRequest) (ViewModel, error) {
	catalog, err := s.provider.Catalog(ctx)
	if err != nil {
		return ViewModel{}, err
	}

	session, ok := s.store.Get(req.SessionKey)
	if !ok {
		session = Session{
			Current:  startScreenForScenario(req.StartScenario),
			History:  nil,
			Pending:  PendingInput{},
			CanAdmin: req.CanAdmin,
		}
	} else {
		session = s.syncSessionAccess(req.SessionKey, session, req.CanAdmin)
	}

	switch req.ActionID {
	case ActionBack:
		if len(session.History) == 0 {
			if session.Pending.Active() {
				session.Pending = PendingInput{}
				s.store.Put(req.SessionKey, session)
			}
			return s.renderScreen(catalog, session.Current), nil
		}

		prev := session.History[len(session.History)-1]
		session.History = session.History[:len(session.History)-1]
		session.Current = prev
		session.Pending = PendingInput{}
		s.store.Put(req.SessionKey, session)

		return s.renderScreen(catalog, prev), nil

	case ActionCatalogStart:
		next := catalogRootForScenario(req.StartScenario)

		if next != session.Current {
			session.History = append(session.History, session.Current)
			session.Current = next
		}
		session.Pending = PendingInput{}
		s.store.Put(req.SessionKey, session)

		return s.renderScreen(catalog, next), nil

	case ActionAdminCategoryCreateStart:
		next := ScreenAdminCategoryCreate

		if next != session.Current {
			session.History = append(session.History, session.Current)
		}
		session.Current = next
		session.Pending = PendingInput{
			Kind:    PendingInputCategoryName,
			Payload: nil,
		}
		s.store.Put(req.SessionKey, session)

		return s.renderScreen(catalog, next), nil
	}

	if next, err := s.resolveCatalogScreen(catalog, session.Current, req.ActionID); err == nil {
		if next != session.Current {
			session.History = append(session.History, session.Current)
			session.Current = next
		}
		session.Pending = PendingInput{}
		s.store.Put(req.SessionKey, session)

		return s.renderScreen(catalog, next), nil
	}

	next, err := resolveNextScreen(req.ActionID)
	if err != nil {
		return ViewModel{}, err
	}

	if next != session.Current {
		session.History = append(session.History, session.Current)
		session.Current = next
	}
	session.Pending = PendingInput{}
	s.store.Put(req.SessionKey, session)

	return s.renderScreen(catalog, next), nil
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

// HandleText resolves a plain text message relative to current session state.
//
// If no pending input exists, the current screen is rendered again.
// If pending input exists, text is handled as a continuation of that flow step.
//
// Current behavior supports admin category creation with automatic code
// suggestion and manual code fallback.
func (s *Service) HandleText(ctx context.Context, req TextRequest) (ViewModel, error) {
	catalog, err := s.provider.Catalog(ctx)
	if err != nil {
		return ViewModel{}, err
	}

	session, ok := s.store.Get(req.SessionKey)
	if !ok {
		session = Session{
			Current:  startScreenForScenario(req.StartScenario),
			History:  nil,
			Pending:  PendingInput{},
			CanAdmin: req.CanAdmin,
		}
	} else {
		session = s.syncSessionAccess(req.SessionKey, session, req.CanAdmin)
	}

	if !session.Pending.Active() {
		return s.renderScreen(catalog, session.Current), nil
	}

	switch session.Pending.Kind {
	case PendingInputCategoryName:
		name := strings.TrimSpace(req.Text)
		if name == "" {
			return buildAdminCategoryCreateInputView("Название категории не может быть пустым."), nil
		}

		session.Pending.SetValue(PendingValueName, name)

		suggestedCode := catalogservice.SuggestCategoryCode(name)
		if suggestedCode == "" {
			session.Current = ScreenAdminCategoryCode
			session.Pending.Kind = PendingInputCategoryCode
			s.store.Put(req.SessionKey, session)

			return buildAdminCategoryCodeInputView(
				"Не удалось автоматически подобрать code.",
				"",
			), nil
		}

		session.Pending.SetValue(PendingValueCode, suggestedCode)

		if s.categories == nil {
			return ViewModel{}, errors.New("flow category creator is nil")
		}

		err := s.categories.CreateCategory(ctx, CreateCategoryParams{
			Code: suggestedCode,
			Name: name,
		})
		if err != nil {
			session.Current = ScreenAdminCategoryCode
			session.Pending.Kind = PendingInputCategoryCode
			s.store.Put(req.SessionKey, session)

			return buildAdminCategoryCodeInputView(
				"Не удалось создать категорию с автоматическим code.",
				"",
			), nil
		}

		session.Pending = PendingInput{}
		session.Current = ScreenAdminCategoryCreateDone
		s.store.Put(req.SessionKey, session)

		return s.renderScreen(catalog, session.Current), nil

	case PendingInputCategoryCode:
		code := strings.TrimSpace(req.Text)
		if code == "" {
			return buildAdminCategoryCodeInputView("Code категории не может быть пустым.", ""), nil
		}

		name := strings.TrimSpace(session.Pending.Value(PendingValueName))
		if name == "" {
			return ViewModel{}, errors.New("pending category name is empty")
		}

		session.Pending.SetValue(PendingValueCode, code)

		if s.categories == nil {
			return ViewModel{}, errors.New("flow category creator is nil")
		}

		err := s.categories.CreateCategory(ctx, CreateCategoryParams{
			Code: session.Pending.Value(PendingValueCode),
			Name: name,
		})
		if err != nil {
			return buildAdminCategoryCodeInputView("Не удалось создать категорию. Попробуйте другой code.", code), nil
		}

		session.Pending = PendingInput{}
		session.Current = ScreenAdminCategoryCreateDone
		s.store.Put(req.SessionKey, session)

		return s.renderScreen(catalog, session.Current), nil

	default:
		return ViewModel{}, ErrUnknownPendingInput
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

func buildCompactRootSelectionView(roots []CatalogNode) ViewModel {
	return buildRootSelectionView(DefaultCompactRootColumns, RootVariantCompact, roots)
}

func buildExtendedRootSelectionView(roots []CatalogNode) ViewModel {
	return buildRootSelectionView(DefaultExtendedRootColumns, RootVariantExtended, roots)
}

// buildRootSelectionView renders the root inline selection screen.
//
// The compact variant renders only the main selectable entities.
// The extended variant renders the same entities plus utility action below.
func buildRootSelectionView(columns int, variant RootVariant, roots []CatalogNode) ViewModel {
	cols := normalizeColumns(columns)

	actions := make([]ActionButton, 0, len(roots))
	for _, node := range roots {
		actions = append(actions, ActionButton{
			ID:    catalogSelectAction(node.Level, node.ID),
			Label: node.Label,
		})
	}

	sections := []ActionSection{
		{
			Columns: cols,
			Actions: actions,
		},
	}

	if variant == RootVariantExtended {
		sections = append(sections, ActionSection{
			Columns: 1,
			Actions: []ActionButton{
				{ID: ActionBalanceOpen, Label: "Баланс"},
				{ID: ActionBotsMine, Label: "Мои боты"},
				{ID: ActionOrderLast, Label: "Последний заказ"},
				{ID: ActionAdminOpen, Label: "Админка"},
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

// resolveCatalogScreen resolves one generic catalog selection action
// relative to the current screen state.
//
// It validates:
//   - action payload format
//   - current catalog path
//   - expected next schema level
//   - existence of target node in catalog tree
func (s *Service) resolveCatalogScreen(catalog Catalog, current ScreenID, actionID ActionID) (ScreenID, error) {
	level, id, ok := parseCatalogSelectAction(actionID)
	if !ok {
		return "", ErrUnknownAction
	}

	var currentPath CatalogPath

	switch current {
	case ScreenRootCompact, ScreenRootExtended:
		currentPath = nil

	default:
		path, ok := parseCatalogScreen(current)
		if !ok {
			return "", ErrUnknownAction
		}
		currentPath = path
	}

	expectedLevel, ok := s.expectedNextCatalogLevel(catalog, currentPath)
	if !ok {
		return "", ErrUnknownAction
	}
	if level != expectedLevel {
		return "", ErrUnknownAction
	}

	nextPath := currentPath.Append(level, id)

	if _, ok := catalog.FindNode(nextPath); !ok {
		return "", ErrUnknownAction
	}

	return catalogScreen(nextPath), nil
}

// expectedNextCatalogLevel returns which catalog level may be selected next
// for the provided path.
func (s *Service) expectedNextCatalogLevel(catalog Catalog, path CatalogPath) (CatalogLevel, bool) {
	if len(path) == 0 {
		return catalog.RootLevel()
	}

	last, ok := path.Last()
	if !ok {
		return "", false
	}

	return catalog.Schema.Next(last.Level)
}

func resolveNextScreen(actionID ActionID) (ScreenID, error) {
	switch actionID {
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

	case ActionAdminOpen:
		return ScreenAdminRoot, nil

	case ActionAdminCatalogOpen:
		return ScreenAdminCatalog, nil

	default:
		return "", ErrUnknownAction
	}
}

// renderScreen converts logical screen identifiers into transport-agnostic view models.
//
// Stable root/detail screens are handled directly.
// Dynamic catalog drill-down screens are rendered from CatalogPath.
func (s *Service) renderScreen(catalog Catalog, screen ScreenID) ViewModel {
	switch screen {
	case ScreenReplyWelcome:
		return buildReplyWelcomeStart()

	case ScreenRootCompact:
		return buildCompactRootSelectionView(catalog.RootNodes())

	case ScreenRootExtended:
		return buildExtendedRootSelectionView(catalog.RootNodes())

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

	case ScreenAdminRoot:
		return buildAdminRootView()

	case ScreenAdminCatalog:
		return buildAdminCatalogView()

	case ScreenAdminCategoryCreate:
		return buildAdminCategoryCreateInputView("")

	case ScreenAdminCategoryCode:
		return buildAdminCategoryCodeInputView("", "")

	case ScreenAdminCategoryCreateDone:
		return buildAdminCategoryCreateDoneView()
	}

	path, ok := parseCatalogScreen(screen)
	if !ok {
		return buildReplyWelcomeStart()
	}

	node, found := catalog.FindNode(path)
	if !found {
		return buildReplyWelcomeStart()
	}

	if len(node.Children) > 0 {
		return buildCatalogNodeView(node)
	}

	return buildCatalogLeafView(node)
}

func (s *Service) syncSessionAccess(key SessionKey, session Session, canAdmin bool) Session {
	if session.CanAdmin == canAdmin {
		return session
	}

	session.CanAdmin = canAdmin
	s.store.Put(key, session)

	return session
}

func buildAdminRootView() ViewModel {
	return ViewModel{
		Text: "Админка\n\nВыберите раздел:",
		Inline: &InlineKeyboardView{
			Sections: []ActionSection{
				{
					Columns: 1,
					Actions: []ActionButton{
						{ID: ActionAdminCatalogOpen, Label: "Каталог"},
					},
				},
			},
		},
		RemoveReply: true,
	}
}

func buildAdminCatalogView() ViewModel {
	return ViewModel{
		Text: "Админка · Каталог\n\nВыберите действие:",
		Inline: &InlineKeyboardView{
			Sections: []ActionSection{
				{
					Columns: 1,
					Actions: []ActionButton{
						{ID: ActionAdminCategoryCreateStart, Label: "Создать категорию"},
						{ID: ActionBack, Label: "Назад"},
					},
				},
			},
		},
		RemoveReply: true,
	}
}

func buildAdminCategoryCreateInputView(validation string) ViewModel {
	text := "Новая категория\n\nВведите название категории сообщением."
	if validation != "" {
		text = "Новая категория\n\n" + validation + "\n\nВведите название категории сообщением."
	}

	return ViewModel{
		Text: text,
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

func buildAdminCategoryCodeInputView(validation, suggested string) ViewModel {
	text := "Новая категория\n\nВведите code категрии сообщением."
	if suggested != "" {
		text = "Новая категория\n\nАвто-код: " + suggested + "\n\nВведите code категории сообщением."
	}
	if validation != "" {
		text = "Новая категория\n\n" + validation + "\n\nВведите code категории сообщением."
	}

	return ViewModel{
		Text: text,
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

func buildAdminCategoryCreateDoneView() ViewModel {
	return ViewModel{
		Text: "Новая категория\n\nКатегория создана.",
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

// HasPendingInput reports whether the session currently expects plain text input.
func (s *Service) HasPendingInput(key SessionKey) bool {
	if s == nil || s.store == nil {
		return false
	}

	session, ok := s.store.Get(key)
	if !ok {
		return false
	}

	return session.Pending.Active()
}
