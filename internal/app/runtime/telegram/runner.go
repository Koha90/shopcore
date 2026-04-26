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
type FlowServiceFactory func(spec manager.BotSpec) (*flow.Service, error)

// Runner implements manager.Runner using Telegram Bot API.
type Runner struct {
	cfg                 Config
	log                 *slog.Logger
	flowFactory         FlowServiceFactory
	orderCreatorFactory OrderCreatorFactory
	adminAccess         AdminAccessResolver
	activeMessageMu     sync.RWMutex
	activeMessageID     map[flow.SessionKey]int
}

// NewRunner cunstructs Telegram runtime runner with default flow wiring.
func NewRunner(cfg Config, log *slog.Logger) *Runner {
	return NewRunnerWithDeps(cfg, log, nil, nil, nil)
}

// NewRunnerWithFlowFactory constructs Telegram runtime runner
// with explicit flow service factory.
//
// This constructor is intended for tests and future wiring with
// per-bot catalog providers.
func NewRunnerWithFlowFactory(cfg Config, log *slog.Logger, factory FlowServiceFactory) *Runner {
	return NewRunnerWithDeps(cfg, log, factory, nil, nil)
}

// NewRunnerWithDeps constructs Telegram runtime runner with explicit runtime dependecies.
//
// It allows wiring a custom flow factory and Telegram admin access resolver.
func NewRunnerWithDeps(
	cfg Config,
	log *slog.Logger,
	flowFactory FlowServiceFactory,
	orderFactory OrderCreatorFactory,
	adminAccess AdminAccessResolver,
) *Runner {
	if log == nil {
		log = slog.Default()
	}
	if flowFactory == nil {
		flowFactory = func(spec manager.BotSpec) (*flow.Service, error) {
			return flow.NewService(nil), nil
		}
	}

	return &Runner{
		cfg:                 cfg,
		log:                 log,
		flowFactory:         flowFactory,
		orderCreatorFactory: orderFactory,
		adminAccess:         normalizeAdminAccessResolver(adminAccess),
	}
}

// Run starts Telegram bot runtime for a single managed bot instance.
//
// Bot token is taken from spec.Token.
// Shared runtime settings such as proxy and timeouts are taken Runner config.
func (r *Runner) Run(ctx context.Context, spec manager.BotSpec, ready func()) error {
	if r == nil {
		return fmt.Errorf("telegram runner is nil")
	}

	if strings.TrimSpace(spec.Token) == "" {
		return errors.New("telegram token is required")
	}

	runner := &Runner{
		cfg: r.cfg,
		log: r.log.With(
			"bot_id", spec.ID,
			"bot_name", spec.Name,
		),
		flowFactory:     r.flowFactory,
		adminAccess:     r.adminAccess,
		activeMessageID: make(map[flow.SessionKey]int),
	}

	svc, err := runner.flowFactory(spec)
	if err != nil {
		return fmt.Errorf("build flow service: %w", err)
	}
	if svc == nil {
		return fmt.Errorf("telegram flow service factory returned nil")
	}

	client, err := NewHTTPClient(runner.cfg.ProxyURL)
	if err != nil {
		return fmt.Errorf("telegram http client: %w", err)
	}

	opts := []tgbot.Option{
		tgbot.WithHTTPClient(runner.cfg.PollTimeout, client),
		tgbot.WithCheckInitTimeout(runner.cfg.CheckInitTimeout),
		tgbot.WithDefaultHandler(runner.defaultHandler(spec, svc)),
		tgbot.WithErrorsHandler(runner.errorsHandler(spec)),
		tgbot.WithAllowedUpdates(tgbot.AllowedUpdates{
			"message",
			"callback_query",
		}),
		tgbot.WithNotAsyncHandlers(),
	}

	if runner.cfg.Debug {
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
		runner.startHandler(spec, svc),
	)

	b.RegisterHandler(
		tgbot.HandlerTypeCallbackQueryData,
		callbackPrefix,
		tgbot.MatchTypePrefix,
		runner.callbackHandler(spec, svc),
	)

	// var once sync.Once
	// once.Do(ready)

	runner.log.Info("telegram runtime started")
	defer runner.log.Info("telegram runtime stopped")

	if ready != nil {
		ready()
	}

	b.Start(ctx)

	return nil
}
