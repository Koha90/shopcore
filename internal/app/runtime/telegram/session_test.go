package telegram

import (
	"testing"

	"github.com/go-telegram/bot/models"
	"github.com/stretchr/testify/require"

	"github.com/koha90/shopcore/internal/flow"
)

func TestMessageSessionKey(t *testing.T) {
	t.Run("with sender", func(t *testing.T) {
		msg := &models.Message{
			Chat: models.Chat{
				ID: 1001,
			},
			From: &models.User{
				ID: 2002,
			},
		}

		got := messageSessionKey("bot-1", msg)

		require.Equal(t, flow.SessionKey{
			BotID:  "bot-1",
			ChatID: 1001,
			UserID: 2002,
		}, got)
	})

	t.Run("without sender", func(t *testing.T) {
		msg := &models.Message{
			Chat: models.Chat{
				ID: 3003,
			},
			From: nil,
		}

		got := messageSessionKey("bot-2", msg)

		require.Equal(t, flow.SessionKey{
			BotID:  "bot-2",
			ChatID: 3003,
			UserID: 0,
		}, got)
	})
}

func TestCallbackSessionKey(t *testing.T) {
	t.Run("accessible message", func(t *testing.T) {
		cq := &models.CallbackQuery{
			From: models.User{
				ID: 4004,
			},
			Message: models.MaybeInaccessibleMessage{
				Type: models.MaybeInaccessibleMessageTypeMessage,
				Message: &models.Message{
					Chat: models.Chat{
						ID: 5005,
					},
				},
			},
		}

		got := callbackSessionKey("bot-3", cq)

		require.Equal(t, flow.SessionKey{
			BotID:  "bot-3",
			ChatID: 5005,
			UserID: 4004,
		}, got)
	})

	t.Run("inaccessible message", func(t *testing.T) {
		cq := &models.CallbackQuery{
			From: models.User{
				ID: 6006,
			},
			Message: models.MaybeInaccessibleMessage{
				Type: models.MaybeInaccessibleMessageTypeInaccessibleMessage,
				InaccessibleMessage: &models.InaccessibleMessage{
					Chat: models.Chat{
						ID: 7007,
					},
					MessageID: 88,
					Date:      0,
				},
			},
		}

		got := callbackSessionKey("bot-4", cq)

		require.Equal(t, flow.SessionKey{
			BotID:  "bot-4",
			ChatID: 7007,
			UserID: 6006,
		}, got)
	})

	t.Run("missing callback message payload", func(t *testing.T) {
		cq := &models.CallbackQuery{
			From: models.User{
				ID: 8008,
			},
			Message: models.MaybeInaccessibleMessage{
				Type: models.MaybeInaccessibleMessageTypeMessage,
			},
		}

		got := callbackSessionKey("bot-5", cq)

		require.Equal(t, flow.SessionKey{
			BotID:  "bot-5",
			ChatID: 0,
			UserID: 8008,
		}, got)
	})
}
