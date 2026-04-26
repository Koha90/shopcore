package telegram

import (
	"context"
	"errors"
	"fmt"
	"strings"

	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/koha90/shopcore/internal/flow"
	"github.com/koha90/shopcore/internal/manager"
)

// OrderNotificationMeta contains Telegram-side metadata attached to one order notification.
type OrderNotificationMeta struct {
	BotName     string
	BotUsername string

	UserID    int64
	ChatID    int64
	UserName  string
	UserLogin string
}

// notifyOrderConfirmed sends order notification to configured admin chat.
//
// Notification delivery belongs to Telegram runtime layer.
// Flow remains transport-agnostic and only provides resolved order context.
func (r *Runner) notifyOrderConfirmed(
	ctx context.Context,
	b *tgbot.Bot,
	spec manager.BotSpec,
	svc *flow.Service,
	key flow.SessionKey,
	meta OrderNotificationMeta,
) error {
	if spec.AdminOrdersChatID == 0 {
		return nil
	}
	if svc == nil {
		return fmt.Errorf("flow service is nil")
	}

	orderCtx, err := svc.CurrentOrderContext(ctx, key)
	if err != nil {
		if errors.Is(err, flow.ErrOrderContextUnavailable) {
			return nil
		}
		return fmt.Errorf("resolve order context: %w", err)
	}

	vm := buildAdminOrderNotificationView(spec, meta, orderCtx)

	if _, err := r.sendTextMessage(ctx, b, spec.AdminOrdersChatID, vm); err != nil {
		return fmt.Errorf("send admin order notification: %w", err)
	}

	return nil
}

// buildAdminOrderNotificationView builds one admin-facing order notification card.
func buildAdminOrderNotificationView(
	spec manager.BotSpec,
	meta OrderNotificationMeta,
	order flow.OrderContext,
) flow.ViewModel {
	var text strings.Builder

	text.WriteString("🛒 Новый заказ\n\n")

	text.WriteString("Бот: ")
	text.WriteString(formatBotLabel(spec))
	text.WriteString("\n")

	if meta.BotUsername != "" {
		text.WriteString("Username бота: @")
		text.WriteString(strings.TrimPrefix(meta.BotUsername, "@"))
		text.WriteString("\n")
	}

	if meta.UserName != "" {
		text.WriteString("Пользователь: ")
		text.WriteString(meta.UserName)
		text.WriteString("\n")
	}

	if meta.UserLogin != "" {
		text.WriteString("Логин: @")
		text.WriteString(strings.TrimPrefix(meta.UserLogin, "@"))
		text.WriteString("\n")
	}

	text.WriteString(fmt.Sprintf("User ID: %d\n", meta.UserID))
	text.WriteString(fmt.Sprintf("Chat ID: %d\n\n", meta.ChatID))

	text.WriteString("Город: ")
	text.WriteString(order.CityName)
	text.WriteString("\n")

	text.WriteString("Район: ")
	text.WriteString(order.DistrictName)
	text.WriteString("\n")

	text.WriteString("Товар: ")
	text.WriteString(order.ProductLabel)
	text.WriteString("\n")

	text.WriteString("Вариант: ")
	text.WriteString(order.VariantLabel)
	text.WriteString("\n")

	text.WriteString("Цена: ")
	if order.BasePriceText != "" {
		text.WriteString(order.BasePriceText)
	} else {
		text.WriteString("уточняется")
	}

	return flow.ViewModel{
		Text: text.String(),
	}
}

// buildTelegramDisplayName returns best-effort human-readable Telegram user name.
func buildTelegramDisplayName(user *models.User) string {
	if user == nil {
		return ""
	}

	first := strings.TrimSpace(user.FirstName)
	last := strings.TrimSpace(user.LastName)

	switch {
	case first != "" && last != "":
		return first + " " + last
	case first != "":
		return first
	case user.Username != "":
		return "@" + user.Username
	default:
		return ""
	}
}

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
