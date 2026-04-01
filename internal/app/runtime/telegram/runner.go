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

// FlowServiceFactory builds one flow service for one bot runtime instance.
//
// The factory allows runtime wiring to choose catalog provider per bot spec.
type FlowServiceFactory func(spec manager.BotSpec) *flow.Service

// Runner implements manager.Runner using Telegram Bot API.
type Runner struct {
	cfg         Config
	log         *slog.Logger
	flowFactory FlowServiceFactory
}

// NewRunner cunstructs Telegram runtime runner with default flow wiring.
func NewRunner(cfg Config, log *slog.Logger) *Runner {
	return NewRunnerWithFlowFactory(cfg, log, nil)
}

// NewRunnerWithFlowFactory constructs Telegram runtime runner
// with explicit flow service factory.
//
// This constructor is intended for tests and future wiring with
// per-bot catalog providers.
func NewRunnerWithFlowFactory(cfg Config, log *slog.Logger, factory FlowServiceFactory) *Runner {
	if log == nil {
		log = slog.Default()
	}
	if factory == nil {
		factory = func(spec manager.BotSpec) *flow.Service {
			return flow.NewService(nil)
		}
	}

	return &Runner{
		cfg:         cfg,
		log:         log,
		flowFactory: factory,
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

	svc := r.flowFactory(spec)
	if svc == nil {
		return fmt.Errorf("telegram flow service factory returned nil")
	}

	client, err := NewHTTPClient(r.cfg.ProxyURL)
	if err != nil {
		return fmt.Errorf("telegram http client: %w", err)
	}

	opts := []tgbot.Option{
		tgbot.WithHTTPClient(r.cfg.PollTimeout, client),
		tgbot.WithCheckInitTimeout(r.cfg.CheckInitTimeout),
		tgbot.WithDefaultHandler(r.defaultHandler(spec, svc)),
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
		r.startHandler(spec, svc),
	)

	b.RegisterHandler(
		tgbot.HandlerTypeCallbackQueryData,
		callbackPrefix,
		tgbot.MatchTypePrefix,
		r.callbackHandler(spec, svc),
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
