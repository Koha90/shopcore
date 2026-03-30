package telegram

import (
	"testing"

	"github.com/go-telegram/bot/models"
	"github.com/stretchr/testify/require"
)

func TestCallbackMessageContext(t *testing.T) {
	t.Run("nil update", func(t *testing.T) {
		chatID, messageID, ok := callbackMessageContext(nil)

		require.False(t, ok)
		require.Zero(t, chatID)
		require.Zero(t, messageID)
	})

	t.Run("nil callback query", func(t *testing.T) {
		update := &models.Update{}

		chatID, messageID, ok := callbackMessageContext(update)

		require.False(t, ok)
		require.Zero(t, chatID)
		require.Zero(t, messageID)
	})

	t.Run("accessible message", func(t *testing.T) {
		update := &models.Update{
			CallbackQuery: &models.CallbackQuery{
				Message: models.MaybeInaccessibleMessage{
					Type: models.MaybeInaccessibleMessageTypeMessage,
					Message: &models.Message{
						ID: 321,
						Chat: models.Chat{
							ID: 12345,
						},
					},
				},
			},
		}

		chatID, messageID, ok := callbackMessageContext(update)

		require.True(t, ok)
		require.Equal(t, int64(12345), chatID)
		require.Equal(t, 321, messageID)
	})

	t.Run("accessible type but nil message", func(t *testing.T) {
		update := &models.Update{
			CallbackQuery: &models.CallbackQuery{
				Message: models.MaybeInaccessibleMessage{
					Type: models.MaybeInaccessibleMessageTypeMessage,
				},
			},
		}

		chatID, messageID, ok := callbackMessageContext(update)

		require.False(t, ok)
		require.Zero(t, chatID)
		require.Zero(t, messageID)
	})

	t.Run("inaccessible message", func(t *testing.T) {
		update := &models.Update{
			CallbackQuery: &models.CallbackQuery{
				Message: models.MaybeInaccessibleMessage{
					Type: models.MaybeInaccessibleMessageTypeInaccessibleMessage,
					InaccessibleMessage: &models.InaccessibleMessage{
						Chat: models.Chat{
							ID: 999,
						},
						MessageID: 111,
						Date:      0,
					},
				},
			},
		}

		chatID, messageID, ok := callbackMessageContext(update)

		require.False(t, ok)
		require.Zero(t, chatID)
		require.Zero(t, messageID)
	})
}
