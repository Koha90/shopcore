package flow

import (
	"context"
	"errors"
	"strings"
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
	store                 Store
	provider              CatalogProvider
	categories            CategoryCreator
	cities                CityCreator
	cityLister            CityLister
	districts             DistrictCreator
	categoryLister        CategoryLister
	products              ProductCreator
	productLister         ProductLister
	variants              VariantCreator
	districtLister        DistrictLister
	variantLister         VariantLister
	districtVariants      DistrictVariantCreator
	districtVariantPrices DistrictVariantPriceUpdater
	districtPlacements    DistrictPlacementReader
}

// ServiceDeps contains optional flow dependencies used by admin/catalog actions.
type ServiceDeps struct {
	Categories            CategoryCreator
	Cities                CityCreator
	CityLister            CityLister
	Districts             DistrictCreator
	CategoryLister        CategoryLister
	Products              ProductCreator
	ProductLister         ProductLister
	Variants              VariantCreator
	DistrictLister        DistrictLister
	VariantLister         VariantLister
	DistrictVariants      DistrictVariantCreator
	DistrictVariantPrices DistrictVariantPriceUpdater
	DistrictPlacements    DistrictPlacementReader
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
func NewServiceWithDeps(store Store, provider CatalogProvider, deps ServiceDeps) *Service {
	if store == nil {
		store = NewMemoryStore()
	}
	if provider == nil {
		provider = NewStaticCatalogProvider(DemoCatalog())
	}

	return &Service{
		store:                 store,
		provider:              provider,
		categories:            deps.Categories,
		cities:                deps.Cities,
		cityLister:            deps.CityLister,
		districts:             deps.Districts,
		categoryLister:        deps.CategoryLister,
		products:              deps.Products,
		productLister:         deps.ProductLister,
		variants:              deps.Variants,
		districtLister:        deps.DistrictLister,
		variantLister:         deps.VariantLister,
		districtVariants:      deps.DistrictVariants,
		districtVariantPrices: deps.DistrictVariantPrices,
		districtPlacements:    deps.DistrictPlacements,
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

	return s.renderScreen(catalog, screen, req.CanAdmin), nil
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

func buildCompactRootSelectionView(roots []CatalogNode) ViewModel {
	return buildRootSelectionView(DefaultCompactRootColumns, RootVariantCompact, roots, false)
}

func buildExtendedRootSelectionView(roots []CatalogNode, canAdmin bool) ViewModel {
	return buildRootSelectionView(DefaultExtendedRootColumns, RootVariantExtended, roots, canAdmin)
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
