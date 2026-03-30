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

func newFlowTestRunner() *Runner {
	return &Runner{
		log:  slog.New(slog.NewTextHandler(io.Discard, nil)),
		flow: flow.NewService(nil),
	}
}

func TestResolveReplyView_UnknownText(t *testing.T) {
	r := newFlowTestRunner()

	spec := manager.BotSpec{
		ID:            "bot-1",
		Name:          "Shop Bot",
		StartScenario: string(flow.StartScenarioReplyWelcome),
	}

	msg := &models.Message{
		Text: "что-то левое",
		Chat: models.Chat{ID: 100},
		From: &models.User{ID: 200},
	}

	vm, ok, err := r.resolveReplyView(context.Background(), spec, msg)

	require.NoError(t, err)
	require.False(t, ok)
	require.Equal(t, flow.ViewModel{}, vm)
}

func TestResolveReplyView_KnownText(t *testing.T) {
	r := newFlowTestRunner()

	spec := manager.BotSpec{
		ID:            "bot-2",
		Name:          "Shop Bot",
		StartScenario: string(flow.StartScenarioReplyWelcome),
	}

	msg := &models.Message{
		Text: "♻️ Каталог",
		Chat: models.Chat{ID: 101},
		From: &models.User{ID: 201},
	}

	vm, ok, err := r.resolveReplyView(context.Background(), spec, msg)

	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, "Каталог\n\nВыберите раздел:", vm.Text)
	require.NotNil(t, vm.Inline)
	require.True(t, vm.RemoveReply)
}

func TestResolveCallbackView_InvalidPayload(t *testing.T) {
	r := newFlowTestRunner()

	spec := manager.BotSpec{
		ID:            "bot-3",
		Name:          "Inline Bot",
		StartScenario: string(flow.StartScenarioInlineCatalog),
	}

	cq := &models.CallbackQuery{
		Data: "bad-payload",
		From: models.User{ID: 300},
		Message: models.MaybeInaccessibleMessage{
			Type: models.MaybeInaccessibleMessageTypeMessage,
			Message: &models.Message{
				Chat: models.Chat{ID: 400},
			},
		},
	}

	vm, actionID, ok, err := r.resolveCallbackView(context.Background(), spec, cq)

	require.NoError(t, err)
	require.False(t, ok)
	require.Empty(t, actionID)
	require.Equal(t, flow.ViewModel{}, vm)
}

func TestResolveCallbackView_ValidPayload(t *testing.T) {
	r := newFlowTestRunner()

	spec := manager.BotSpec{
		ID:            "bot-4",
		Name:          "Inline Bot",
		StartScenario: string(flow.StartScenarioInlineCatalog),
	}

	startCQ := &models.CallbackQuery{
		Data: encodeActionID(flow.ActionEntity1),
		From: models.User{ID: 301},
		Message: models.MaybeInaccessibleMessage{
			Type: models.MaybeInaccessibleMessageTypeMessage,
			Message: &models.Message{
				Chat: models.Chat{ID: 401},
			},
		},
	}

	vm, actionID, ok, err := r.resolveCallbackView(context.Background(), spec, startCQ)

	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, flow.ActionEntity1, actionID)
	require.Equal(t, "Москва\n\nЗдесь будет следующий шаг сценария для выбранной сущности.", vm.Text)
	require.NotNil(t, vm.Inline)
}

func TestResolveCallbackView_BackUsesHistory(t *testing.T) {
	r := newFlowTestRunner()

	spec := manager.BotSpec{
		ID:            "bot-5",
		Name:          "Inline Bot",
		StartScenario: string(flow.StartScenarioInlineCatalog),
	}

	keyedCQ := func(action flow.ActionID) *models.CallbackQuery {
		return &models.CallbackQuery{
			Data: encodeActionID(action),
			From: models.User{ID: 777},
			Message: models.MaybeInaccessibleMessage{
				Type: models.MaybeInaccessibleMessageTypeMessage,
				Message: &models.Message{
					Chat: models.Chat{ID: 888},
				},
			},
		}
	}

	_, _, ok, err := r.resolveCallbackView(context.Background(), spec, keyedCQ(flow.ActionEntity1))
	require.NoError(t, err)
	require.True(t, ok)

	vm, actionID, ok, err := r.resolveCallbackView(context.Background(), spec, keyedCQ(flow.ActionBack))
	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, flow.ActionBack, actionID)
	require.Equal(t, "Каталог\n\nВыберите раздел:", vm.Text)
	require.NotNil(t, vm.Inline)
}
