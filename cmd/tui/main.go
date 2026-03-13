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

	must(m.Register(manager.BotSpec{
		ID:    "1",
		Name:  "1",
		Token: "token-1",
	}))

	must(m.Register(manager.BotSpec{
		ID:    "2",
		Name:  "2",
		Token: "token-2",
	}))

	must(m.Register(manager.BotSpec{
		ID:    "3",
		Name:  "3",
		Token: "token-3",
	}))
	must(m.Register(manager.BotSpec{
		ID:    "4",
		Name:  "4",
		Token: "token-4",
	}))

	must(m.Register(manager.BotSpec{
		ID:    "5",
		Name:  "5",
		Token: "token-5",
	}))

	must(m.Register(manager.BotSpec{
		ID:    "6",
		Name:  "6",
		Token: "token-6",
	}))

	must(m.Register(manager.BotSpec{
		ID:    "7",
		Name:  "7",
		Token: "token-7",
	}))

	must(m.Register(manager.BotSpec{
		ID:    "8",
		Name:  "8",
		Token: "token-8",
	}))

	must(m.Register(manager.BotSpec{
		ID:    "9",
		Name:  "9",
		Token: "token-9",
	}))
	must(m.Register(manager.BotSpec{
		ID:    "10",
		Name:  "10",
		Token: "token-10",
	}))

	must(m.Register(manager.BotSpec{
		ID:    "11",
		Name:  "11",
		Token: "token-11",
	}))

	must(m.Register(manager.BotSpec{
		ID:    "12",
		Name:  "12",
		Token: "token-12",
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
