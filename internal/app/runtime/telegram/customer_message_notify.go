package telegram

import (
	"context"
	"fmt"
	"strings"

	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/koha90/shopcore/internal/flow"
	"github.com/koha90/shopcore/internal/manager"
)

// customerMessageNotification carries on plain customer text message to the
// operator-facing notification builder.
type customerMessageNotification struct {
	UserID       int64
	ChatID       int64
	UserName     string
	UserUsername string
	Text         string
}

// notifyCustomerTextMessage sends a plain customer text message to the
// configured admin chat.
//
// Admin users are intentionally ignored here. Their text message are either
// handled as pending admin input or left untouched for future admin commands.
func (r *Runner) notifyCustomerTextMessage(
	ctx context.Context,
	b *tgbot.Bot,
	spec manager.BotSpec,
	msg *models.Message,
) error {
	if spec.AdminOrdersChatID == 0 || msg == nil || msg.From == nil {
		return nil
	}
	if r.canAdminTelegram(spec, msg.From.ID) {
		return nil
	}

	text := strings.TrimSpace(msg.Text)
	if text == "" {
		return nil
	}

	vm := buildAdminCustomerMessageNotificationView(
		spec,
		customerMessageNotification{
			UserID:       msg.From.ID,
			ChatID:       msg.Chat.ID,
			UserName:     buildTelegramDisplayName(msg.From),
			UserUsername: msg.From.Username,
			Text:         text,
		},
	)

	if _, err := r.sendTextMessage(ctx, b, spec.AdminOrdersChatID, vm); err != nil {
		return fmt.Errorf("send admin customer message notification: %w", err)
	}

	return nil
}

func buildAdminCustomerReplyActions(chatID, userID int64) []flow.ActionButton {
	return []flow.ActionButton{
		{
			ID:    flow.AdminCustomerReplyStartAction(chatID, userID),
			Label: "Ответить текстом",
		},
		{
			ID:    flow.AdminCustomerPhotoReplyStartAction(chatID, userID),
			Label: "Ответить с фото",
		},
	}
}

// buildAdminCustomerMessageNotificationView builds one admin-facing customer
// message card.
func buildAdminCustomerMessageNotificationView(
	spec manager.BotSpec,
	message customerMessageNotification,
) flow.ViewModel {
	var text strings.Builder

	text.WriteString("💬 Новое сообщение\n\n")
	text.WriteString("Бот: ")
	text.WriteString(formatBotLabel(spec))
	text.WriteString("\n\n")

	if message.UserName != "" {
		text.WriteString("Пользователь: ")
		text.WriteString(message.UserName)
		text.WriteString("\n")
	}

	if message.UserUsername != "" {
		text.WriteString("Логин: @")
		text.WriteString(strings.TrimPrefix(message.UserUsername, "@"))
		text.WriteString("\n")
	}

	text.WriteString(fmt.Sprintf("User ID: %d\n", message.UserID))
	text.WriteString(fmt.Sprintf("Chat ID: %d\n", message.ChatID))

	text.WriteString("Сообщение:\n")
	text.WriteString(message.Text)

	return flow.ViewModel{
		Text: text.String(),
		Inline: &flow.InlineKeyboardView{
			Sections: []flow.ActionSection{
				{
					Columns: 1,
					Actions: buildAdminCustomerReplyActions(
						message.ChatID,
						message.UserID,
					),
				},
			},
		},
	}
}
