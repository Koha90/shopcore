package telegram

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	tgbot "github.com/go-telegram/bot"

	"github.com/koha90/shopcore/internal/manager"
)

// Runner implements manager.Runner using Telegram Bot API.
type Runner struct {
	cfg Config
	log *slog.Logger
}

// NewRunner implements manager.Runner using Telegram runtime runner.
func NewRunner(cfg Config, log *slog.Logger) *Runner {
	if log == nil {
		log = slog.Default()
	}

	return &Runner{
		cfg: cfg,
		log: log,
	}
}

// Run starts Telegram bot runtime and blocks until it stops or fails.
func (r *Runner) Run(ctx context.Context, spec manager.BotSpec, ready func()) error {
	client, err := NewHTTPClient(r.cfg.ProxyURL)
	if err != nil {
		return fmt.Errorf("telegram http client: %w", err)
	}

	var once sync.Once

	opts := []tgbot.Option{
		tgbot.WithHTTPClient(r.cfg.PollTimeout, client),
		tgbot.WithCheckInitTimeout(r.cfg.CheckInitTimeout),
		tgbot.WithDefaultHandler(r.defaultHandler(spec)),
		tgbot.WithErrorsHandler(r.errorsHandler(spec)),
	}

	b, err := tgbot.New(r.cfg.Token, opts...)
	if err != nil {
		return fmt.Errorf("create telegram bot: %w", err)
	}

	b.RegisterHandler(tgbot.HandlerTypeMessageText, "/start", tgbot.MatchTypeExact, r.startHandler(spec))

	once.Do(ready)

	r.log.Info("telegram runtime started", "bot_id", spec.ID, "name", spec.Name)

	b.Start(ctx)

	r.log.Info("telegram runtime stopped", "bot_id", spec.ID, "name", spec.Name)
	return nil
}
