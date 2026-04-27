package telegram

import (
	"context"
	"fmt"
	"strings"

	tgbot "github.com/go-telegram/bot"

	"github.com/koha90/shopcore/internal/flow"
	"github.com/koha90/shopcore/internal/manager"
	ordersvc "github.com/koha90/shopcore/internal/order/service"
)

// notifyOrderConfirmed sends persisted order notification to configured admin chat.
func (r *Runner) notifyOrderConfirmed(
	ctx context.Context,
	b *tgbot.Bot,
	spec manager.BotSpec,
	order ordersvc.Order,
) error {
	if spec.AdminOrdersChatID == 0 {
		return nil
	}

	vm := buildAdminOrderNotificationView(spec, order)

	if _, err := r.sendTextMessage(ctx, b, spec.AdminOrdersChatID, vm); err != nil {
		return fmt.Errorf("send admin order notification: %w", err)
	}

	return nil
}

// buildAdminOrderNotificationView builds one admin-facing order notification card.
func buildAdminOrderNotificationView(
	spec manager.BotSpec,
	order ordersvc.Order,
) flow.ViewModel {
	var text strings.Builder

	text.WriteString("🛒 Новый заказ\n\n")
	text.WriteString(fmt.Sprintf("Заказ: #%d\n", order.ID))
	text.WriteString("Статус: ")
	text.WriteString(formatOrderStatusLabel(order.Status))
	text.WriteString("\n")

	text.WriteString("Бот: ")
	text.WriteString(formatBotLabel(spec))
	text.WriteString("\n\n")

	if order.UserName != "" {
		text.WriteString("Пользователь: ")
		text.WriteString(order.UserName)
		text.WriteString("\n")
	}

	if order.UserUsername != "" {
		text.WriteString("Логин: @")
		text.WriteString(strings.TrimPrefix(order.UserUsername, "@"))
		text.WriteString("\n")
	}

	text.WriteString(fmt.Sprintf("User ID: %d\n", order.UserID))
	text.WriteString(fmt.Sprintf("Chat ID: %d\n\n", order.ChatID))

	text.WriteString("Город: ")
	text.WriteString(order.CityName)
	text.WriteString("\n")

	text.WriteString("Район: ")
	text.WriteString(order.DistrictName)
	text.WriteString("\n")

	text.WriteString("Товар: ")
	text.WriteString(order.ProductName)
	text.WriteString("\n")

	text.WriteString("Вариант: ")
	text.WriteString(order.VariantName)
	text.WriteString("\n")

	text.WriteString("Цена: ")
	if order.PriceText != "" {
		text.WriteString(order.PriceText)
	} else {
		text.WriteString("уточняется")
	}

	actions := buildAdminOrderActions(order)

	vm := flow.ViewModel{
		Text: text.String(),
	}

	if len(actions) > 0 {
		vm.Inline = &flow.InlineKeyboardView{
			Sections: []flow.ActionSection{
				{
					Columns: 1,
					Actions: actions,
				},
			},
		}
	}

	return vm
}

// buildAdminOrderActions returns operator actions allowed for current order status.
func buildAdminOrderActions(order ordersvc.Order) []flow.ActionButton {
	switch order.Status {
	case ordersvc.OrderStatusNew:
		return []flow.ActionButton{
			{
				ID:    buildAdminOrderActionTake(order.ID),
				Label: "Взять в работу",
			},
			{
				ID:    buildAdminOrderActionClose(order.ID),
				Label: "Закрыть",
			},
		}

	case ordersvc.OrderStatusInProgress:
		return []flow.ActionButton{
			{
				ID:    buildAdminOrderActionClose(order.ID),
				Label: "Закрыть",
			},
		}

	default:
		return nil
	}
}

// formatBotLabel returns best available operator-facing bot label.
func formatBotLabel(spec manager.BotSpec) string {
	username := strings.TrimPrefix(strings.TrimSpace(spec.TelegramUsername), "@")

	switch {
	case username != "":
		return "@" + username
	case strings.TrimSpace(spec.TelegramBotName) != "":
		return strings.TrimSpace(spec.TelegramBotName)
	default:
		return strings.TrimSpace(spec.Name)
	}
}

// formatOrderStatusLabel renders operator-facing order status text.
func formatOrderStatusLabel(status ordersvc.OrderStatus) string {
	switch status {
	case ordersvc.OrderStatusNew:
		return "new"
	case ordersvc.OrderStatusInProgress:
		return "in_progress"
	case ordersvc.OrderStatusClosed:
		return "closed"
	default:
		return string(status)
	}
}
