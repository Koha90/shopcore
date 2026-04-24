package flow

import (
	"context"
	"errors"
	"fmt"
)

// ErrOrderContextUnavailable is returned when current session state does not
// contains resolved order context.
var ErrOrderContextUnavailable = errors.New("order context unavailable")

// OrderContext contains catalog data required to render confirmation.
//
// This is a flow-local snapshot derived from current catalog selection.
// It is intentionally storage-agnostic and does not represent persisted order data.
type OrderContext struct {
	CityID       string
	CityName     string
	DistrictID   string
	DistrictName string
	ProductID    string
	ProductLabel string
	VariantID    string
	VariantLabel string

	// BasePriceText keeps current display-formatted base price from catalog node.
	//
	// The first flow-only order step uses UI-ready price text becouse current
	// catalog navigation already expose PriceText on leaf nodes.
	// Numeric base price should be added as a sepparate small step before
	// payment/promo quote logic is introduced.
	BasePriceText string
}

// CurrentOrderContext returns current order context for provided session key.
//
// It is intended for outer application/runtime layer that need to react to
// successful order confirmation without coupling transport logic to internal
// flow state representation.
func (s *Service) CurrentOrderContext(
	ctx context.Context,
	key SessionKey,
) (OrderContext, error) {
	if s == nil {
		return OrderContext{}, ErrOrderContextUnavailable
	}

	session, ok := s.store.Get(key)
	if !ok {
		return OrderContext{}, ErrOrderContextUnavailable
	}

	catalog, err := s.provider.Catalog(ctx)
	if err != nil {
		return OrderContext{}, err
	}

	orderCtx, ok := orderContext(catalog, session)
	if !ok {
		return OrderContext{}, ErrOrderContextUnavailable
	}

	return orderCtx, nil
}

// handleOrderAction resolves order-specific flow actions.
//
// It is intentionally isolated from generic catalog navigation and from admin
// action handlers. That keeps order flow as a small explicit branch on top of
// existing flow/session/history mechanics.
func (s *Service) handleOrderAction(
	catalog Catalog,
	session Session,
	req ActionRequest,
) (ViewModel, Session, bool, error) {
	switch req.ActionID {
	case ActionOrderStart:
		ctx, ok := catalogLeafOrderContext(catalog, session.Current)
		if !ok {
			return ViewModel{}, session, true, ErrUnknownAction
		}

		next := session
		if next.Current != ScreenOrderConfirm {
			next.History = append(next.History, next.Current)
			next.Current = ScreenOrderConfirm
		}
		next.Pending = PendingInput{}

		return buildOrderConfirmView(ctx), next, true, nil

	case ActionOrderConfirm:
		if session.Current != ScreenOrderConfirm {
			return ViewModel{}, session, true, ErrUnknownAction
		}

		if _, ok := orderContext(catalog, session); !ok {
			return ViewModel{}, session, true, ErrUnknownAction
		}

		next := session
		if next.Current != ScreenOrderDone {
			next.History = append(next.History, next.Current)
			next.Current = ScreenOrderDone
		}
		next.Pending = PendingInput{}

		return buildOrderDoneView(), next, true, nil

	default:
		return ViewModel{}, session, false, nil
	}
}

// renderOrderScreen renders order-specific screens.
//
// Order flow stays small on purpose:
//   - confirmation screen
//   - done screen
//
// Anything more complex should be introduced later through dedicated steps.
func (s *Service) renderOrderScreen(
	catalog Catalog,
	session Session,
) (ViewModel, bool) {
	switch session.Current {
	case ScreenOrderConfirm:
		ctx, ok := orderContext(catalog, session)
		if !ok {
			return buildDetailView(
				"Оформление заказа",
				"Не удалось собрать данные заказа. Вернитесь назад и попробуйте снова.",
				ActionBack,
			), true
		}

		return buildOrderConfirmView(ctx), true

	case ScreenOrderDone:
		return buildOrderDoneView(), true

	default:
		return ViewModel{}, false
	}
}

