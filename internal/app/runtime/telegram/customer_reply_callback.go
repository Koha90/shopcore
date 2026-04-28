package telegram

import (
	"context"

	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/koha90/shopcore/internal/flow"
	"github.com/koha90/shopcore/internal/manager"
)

// handleAdminCustomerReplyStartCallback starts customer reply input in the
// admin private chat.
//
// The notification button may live in an admin channel or group, but telegram
// channel posts are not reliable plain user messages. Keeping text input in
// admin private chat preserves the normal flow SessionKey semantics.
func (r *Runner) handleAdminCustomerReplyStartCallback(
	ctx context.Context,
	b *tgbot.Bot,
	spec manager.BotSpec,
	svc *flow.Service,
	update *models.Update,
	actionID flow.ActionID,
) {
	if update == nil || update.CallbackQuery == nil {
		return
	}

	callback := update.CallbackQuery
	adminID := callback.From.ID

	if !r.canAdminTelegram(spec, adminID) {
		r.answerCallback(ctx, b, callback.ID, "Нет доступа")
		return
	}

	key := flow.SessionKey{
		BotID:  spec.ID,
		ChatID: adminID,
		UserID: adminID,
	}

	vm, err := svc.HandleAction(ctx, flow.ActionRequest{
		BotID:         spec.ID,
		BotName:       spec.Name,
		StartScenario: spec.StartScenario,
		ActionID:      actionID,
		SessionKey:    key,
		CanAdmin:      true,
	})
	if err != nil {
		r.answerCallback(ctx, b, callback.ID, "Не удалось открыть ответ")
		r.log.Error(
			"flow admin customer reply start failed",
			"bot_id", spec.ID,
			"admin_user_id", adminID,
			"action_id", actionID,
			"err", err,
		)
		return
	}

	activeID, err := r.sendView(ctx, b, adminID, vm)
	if err != nil {
		r.answerCallback(ctx, b, callback.ID, "Откройте личный чат с ботом")
		r.log.Error(
			"send admin customer reply prompt failed",
			"bot_id", spec.ID,
			"admin_user_id", adminID,
			"action_id", actionID,
			"err", err,
		)
		return
	}

	r.rememberActiveMessage(key, activeID)
	r.answerCallback(ctx, b, callback.ID, "Ответьте в личном чате с ботом")
}
