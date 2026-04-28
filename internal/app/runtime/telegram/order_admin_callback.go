package telegram

import (
	"context"

	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/koha90/shopcore/internal/manager"
	ordersvc "github.com/koha90/shopcore/internal/order/service"
)

func (r *Runner) handleAdminOrderCallback(
	ctx context.Context,
	b *tgbot.Bot,
	spec manager.BotSpec,
	update *models.Update,
	orderID int64,
	targetStatus ordersvc.OrderStatus,
) {
	msg, ok := callbackMessageContext(update)
	if !ok {
		r.answerCallback(ctx, b, update.CallbackQuery.ID, "message unavailable")
		return
	}

	if r.orderFactory == nil {
		r.answerCallback(ctx, b, update.CallbackQuery.ID, "order service unavailable")
	}

	orders, err := r.orderFactory(spec)
	if err != nil {
		r.answerCallback(ctx, b, update.CallbackQuery.ID, "order service failed")
		r.log.Error(
			"build order service",
			"order_id", orderID,
			"err", err,
		)
		return
	}
	if orders == nil {
		r.answerCallback(ctx, b, update.CallbackQuery.ID, "order service unavailable")
		return
	}

	if err := orders.UpdateStatus(ctx, orderID, targetStatus); err != nil {
		r.answerCallback(ctx, b, update.CallbackQuery.ID, "status update failed")
		r.log.Error(
			"update order status",
			"order_id", orderID,
			"target_status", targetStatus,
			"err", err,
		)
		return
	}

	order, err := orders.ByID(ctx, orderID)
	if err != nil {
		r.answerCallback(ctx, b, update.CallbackQuery.ID, "order reload failed")
		r.log.Error(
			"reload order",
			"order_id", orderID,
			"err", err,
		)
		return
	}

	vm := buildAdminOrderNotificationView(spec, order)

	if _, err := r.editView(ctx, b, msg, vm); err != nil {
		if isTelegramMessageNotModified(err) {
			r.answerCallback(ctx, b, update.CallbackQuery.ID, "")
			return
		}

		r.answerCallback(ctx, b, update.CallbackQuery.ID, "render failed")
		r.log.Error(
			"render update admin order card",
			"order_id", orderID,
			"err", err,
		)
		return
	}

	r.answerCallback(ctx, b, update.CallbackQuery.ID, "updated")
}
