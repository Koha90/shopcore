package telegram

import (
	"context"

	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/koha90/shopcore/internal/flow"
	"github.com/koha90/shopcore/internal/manager"
)

func (r *Runner) startHandler(spec manager.BotSpec) func(context.Context, *tgbot.Bot, *models.Update) {
	return func(ctx context.Context, b *tgbot.Bot, update *models.Update) {
		if update.Message == nil {
			return
		}

		vm, err := r.flow.Start(ctx, flow.StartRequest{
			BotID:         spec.ID,
			BotName:       spec.Name,
			StartScenario: spec.StartScenario,
		})
		if err != nil {
			r.log.Error(
				"flow start failed",
				"bot_id", spec.ID,
				"err", err,
			)
			return
		}

		if err := r.sendView(ctx, b, update.Message.Chat.ID, vm); err != nil {
			r.log.Error(
				"telegram send start view failed",
				"bot_id", spec.ID,
				"err", err,
			)
		}
	}
}

func (r *Runner) callbackHandler(spec manager.BotSpec) func(context.Context, *tgbot.Bot, *models.Update) {
	return func(ctx context.Context, b *tgbot.Bot, update *models.Update) {
		if update.CallbackQuery == nil {
			return
		}

		actionID, ok := decodeActionID(update.CallbackQuery.Data)
		if !ok {
			r.answerCallback(ctx, b, update.CallbackQuery.ID, "unknown action")
			r.log.Error(
				"telegram callback decode failed",
				"bot_id", spec.ID,
				"data", update.CallbackQuery.Data,
			)
			return
		}

		vm, err := r.flow.HandleAction(ctx, flow.ActionRequest{
			BotID:         spec.ID,
			BotName:       spec.Name,
			StartScenario: spec.StartScenario,
			ActionID:      actionID,
		})
		if err != nil {
			r.answerCallback(ctx, b, update.CallbackQuery.ID, "action failed")
			r.log.Error(
				"flow action failed",
				"bot_id", spec.ID,
				"action_id", actionID,
				"err", err,
			)
			return
		}

		chatID, messageID, ok := callbackMessageContext(update)
		if !ok {
			r.answerCallback(ctx, b, update.CallbackQuery.ID, "message unavailable")
			r.log.Error(
				"callback message context unavailable",
				"bot_id", spec.ID,
				"action_id", actionID,
			)
			return
		}

		if err := r.editView(ctx, b, chatID, messageID, vm); err != nil {
			r.answerCallback(ctx, b, update.CallbackQuery.ID, "render failed")
			r.log.Error(
				"telegram edit callback view failed",
				"bot_id", spec.ID,
				"action_id", actionID,
				"err", err,
			)
			return
		}

		r.answerCallback(ctx, b, update.CallbackQuery.ID, "")
	}
}

func (r *Runner) defaultHandler(spec manager.BotSpec) func(context.Context, *tgbot.Bot, *models.Update) {
	return func(ctx context.Context, b *tgbot.Bot, update *models.Update) {
		if update.Message == nil {
			return
		}

		actionID, ok := r.flow.ResolveReplyAction(update.Message.Text)
		if !ok {
			r.log.Info(
				"telegram update received",
				"bot_id", spec.ID,
				"chat_id", update.Message.Chat.ID,
				"text", update.Message.Text,
			)
			return
		}

		vm, err := r.flow.HandleAction(ctx, flow.ActionRequest{
			BotID:         spec.ID,
			BotName:       spec.Name,
			StartScenario: spec.StartScenario,
			ActionID:      actionID,
		})
		if err != nil {
			r.log.Error(
				"flow reply action failed",
				"bot_id", spec.ID,
				"action_id", actionID,
				"err", err,
			)
			return
		}

		if err := r.sendView(ctx, b, update.Message.Chat.ID, vm); err != nil {
			r.log.Error(
				"telegram send reply action view failed",
				"bot_id", spec.ID,
				"action_id", actionID,
				"err", err)
		}
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

// callbackMessageContext extracts editable message coordinates from callback query.
//
// go-telegram/bot wraps callback messages into MaybeInaccessibleMessage, so we
// only proceed when the accessible message is present.
func callbackMessageContext(update *models.Update) (chatID int64, messageID int, ok bool) {
	if update == nil || update.CallbackQuery == nil {
		return 0, 0, false
	}

	msg := update.CallbackQuery.Message.Message
	if msg == nil {
		return 0, 0, false
	}

	return msg.Chat.ID, msg.ID, true
}
