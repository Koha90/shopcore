package demo

import (
	"context"
	"errors"
	"time"

	"github.com/koha90/shopcore/internal/manager"
)

// Runner simulates bot runtime lifecycle for local development.
//
// It is intentionally simple:
//   - broken-bot fails during startup
//   - slow-bot becomes ready after a delay
//   - other bots become ready immediately
//
// This runner is useful while wiring storage, bootstrap, and TUI together
// before real Telegram runtime is connected.
type Runner struct{}

// NewRunner constructs a demo runtime runner.
func NewRunner() *Runner {
	return &Runner{}
}

// Run starts demo bot runtime and reports readiness through ready callback.
func (r *Runner) Run(ctx context.Context, spec manager.BotSpec, ready func()) error {
	switch spec.ID {
	case "broken-bot":
		select {
		case <-time.After(700 * time.Millisecond):
			return errors.New("telegram auth failed")
		case <-ctx.Done():
			return ctx.Err()
		}

	case "slow-bot":
		select {
		case <-time.After(4 * time.Second):
			ready()
		case <-ctx.Done():
			return ctx.Err()
		}

		<-ctx.Done()
		return nil

	default:
		ready()
		<-ctx.Done()
		return nil
	}
}
