package telegram

import (
	"context"
	"fmt"

	tgbot "github.com/go-telegram/bot"

	"github.com/koha90/shopcore/internal/manager"
)

// syncBotMetadata fetches Telegram bot metadata and stores it in bot config.
//
// Metadata sync is best-effort: failure must not stop bot runtime startup.
func (r *Runner) syncBotMetadata(
	ctx context.Context,
	bot *tgbot.Bot,
	spec manager.BotSpec,
) error {
	if r == nil || r.botMetadataUpdater == nil {
		return nil
	}

	meta, err := fetchBotMetadata(ctx, bot)
	if err != nil {
		return err
	}

	if err := r.botMetadataUpdater.UpdateTelegramBotMetadata(
		ctx,
		spec.ID,
		meta.ID,
		meta.Username,
		meta.Name,
	); err != nil {
		return fmt.Errorf("update telegram bot metadata: %w", err)
	}

	return nil
}

func (r *Runner) applyBotMetadata(
	ctx context.Context,
	spec manager.BotSpec,
	meta BotMetadata,
) error {
	if r == nil || r.botMetadataUpdater == nil {
		return nil
	}

	return r.botMetadataUpdater.UpdateTelegramBotMetadata(
		ctx,
		spec.ID,
		meta.ID,
		meta.Username,
		meta.Name,
	)
}
