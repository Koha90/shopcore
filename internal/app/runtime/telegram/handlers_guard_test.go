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

func newTestRunner() *Runner {
	return &Runner{
		log: slog.New(slog.NewTextHandler(io.Discard, nil)),
	}
}

func TestStartHandler_IgnoresNilMessage(t *testing.T) {
	r := newTestRunner()
	svc := flow.NewService(nil)
	h := r.startHandler(manager.BotSpec{
		ID:            "bot-1",
		Name:          "Shop Bot",
		StartScenario: "reply_welcome",
	}, svc)

	require.NotPanics(t, func() {
		h(context.Background(), nil, &models.Update{})
	})
}

func TestCallbackHandler_IgnoresNilCallbackQuery(t *testing.T) {
	r := newTestRunner()
	svc := flow.NewService(nil)
	h := r.callbackHandler(manager.BotSpec{
		ID:            "bot-2",
		Name:          "Shop Bot",
		StartScenario: "inline_catalog",
	}, svc)

	require.NotPanics(t, func() {
		h(context.Background(), nil, &models.Update{})
	})
}

func TestDefaultHandler_IgnoresNilMessage(t *testing.T) {
	r := newTestRunner()
	svc := flow.NewService(nil)
	h := r.defaultHandler(manager.BotSpec{
		ID:            "bot-3",
		Name:          "Shop Bot",
		StartScenario: "reply_welcome",
	}, svc)

	require.NotPanics(t, func() {
		h(context.Background(), nil, &models.Update{})
	})
}

func TestDefaultHandler_IgnoresUnknownReplyText(t *testing.T) {
	r := newTestRunner()
	svc := flow.NewService(nil)
	h := r.defaultHandler(manager.BotSpec{
		ID:            "bot-4",
		Name:          "Shop Bot",
		StartScenario: "reply_welcome",
	}, svc)

	update := &models.Update{
		Message: &models.Message{
			Text: "какая-то непонятная команда",
			Chat: models.Chat{
				ID: 101,
			},
		},
	}

	require.NotPanics(t, func() {
		h(context.Background(), nil, update)
	})
}

func TestErrorsHandler_DoesNotPanic(t *testing.T) {
	r := newTestRunner()
	h := r.errorsHandler(manager.BotSpec{
		ID:   "bot-5",
		Name: "Shop Bot",
	})

	require.NotPanics(t, func() {
		h(assertErr{})
	})
}

type assertErr struct{}

func (assertErr) Error() string {
	return "boom"
}
