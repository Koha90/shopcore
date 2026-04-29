package telegram

import (
	"testing"

	"github.com/go-telegram/bot/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCallbackMessageContext(t *testing.T) {
	t.Run("nil update", func(t *testing.T) {
		msg, ok := callbackMessageContext(nil)

		require.False(t, ok)
		require.Nil(t, msg)
	})

	t.Run("nil callback query", func(t *testing.T) {
		update := &models.Update{}

		msg, ok := callbackMessageContext(update)

		require.False(t, ok)
		require.Nil(t, msg)
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

		msg, ok := callbackMessageContext(update)

		require.True(t, ok)
		require.NotNil(t, msg)
		require.Equal(t, int64(12345), msg.Chat.ID)
		require.Equal(t, 321, msg.ID)
	})

	t.Run("accessible type but nil message", func(t *testing.T) {
		update := &models.Update{
			CallbackQuery: &models.CallbackQuery{
				Message: models.MaybeInaccessibleMessage{
					Type: models.MaybeInaccessibleMessageTypeMessage,
				},
			},
		}

		msg, ok := callbackMessageContext(update)

		require.False(t, ok)
		require.Nil(t, msg)
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

		msg, ok := callbackMessageContext(update)

		require.False(t, ok)
		require.Nil(t, msg)
	})
}

func TestTelegramPhotoFileTokenReturnsLargestPhotoID(t *testing.T) {
	t.Parallel()

	got := telegramPhotoFileToken(&models.Message{
		Photo: []models.PhotoSize{
			{FileID: "small"},
			{FileID: "medium"},
			{FileID: "large"},
		},
	})

	assert.Equal(t, "large", got)
}

func TestTelegramPhotoFileTokenEmptyMessage(t *testing.T) {
	t.Parallel()

	assert.Empty(t, telegramPhotoFileToken(nil))
}
