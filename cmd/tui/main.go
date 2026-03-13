package main

import (
	"context"
	"errors"
	"time"

	"botmanager/internal/manager"
	"botmanager/internal/tui"
)

type demoRunner struct{}

func (r *demoRunner) Run(ctx context.Context, spec manager.BotSpec) error {
	switch spec.ID {
	case "broken-bot":
		time.Sleep(700 * time.Millisecond)
		return errors.New("telegram auth failed")
	case "slow-bot":
		<-ctx.Done()
		return nil
	default:
		<-ctx.Done()
		return nil
	}
}

func main() {
	m := manager.New(&demoRunner{})

	must(m.Register(manager.BotSpec{
		ID:    "shop-main",
		Name:  "Shop Main",
		Token: "token-main",
	}))

	must(m.Register(manager.BotSpec{
		ID:    "slow-bot",
		Name:  "Slow Bot",
		Token: "token-slow",
	}))

	must(m.Register(manager.BotSpec{
		ID:    "broken-bot",
		Name:  "Broken Bot",
		Token: "token-broken",
	}))

	if err := tui.Run(m); err != nil {
		panic(err)
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
