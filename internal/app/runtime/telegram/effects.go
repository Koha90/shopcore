package telegram

import (
	"context"
	"fmt"
	"strings"

	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/koha90/shopcore/internal/flow"
)

// applyViewEffects execute transport-agnostic flow effects in Telegram runtime.
func (r *Runner) applyViewEffects(
	ctx context.Context,
	b *tgbot.Bot,
	vm flow.ViewModel,
) error {
	for _, effect := range vm.Effects {
		switch effect.Kind {
		case flow.EffectSendText:
			if err := r.applySendTextEffect(ctx, b, effect); err != nil {
				return err
			}

		case flow.EffectSendPhoto:
			if err := r.applySendPhotoEffect(ctx, b, effect); err != nil {
				return err
			}

		default:
			return fmt.Errorf("uknown flow effect kind %q", effect.Kind)
		}
	}

	return nil
}

func (r *Runner) applySendTextEffect(
	ctx context.Context,
	b *tgbot.Bot,
	effect flow.Effect,
) error {
	text := strings.TrimSpace(effect.Text)
	if text == "" {
		return nil
	}
	if effect.Target.ChatID == 0 {
		return fmt.Errorf("send text effect target chat id is empty")
	}

	_, err := r.sendTextMessage(ctx, b, effect.Target.ChatID, flow.ViewModel{
		Text: text,
	})
	if err != nil {
		return fmt.Errorf("send text effect: %w", err)
	}

	return nil
}

func (r *Runner) applySendPhotoEffect(
	ctx context.Context,
	b *tgbot.Bot,
	effect flow.Effect,
) error {
	const op = "send photo effect"

	if effect.Target.ChatID == 0 {
		return fmt.Errorf("%s target chat id is empty", op)
	}
	if effect.Media == nil {
		return fmt.Errorf("%s media is nil", op)
	}
	if effect.Media.Kind != flow.EffectMediaPhoto {
		return fmt.Errorf("%s media kind %q is unsupported", op, effect.Media.Kind)
	}

	fileToken := strings.TrimSpace(effect.Media.FileToken)
	if fileToken == "" {
		return fmt.Errorf("%s filw token is empty", op)
	}

	params := &tgbot.SendPhotoParams{
		ChatID: effect.Target.ChatID,
		Photo: &models.InputFileString{
			Data: fileToken,
		},
	}

	caption := strings.TrimSpace(effect.Text)
	if caption != "" {
		params.Caption = caption
	}

	if _, err := b.SendPhoto(ctx, params); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
