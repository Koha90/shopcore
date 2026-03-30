package telegram

import (
	"github.com/go-telegram/bot/models"

	"github.com/koha90/shopcore/internal/flow"
)

func messageSessionKey(botID string, msg *models.Message) flow.SessionKey {
	var userID int64
	if msg.From != nil {
		userID = msg.From.ID
	}

	return flow.SessionKey{
		BotID:  botID,
		ChatID: msg.Chat.ID,
		UserID: userID,
	}
}

func callbackSessionKey(botID string, cq *models.CallbackQuery) flow.SessionKey {
	var chatID int64

	switch cq.Message.Type {
	case models.MaybeInaccessibleMessageTypeMessage:
		if cq.Message.Message != nil {
			chatID = cq.Message.Message.Chat.ID
		}
	case models.MaybeInaccessibleMessageTypeInaccessibleMessage:
		if cq.Message.InaccessibleMessage != nil {
			chatID = cq.Message.InaccessibleMessage.Chat.ID
		}
	}
	return flow.SessionKey{
		BotID:  botID,
		ChatID: chatID,
		UserID: cq.From.ID,
	}
}
