package telegram

import (
	"context"
	"errors"
	"strings"

	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/koha90/shopcore/internal/flow"
	"github.com/koha90/shopcore/internal/manager"
)

func (r *Runner) startHandler(
	spec manager.BotSpec,
	svc *flow.Service,
) func(context.Context, *tgbot.Bot, *models.Update) {
	return func(ctx context.Context, b *tgbot.Bot, update *models.Update) {
		if update.Message == nil {
			return
		}

		vm, err := r.resolveStartView(ctx, svc, spec, update.Message)
		if err != nil {
			r.log.Error(
				"flow start failed",
				"bot_id", spec.ID,
				"err", err,
			)
			return
		}

		activeID, err := r.sendView(ctx, b, update.Message.Chat.ID, vm)
		if err != nil {
			r.log.Error(
				"telegram send start view failed",
				"bot_id", spec.ID,
				"err", err,
			)
			return
		}

		if update.Message.From != nil {
			r.rememberActiveMessage(flow.SessionKey{
				BotID:  spec.ID,
				ChatID: update.Message.Chat.ID,
				UserID: update.Message.From.ID,
			}, activeID)
		}
	}
}

func (r *Runner) callbackHandler(
	spec manager.BotSpec,
	svc *flow.Service,
) func(context.Context, *tgbot.Bot, *models.Update) {
	return func(ctx context.Context, b *tgbot.Bot, update *models.Update) {
		if update.CallbackQuery == nil {
			return
		}

		actionID, ok := decodeCallbackActionID(update.CallbackQuery.Data)
		if ok {
			if orderID, targetStatus, ok := parseAdminOrderAction(actionID); ok {
				r.handleAdminOrderCallback(ctx, b, spec, update, orderID, targetStatus)
				return
			}
			if _, _, ok := flow.AdminCustomerReplyStartTarget(actionID); ok {
				r.handleAdminCustomerReplyStartCallback(ctx, b, spec, svc, update, actionID)
				return
			}
		}

		vm, actionID, ok, err := r.resolveCallbackView(ctx, svc, spec, update.CallbackQuery)
		if !ok {
			r.answerCallback(ctx, b, update.CallbackQuery.ID, "unknown action")
			r.log.Error(
				"telegram callback decode failed",
				"bot_id", spec.ID,
				"data", update.CallbackQuery.Data,
				"err", err,
			)
			return
		}

		msg, ok := callbackMessageContext(update)
		if !ok {
			r.answerCallback(ctx, b, update.CallbackQuery.ID, "message unavailable")
			r.log.Error(
				"callback message context unavailable",
				"bot_id", spec.ID,
				"action_id", actionID,
			)
			return
		}

		key := callbackSessionKey(spec.ID, update.CallbackQuery)

		if activeID, ok := r.activeMessageFor(key); ok && activeID != 0 && activeID != msg.ID {
			r.answerCallback(ctx, b, update.CallbackQuery.ID, "Это старый экран. Используйте последнее сообщение.")
			return
		}

		activeID, err := r.editView(ctx, b, msg, vm)
		if err != nil {
			if isTelegramMessageNotModified(err) {
				r.answerCallback(ctx, b, update.CallbackQuery.ID, "")
				return
			}

			r.answerCallback(ctx, b, update.CallbackQuery.ID, "render failed")
			r.log.Error(
				"telegram edit callback view failed",
				"bot_id", spec.ID,
				"action_id", actionID,
				"err", err,
			)
			return
		}

		r.log.Debug(
			"telegram callback resolved",
			"action_id", actionID,
			"chat_id", msg.Chat.ID,
			"message_id", msg.ID,
			"callback_data", update.CallbackQuery.Data,
		)

		r.log.Debug(
			"telegram callback render view",
			"action_id", actionID,
			"chat_id", msg.Chat.ID,
			"message_id", msg.ID,
			"view_has_media", hasImage(vm),
			"text_len", len(vm.Text),
		)

		r.rememberActiveMessage(key, activeID)
		r.answerCallback(ctx, b, update.CallbackQuery.ID, "")

		if actionID == flow.ActionOrderConfirm {
			meta := OrderNotificationMeta{
				UserID:    update.CallbackQuery.From.ID,
				ChatID:    msg.Chat.ID,
				UserName:  buildTelegramDisplayName(&update.CallbackQuery.From),
				UserLogin: update.CallbackQuery.From.Username,
			}

			order, persistErr := r.persistConfirmedOrder(ctx, spec, svc, key, meta)
			if persistErr != nil {
				r.log.Error(
					"persist confirmed order",
					"bot_id", spec.ID,
					"user_id", update.CallbackQuery.From.ID,
					"chat_id", msg.Chat.ID,
					"err", persistErr,
				)
				return
			}

			notifyErr := r.notifyOrderConfirmed(ctx, b, spec, order)
			if notifyErr != nil {
				r.log.Error(
					"send admin order notification",
					"bot_id", spec.ID,
					"admin_orders_chat_id", spec.AdminOrdersChatID,
					"user_id", update.CallbackQuery.From.ID,
					"chat_id", msg.Chat.ID,
					"err", notifyErr,
				)
			}
		}
	}
}

