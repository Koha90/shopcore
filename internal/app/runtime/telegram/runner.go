package telegram

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"sync"

	tgbot "github.com/go-telegram/bot"

	"github.com/koha90/shopcore/internal/flow"
	"github.com/koha90/shopcore/internal/manager"
)

// Runner implements manager.Runner using Telegram Bot API.
type Runner struct {
	cfg  Config
	log  *slog.Logger
	flow *flow.Service
}

// NewRunner implements manager.Runner using Telegram runtime runner.
func NewRunner(cfg Config, log *slog.Logger) *Runner {
	if log == nil {
		log = slog.Default()
	}

	return &Runner{
		cfg:  cfg,
		log:  log,
		flow: flow.NewService(),
	}
}

// Run starts Telegram bot runtime for a single managed bot instance.
//
// Bot token is taken from spec.Token.
// Shared runtime settings such as proxy and timeouts are taken Runner config.
func (r *Runner) Run(ctx context.Context, spec manager.BotSpec, ready func()) error {
	if strings.TrimSpace(spec.Token) == "" {
		return errors.New("telegram token is required")
	}

	client, err := NewHTTPClient(r.cfg.ProxyURL)
	if err != nil {
		return fmt.Errorf("telegram http client: %w", err)
	}

	opts := []tgbot.Option{
		tgbot.WithHTTPClient(r.cfg.PollTimeout, client),
		tgbot.WithCheckInitTimeout(r.cfg.CheckInitTimeout),
		tgbot.WithDefaultHandler(r.defaultHandler(spec)),
		tgbot.WithErrorsHandler(r.errorsHandler(spec)),
		tgbot.WithAllowedUpdates(tgbot.AllowedUpdates{
			"message",
			"callback_query",
		}),
		tgbot.WithNotAsyncHandlers(),
	}

	if r.cfg.Debug {
		opts = append(opts, tgbot.WithDebug())
	}

	b, err := tgbot.New(spec.Token, opts...)
	if err != nil {
		return fmt.Errorf("create telegram bot: %w", err)
	}

	b.RegisterHandler(
		tgbot.HandlerTypeMessageText,
		"/start",
		tgbot.MatchTypeExact,
		r.startHandler(spec),
	)

	b.RegisterHandler(
		tgbot.HandlerTypeCallbackQueryData,
		callbackPrefix,
		tgbot.MatchTypePrefix,
		r.callbackHandler(spec),
	)

	var once sync.Once
	once.Do(ready)

	r.log.Info(
		"telegram runtime started",
		"bot_id", spec.ID,
		"name", spec.Name,
	)

	b.Start(ctx)

	r.log.Info(
		"telegram runtime stopped",
		"bot_id", spec.ID,
		"name", spec.Name,
	)

	return nil
}
