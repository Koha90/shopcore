package telegram

import (
	"context"
	"errors"
	"fmt"

	"github.com/koha90/shopcore/internal/flow"
	"github.com/koha90/shopcore/internal/manager"
	ordersvc "github.com/koha90/shopcore/internal/order/service"
)

// OrderServiceFactory builds one order runtime service for one bot runtime instance.
//
// Factory uses bot runtime spec to choose the correct database wiring.
type OrderServiceFactory func(spec manager.BotSpec) (ordersvc.RuntimeService, error)

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
) (ordersvc.Order, error) {
	if r == nil || r.orderFactory == nil {
		return ordersvc.Order{}, nil
	}
	if svc == nil {
		return ordersvc.Order{}, fmt.Errorf("flow service is nil")
	}

	orders, err := r.orderFactory(spec)
	if err != nil {
		return ordersvc.Order{}, fmt.Errorf("build order creator: %w", err)
	}
	if orders == nil {
		return ordersvc.Order{}, nil
	}

	orderCtx, err := svc.CurrentOrderContext(ctx, key)
	if err != nil {
		if errors.Is(err, flow.ErrOrderContextUnavailable) {
			return ordersvc.Order{}, nil
		}
		return ordersvc.Order{}, fmt.Errorf("resolve order context: %w", err)
	}

	created, err := orders.Create(ctx, ordersvc.CreateOrderParams{
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
		return ordersvc.Order{}, fmt.Errorf("create order: %w", err)
	}

	order, err := orders.ByID(ctx, created.ID)
	if err != nil {
		return ordersvc.Order{}, fmt.Errorf("get created order: %w", err)
	}

	return order, nil
}
