package telegram

import (
	"context"

	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/koha90/shopcore/internal/manager"
)

func (r *Runner) startHandler(spec manager.BotSpec) func(context.Context, *tgbot.Bot, *models.Update) {
	return func(ctx context.Context, b *tgbot.Bot, update *models.Update) {
		if update.Message == nil {
			return
		}

		_, err := b.SendMessage(ctx, &tgbot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Привет. Я бот платформы Shopcore  ",
		})
		if err != nil {
			r.log.Error(
				"telegram send start reply failed",
				"bot_id", spec.ID,
				"err", err,
			)
		}
	}
}

func (r *Runner) defaultHandler(spec manager.BotSpec) func(context.Context, *tgbot.Bot, *models.Update) {
	return func(ctx context.Context, b *tgbot.Bot, update *models.Update) {
		if update.Message == nil {
			return
		}

		r.log.Info(
			"telegram update received",
			"bot_id", spec.ID,
			"chat_id", update.Message.Chat.ID,
			"text", update.Message.Text,
		)
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
