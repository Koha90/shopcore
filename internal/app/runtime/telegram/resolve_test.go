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

func newFlowTestService() *flow.Service {
	return flow.NewService(nil)
}

func newFlowTestRunner() *Runner {
	return &Runner{
		log: slog.New(slog.NewTextHandler(io.Discard, nil)),
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

	svc := newFlowTestService()
	vm, ok, err := r.resolveReplyView(context.Background(), svc, spec, msg)

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

	svc := newFlowTestService()
	vm, ok, err := r.resolveReplyView(context.Background(), svc, spec, msg)

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

	svc := newFlowTestService()
	vm, actionID, ok, err := r.resolveCallbackView(context.Background(), svc, spec, cq)

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

	action := flow.ActionID("catalog:select:city:moscow")

	startCQ := &models.CallbackQuery{
		Data: encodeActionID(action),
		From: models.User{ID: 301},
		Message: models.MaybeInaccessibleMessage{
			Type: models.MaybeInaccessibleMessageTypeMessage,
			Message: &models.Message{
				Chat: models.Chat{ID: 401},
			},
		},
	}

	svc := newFlowTestService()
	vm, actionID, ok, err := r.resolveCallbackView(context.Background(), svc, spec, startCQ)

	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, action, actionID)
	require.Equal(t, "Москва\n\nВыберите категорию:", vm.Text)
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

	cityAction := flow.ActionID("catalog:select:city:moscow")

	svc := newFlowTestService()
	_, _, ok, err := r.resolveCallbackView(context.Background(), svc, spec, keyedCQ(cityAction))
	require.NoError(t, err)
	require.True(t, ok)

	vm, actionID, ok, err := r.resolveCallbackView(context.Background(), svc, spec, keyedCQ(flow.ActionBack))
	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, flow.ActionBack, actionID)
	require.Equal(t, "Каталог\n\nВыберите раздел:", vm.Text)
	require.NotNil(t, vm.Inline)
}
