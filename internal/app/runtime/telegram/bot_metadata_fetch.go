package telegram

import (
	"context"
	"fmt"

	tgbot "github.com/go-telegram/bot"
)

// fetchBotMetadata loads Telegram bot metadata through getMe.
func fetchBotMetadata(ctx context.Context, b *tgbot.Bot) (BotMetadata, error) {
	if b == nil {
		return BotMetadata{}, fmt.Errorf("telegram bot is nil")
	}

	me, err := b.GetMe(ctx)
	if err != nil {
		return BotMetadata{}, fmt.Errorf("telegram getMe: %w", err)
	}

	return BotMetadata{
		ID:       me.ID,
		Username: normalizeTelegramUsername(me.Username),
		Name:     me.FirstName,
	}, nil
}
