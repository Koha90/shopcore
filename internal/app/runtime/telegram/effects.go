package telegram

import (
	"context"
	"fmt"
	"strings"

	tgbot "github.com/go-telegram/bot"

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