// orderContext returns current order context for order screens.
//
// When current screen is an order screen, source catalog state is resolved from
// session history. This keeps MVP implementation small and avoids adding
// separate order draft state into Session too early.
func orderContext(catalog Catalog, session Session) (OrderContext, bool) {
	screen, ok := orderSourceScreen(session)
	if !ok {
		return OrderContext{}, false
	}

	return catalogLeafOrderContext(catalog, screen)
}

// orderSourceScreen returns catalog screen that acts as source for current order flow.
//
// Resolution strategy:
//   - current catalog screen, if current screen is already a catalog screen
//   - otherwise the latest catalog screen from history
func orderSourceScreen(session Session) (ScreenID, bool) {
	if _, ok := parseCatalogScreen(session.Current); ok {
		return session.Current, true
	}

	for i := len(session.History) - 1; i >= 0; i-- {
		screen := session.History[i]
		if _, ok := parseCatalogScreen(screen); ok {
			return screen, true
		}
	}

	return "", false
}

// catalogLeafOrderContext extracts order context from catalog variant leaf screen.
//
// It validates that the provided screen is a catalog screen for selected
// variant leaf and then resolves all required labels from catalog tree.
func catalogLeafOrderContext(catalog Catalog, screen ScreenID) (OrderContext, bool) {
	path, ok := parseCatalogScreen(screen)
	if !ok || len(path) == 0 {
		return OrderContext{}, false
	}

	last, ok := path.Last()
	if !ok || last.Level != LevelVariant {
		return OrderContext{}, false
	}

	leaf, ok := catalog.FindNode(path)
	if !ok {
		return OrderContext{}, false
	}
	if len(leaf.Children) > 0 {
		return OrderContext{}, false
	}

	var ctx OrderContext

	for i, sel := range path {
		node, ok := catalog.FindNode(path[:i+1])
		if !ok {
			return OrderContext{}, false
		}

		switch sel.Level {
		case LevelCity:
			ctx.CityID = node.ID
			ctx.CityName = node.Label

		case LevelDistrict:
			ctx.DistrictID = node.ID
			ctx.DistrictName = node.Label

		case LevelProduct:
			ctx.ProductID = node.ID
			ctx.ProductLabel = node.Label

		case LevelVariant:
			ctx.VariantID = node.ID
			ctx.VariantLabel = node.Label
			ctx.BasePriceText = node.PriceText
		}
	}

	if ctx.CityID == "" || ctx.DistrictID == "" || ctx.ProductID == "" || ctx.VariantID == "" {
		return OrderContext{}, false
	}

	return ctx, true
}

// buildOrderConfirmView renders order confirmation screen for selected variant.
func buildOrderConfirmView(ctx OrderContext) ViewModel {
	priceText := ctx.BasePriceText
	if priceText == "" {
		priceText = "Цена уточняется"
	}

	text := fmt.Sprintf(
		"Оформление заказа\n\nГород: %s\nРайон: %s\nТовар: %s\nВариант: %s\nЦена: %s\n\nПодтвердить заявку?",
		ctx.CityName,
		ctx.DistrictName,
		ctx.ProductLabel,
		ctx.VariantLabel,
		priceText,
	)

	return ViewModel{
		Text: text,
		Inline: &InlineKeyboardView{
			Sections: []ActionSection{
				{
					Columns: 1,
					Actions: []ActionButton{
						{ID: ActionOrderConfirm, Label: "Подтвердить заказ"},
						{ID: ActionBack, Label: "Назад"},
					},
				},
			},
		},
		RemoveReply: true,
	}
}

// buildOrderDoneView renders terminal success screen for flow-only order request.
func buildOrderDoneView() ViewModel {
	return ViewModel{
		Text: "Заявка принята\n\nМы получили ваш запрос. Дальше можно вернуться в главное меню.",
		Inline: &InlineKeyboardView{
			Sections: []ActionSection{
				{
					Columns: 1,
					Actions: []ActionButton{
						{ID: ActionCatalogStart, Label: "В главное меню"},
					},
				},
			},
		},
		RemoveReply: true,
	}
}