func (r *Runner) defaultHandler(
	spec manager.BotSpec,
	svc *flow.Service,
) func(context.Context, *tgbot.Bot, *models.Update) {
	return func(ctx context.Context, b *tgbot.Bot, update *models.Update) {
		if update.Message == nil {
			return
		}

		vm, ok, err := r.resolveReplyView(ctx, svc, spec, update.Message)
		if err != nil {
			r.log.Error(
				"flow reply action failed",
				"bot_id", spec.ID,
				"text", update.Message.Text,
				"err", err,
			)
			return
		}
		if ok {
			activeID, err := r.sendView(ctx, b, update.Message.Chat.ID, vm)
			if err != nil {
				r.log.Error(
					"telegram send reply action view failed",
					"bot_id", spec.ID,
					"err", err,
				)
				return
			}

			if update.Message.From != nil {
				r.rememberActiveMessage(flow.SessionKey{
					BotID:  spec.ID,
					ChatID: update.Message.Chat.ID,
					UserID: update.Message.From.ID,
				}, activeID)
			}
			return
		}

		if update.Message.From == nil {
			r.log.Info(
				"telegram update received",
				"bot_id", spec.ID,
				"chat_id", update.Message.Chat.ID,
				"text", update.Message.Text,
			)
			return
		}

		key := flow.SessionKey{
			BotID:  spec.ID,
			ChatID: update.Message.Chat.ID,
			UserID: update.Message.From.ID,
		}

		if !svc.HasPendingInput(key) {
			if err := r.notifyCustomerTextMessage(ctx, b, spec, update.Message); err != nil {
				r.log.Error(
					"send admin customer message notification",
					"bot_id", spec.ID,
					"admin_orders_chat_id", spec.AdminOrdersChatID,
					"user_id", update.Message.From.ID,
					"chat_id", update.Message.Chat.ID,
					"err", err,
				)
			}

			r.log.Info(
				"telegram update received",
				"bot_id", spec.ID,
				"chat_id", update.Message.Chat.ID,
				"text", update.Message.Text,
			)
			return
		}

		vm, err = r.resolveTextView(ctx, svc, spec, update.Message)
		if err != nil {
			r.log.Error(
				"flow text input failed",
				"bot_id", spec.ID,
				"text", update.Message.Text,
				"err", err,
			)
			return
		}

		if err := r.applyViewEffects(ctx, b, vm); err != nil {
			r.log.Error(
				"apply flow effects failed",
				"bot_id", spec.ID,
				"user_id", update.Message.From.ID,
				"chat_id", update.Message.Chat.ID,
				"err", err,
			)

			vm = flow.ViewModel{
				Text: "Не удалось отправить ответ пользователю. Проверте лог и попробуйте ещё раз.",
			}
		}

		activeID, err := r.sendView(ctx, b, update.Message.Chat.ID, vm)
		if err != nil {
			r.log.Error(
				"telegram send reply action view failed",
				"bot_id", spec.ID,
				"err", err,
			)
			return
		}

		r.rememberActiveMessage(key, activeID)
	}
}

func (r *Runner) errorsHandler(spec manager.BotSpec) func(error) {
	return func(err error) {
		r.log.Error(
			"telegram runtime error",
			"bot_id", spec.ID,
			"err", err,
		)
	}
}

func (r *Runner) answerCallback(ctx context.Context, b *tgbot.Bot, callbackID, text string) {
	params := &tgbot.AnswerCallbackQueryParams{
		CallbackQueryID: callbackID,
		Text:            text,
	}

	if _, err := b.AnswerCallbackQuery(ctx, params); err != nil {
		r.log.Error("telegram answer callback failed", "err", err)
	}
}

func (r *Runner) resolveTextView(
	ctx context.Context,
	svc *flow.Service,
	spec manager.BotSpec,
	msg *models.Message,
) (flow.ViewModel, error) {
	if msg == nil {
		return flow.ViewModel{}, errors.New("telegram message is nil")
	}
	if msg.From == nil {
		return flow.ViewModel{}, errors.New("telegram message sender is nil")
	}

	return svc.HandleText(ctx, flow.TextRequest{
		BotID:         spec.ID,
		BotName:       spec.Name,
		StartScenario: spec.StartScenario,
		Text:          msg.Text,
		SessionKey: flow.SessionKey{
			BotID:  spec.ID,
			ChatID: msg.Chat.ID,
			UserID: msg.From.ID,
		},
		CanAdmin: r.canAdminTelegram(spec, msg.From.ID),
	})
}

// callbackMessageContext extracts editable message coordinates from callback query.
//
// go-telegram/bot wraps callback messages into MaybeInaccessibleMessage, so we
// only proceed when the accessible message is present.
func callbackMessageContext(update *models.Update) (*models.Message, bool) {
	if update == nil || update.CallbackQuery == nil {
		return nil, false
	}

	msg := update.CallbackQuery.Message.Message
	if msg == nil {
		return nil, false
	}

	return msg, true
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

func isTelegramMessageNotModified(err error) bool {
	if err == nil {
		return false
	}

	return strings.Contains(err.Error(), "message is not modified")
}
