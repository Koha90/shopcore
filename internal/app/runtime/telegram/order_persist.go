package telegram

import (
	"context"
	"errors"
	"fmt"

	"github.com/koha90/shopcore/internal/flow"
	"github.com/koha90/shopcore/internal/manager"
	ordersvc "github.com/koha90/shopcore/internal/order/service"
)

// OrderCreatorFactory builds one order creator for one bot runtime instance.
//
// Factory uses bot runtime spec to choose the correct database wiring.
type OrderCreatorFactory func(spec manager.BotSpec) (ordersvc.OrderCreator, error)

// persistConfirmedOrder stores confirmed order before sending admin notification.
//
// This keeps flow transport-agnostic and lets Telegram runtime orchestrate
// persistence and notification after successful confiramtion.
func (r *Runner) persistConfirmedOrder(
	ctx context.Context,
	spec manager.BotSpec,
	svc *flow.Service,
	key flow.SessionKey,
	meta OrderNotificationMeta,
) error {
	if r == nil || r.orderCreatorFactory == nil {
		return nil
	}
	if svc == nil {
		return fmt.Errorf("flow service is nil")
	}

	creator, err := r.orderCreatorFactory(spec)
	if err != nil {
		return fmt.Errorf("build order creator: %w", err)
	}
	if creator == nil {
		return nil
	}

	orderCtx, err := svc.CurrentOrderContext(ctx, key)
	if err != nil {
		if errors.Is(err, flow.ErrOrderContextUnavailable) {
			return nil
		}
		return fmt.Errorf("resolve order context: %w", err)
	}

	err = creator.Create(ctx, ordersvc.CreateOrderParams{
		BotID:        spec.ID,
		BotName:      spec.Name,
		ChatID:       meta.ChatID,
		UserID:       meta.UserID,
		UserName:     meta.UserName,
		UserUsername: meta.UserLogin,
		CityID:       orderCtx.CityID,
		CityName:     orderCtx.CityName,
		DistrictID:   orderCtx.DistrictID,
		DistrictName: orderCtx.DistrictName,
		ProductID:    orderCtx.ProductID,
		ProductName:  orderCtx.ProductLabel,
		VariantID:    orderCtx.VariantID,
		VariantName:  orderCtx.VariantLabel,
		PriceText:    orderCtx.BasePriceText,
	})
	if err != nil {
		return fmt.Errorf("create order: %w", err)
	}

	return nil
}
