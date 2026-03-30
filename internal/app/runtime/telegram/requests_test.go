package telegram

import (
	"testing"

	"github.com/go-telegram/bot/models"
	"github.com/stretchr/testify/require"

	"github.com/koha90/shopcore/internal/flow"
	"github.com/koha90/shopcore/internal/manager"
)

func TestBuildStartRequest(t *testing.T) {
	spec := manager.BotSpec{
		ID:            "bot-1",
		Name:          "Shop Bot",
		StartScenario: "inline_catalog",
	}

	msg := &models.Message{
		Chat: models.Chat{
			ID: 101,
		},
		From: &models.User{
			ID: 202,
		},
	}

	got := buildStartRequest(spec, msg)

	require.Equal(t, flow.StartRequest{
		BotID:         "bot-1",
		BotName:       "Shop Bot",
		StartScenario: "inline_catalog",
		SessionKey: flow.SessionKey{
			BotID:  "bot-1",
			ChatID: 101,
			UserID: 202,
		},
	}, got)
}

func TestBuildMessageActionRequest(t *testing.T) {
	spec := manager.BotSpec{
		ID:            "bot-2",
		Name:          "Reply Bot",
		StartScenario: "reply_welcome",
	}

	msg := &models.Message{
		Chat: models.Chat{
			ID: 303,
		},
		From: &models.User{
			ID: 404,
		},
	}

	action := flow.ActionID("catalog:select:city:moscow")

	got := buildMessageActionRequest(spec, msg, action)

	require.Equal(t, flow.ActionRequest{
		BotID:         "bot-2",
		BotName:       "Reply Bot",
		StartScenario: "reply_welcome",
		ActionID:      action,
		SessionKey: flow.SessionKey{
			BotID:  "bot-2",
			ChatID: 303,
			UserID: 404,
		},
	}, got)
}

func TestBuildCallbackActionRequest_AccessibleMessage(t *testing.T) {
	spec := manager.BotSpec{
		ID:            "bot-3",
		Name:          "Inline Bot",
		StartScenario: "inline_catalog",
	}

	cq := &models.CallbackQuery{
		From: models.User{
			ID: 505,
		},
		Message: models.MaybeInaccessibleMessage{
			Type: models.MaybeInaccessibleMessageTypeMessage,
			Message: &models.Message{
				Chat: models.Chat{
					ID: 606,
				},
			},
		},
	}

	got := buildCallbackActionRequest(spec, cq, flow.ActionBack)

	require.Equal(t, flow.ActionRequest{
		BotID:         "bot-3",
		BotName:       "Inline Bot",
		StartScenario: "inline_catalog",
		ActionID:      flow.ActionBack,
		SessionKey: flow.SessionKey{
			BotID:  "bot-3",
			ChatID: 606,
			UserID: 505,
		},
	}, got)
}

func TestBuildCallbackActionRequest_InaccessibleMessage(t *testing.T) {
	spec := manager.BotSpec{
		ID:            "bot-4",
		Name:          "Inline Bot 2",
		StartScenario: "inline_catalog",
	}

	cq := &models.CallbackQuery{
		From: models.User{
			ID: 707,
		},
		Message: models.MaybeInaccessibleMessage{
			Type: models.MaybeInaccessibleMessageTypeInaccessibleMessage,
			InaccessibleMessage: &models.InaccessibleMessage{
				Chat: models.Chat{
					ID: 808,
				},
				MessageID: 12,
			},
		},
	}

	action := flow.ActionID("catalog:select:city:moscow")

	got := buildCallbackActionRequest(spec, cq, action)

	require.Equal(t, flow.ActionRequest{
		BotID:         "bot-4",
		BotName:       "Inline Bot 2",
		StartScenario: "inline_catalog",
		ActionID:      action,
		SessionKey: flow.SessionKey{
			BotID:  "bot-4",
			ChatID: 808,
			UserID: 707,
		},
	}, got)
}
