package tui

import (
	"context"
	"fmt"
	"time"

	tea "charm.land/bubbletea/v2"

	"botmanager/internal/botconfig"
)

type actionResultMsg struct {
	message string
	err     error
}

type tickMsg time.Time

type botConfigLoadMsg struct {
	config botconfig.BotView
	err    error
}

func tickCmd() tea.Cmd {
	return tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func loadBotConfigCmd(cfg BotConfigReader, id string) tea.Cmd {
	return func() tea.Msg {
		if cfg == nil || id == "" {
			return botConfigLoadMsg{
				err: fmt.Errorf("config reader unavailable"),
			}
		}

		view, err := cfg.BotByID(context.Background(), id)
		return botConfigLoadMsg{
			config: view,
			err:    err,
		}
	}
}
