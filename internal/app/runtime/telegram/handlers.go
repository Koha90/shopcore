package telegram

import (
	"context"
	"errors"

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

		if err := r.sendView(ctx, b, update.Message.Chat.ID, vm); err != nil {
			r.log.Error(
				"telegram send start view failed",
				"bot_id", spec.ID,
				"err", err,
			)
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

		if err := r.editView(ctx, b, msg, vm); err != nil {
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
			if err := r.sendView(ctx, b, update.Message.Chat.ID, vm); err != nil {
				r.log.Error(
					"telegram send reply action view failed",
					"bot_id", spec.ID,
					"err", err,
				)
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

		if err := r.sendView(ctx, b, update.Message.Chat.ID, vm); err != nil {
			r.log.Error(
				"telegram send reply action view failed",
				"bot_id", spec.ID,
				"err", err,
			)
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
