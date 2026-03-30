package telegram

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"github.com/go-telegram/bot/models"
	"github.com/stretchr/testify/require"

	"github.com/koha90/shopcore/internal/flow"
	"github.com/koha90/shopcore/internal/manager"
)

func newScenarioTestRunner() *Runner {
	return &Runner{
		log:  slog.New(slog.NewTextHandler(io.Discard, nil)),
		flow: flow.NewService(nil),
	}
}

func TestScenarioReplyWelcome_HistoryBack(t *testing.T) {
	r := newScenarioTestRunner()
	ctx := context.Background()

	spec := manager.BotSpec{
		ID:            "bot-reply",
		Name:          "Reply Bot",
		StartScenario: string(flow.StartScenarioReplyWelcome),
	}

	msg := &models.Message{
		Chat: models.Chat{ID: 1001},
		From: &models.User{ID: 2001},
		Text: "/start",
	}

	startVM, err := r.flow.Start(ctx, buildStartRequest(spec, msg))
	require.NoError(t, err)
	require.Equal(t, "Добро пожаловать 👋\nВыберите раздел:", startVM.Text)
	require.NotNil(t, startVM.Reply)
	require.Nil(t, startVM.Inline)
	require.False(t, startVM.RemoveReply)

	replyMsg := &models.Message{
		Chat: models.Chat{ID: 1001},
		From: &models.User{ID: 2001},
		Text: "♻️ Каталог",
	}

	rootVM, ok, err := r.resolveReplyView(ctx, spec, replyMsg)
	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, "Каталог\n\nВыберите раздел:", rootVM.Text)
	require.NotNil(t, rootVM.Inline)
	require.True(t, rootVM.RemoveReply)

	cityAction := flow.ActionID("catalog:select:city:moscow")

	entityCQ := &models.CallbackQuery{
		Data: encodeActionID(cityAction),
		From: models.User{ID: 2001},
		Message: models.MaybeInaccessibleMessage{
			Type: models.MaybeInaccessibleMessageTypeMessage,
			Message: &models.Message{
				Chat: models.Chat{ID: 1001},
			},
		},
	}

	entityVM, actionID, ok, err := r.resolveCallbackView(ctx, spec, entityCQ)
	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, cityAction, actionID)
	require.Equal(t, "Москва\n\nВыберите категорию:", entityVM.Text)

	backCQ := &models.CallbackQuery{
		Data: encodeActionID(flow.ActionBack),
		From: models.User{ID: 2001},
		Message: models.MaybeInaccessibleMessage{
			Type: models.MaybeInaccessibleMessageTypeMessage,
			Message: &models.Message{
				Chat: models.Chat{ID: 1001},
			},
		},
	}

	backVM, actionID, ok, err := r.resolveCallbackView(ctx, spec, backCQ)
	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, flow.ActionBack, actionID)
	require.Equal(t, "Каталог\n\nВыберите раздел:", backVM.Text)
	require.NotNil(t, backVM.Inline)

	require.Len(t, backVM.Inline.Sections, 1)
	require.Equal(t, 1, backVM.Inline.Sections[0].Columns)
}

func TestScenarioInlineCatalog_HistoryBack(t *testing.T) {
	r := newScenarioTestRunner()
	ctx := context.Background()

	spec := manager.BotSpec{
		ID:            "bot-inline",
		Name:          "Inline Bot",
		StartScenario: string(flow.StartScenarioInlineCatalog),
	}

	msg := &models.Message{
		Chat: models.Chat{ID: 3001},
		From: &models.User{ID: 4001},
		Text: "/start",
	}

	startVM, err := r.flow.Start(ctx, buildStartRequest(spec, msg))
	require.NoError(t, err)
	require.Equal(t, "Каталог\n\nВыберите раздел:", startVM.Text)
	require.NotNil(t, startVM.Inline)
	require.True(t, startVM.RemoveReply)

	require.Len(t, startVM.Inline.Sections, 2)
	require.Equal(t, 2, startVM.Inline.Sections[0].Columns)
	require.Equal(t, 1, startVM.Inline.Sections[1].Columns)

	cityAction := flow.ActionID("catalog:select:city:moscow")

	entityCQ := &models.CallbackQuery{
		Data: encodeActionID(cityAction),
		From: models.User{ID: 4001},
		Message: models.MaybeInaccessibleMessage{
			Type: models.MaybeInaccessibleMessageTypeMessage,
			Message: &models.Message{
				Chat: models.Chat{ID: 3001},
			},
		},
	}

	entityVM, actionID, ok, err := r.resolveCallbackView(ctx, spec, entityCQ)
	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, cityAction, actionID)
	require.Equal(t, "Москва\n\nВыберите категорию:", entityVM.Text)

	backCQ := &models.CallbackQuery{
		Data: encodeActionID(flow.ActionBack),
		From: models.User{ID: 4001},
		Message: models.MaybeInaccessibleMessage{
			Type: models.MaybeInaccessibleMessageTypeMessage,
			Message: &models.Message{
				Chat: models.Chat{ID: 3001},
			},
		},
	}

	backVM, actionID, ok, err := r.resolveCallbackView(ctx, spec, backCQ)
	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, flow.ActionBack, actionID)
	require.Equal(t, "Каталог\n\nВыберите раздел:", backVM.Text)
	require.NotNil(t, backVM.Inline)

	require.Len(t, backVM.Inline.Sections, 2)
	require.Equal(t, 2, backVM.Inline.Sections[0].Columns)
	require.Equal(t, 1, backVM.Inline.Sections[1].Columns)
}
