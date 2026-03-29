package tui

import (
	"context"
	"fmt"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/koha90/shopcore/internal/botconfig"
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

func loadBotConfigCmd(cfg BotConfigService, id string) tea.Cmd {
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

type databaseProfilesLoadedMsg struct {
	profiles []botconfig.DatabaseProfileView
	err      error
}

func loadDatabaseProfilesCmd(cfg BotConfigService) tea.Cmd {
	return func() tea.Msg {
		if cfg == nil {
			return databaseProfilesLoadedMsg{
				err: fmt.Errorf("config reader unavailable"),
			}
		}

		profiles, err := cfg.ListDatabaseProfiles(context.Background())
		return databaseProfilesLoadedMsg{
			profiles: profiles,
			err:      err,
		}
	}
}

type botConfigSavedMsg struct {
	id   string
	name string
	err  error
}

func saveBotConfigCmd(cfg BotConfigService, id string, form BotConfigEditForm) tea.Cmd {
	return func() tea.Msg {
		if cfg == nil || id == "" {
			return botConfigSavedMsg{
				err: fmt.Errorf("config service unavailable"),
			}
		}

		err := cfg.UpdateBot(context.Background(), botconfig.UpdateBotParams{
			ID:            id,
			Name:          form.Name,
			Token:         nil,
			DatabaseID:    form.DatabaseID,
			StartScenario: form.StartScenario,
			IsEnabled:     form.IsEnabled,
		})

		return botConfigSavedMsg{
			id:   id,
			name: form.Name,
			err:  err,
		}
	}
}

type editBotConfigLoadedMsg struct {
	config botconfig.BotView
	err    error
}

func loadEditBotConfigCmd(cfg BotConfigService, id string) tea.Cmd {
	return func() tea.Msg {
		if cfg == nil || id == "" {
			return editBotConfigLoadedMsg{
				err: fmt.Errorf("config service unavailable"),
			}
		}

		view, err := cfg.BotByID(context.Background(), id)
		return editBotConfigLoadedMsg{
			config: view,
			err:    err,
		}
	}
}

type botTokenSavedMsg struct {
	id string
}

type botTokenSaveFailedMsg struct {
	err error
}

func updateBotTokenCmd(cfg BotConfigService, id, token string) tea.Cmd {
	return func() tea.Msg {
		if err := cfg.UpdateBotToken(context.Background(), id, token); err != nil {
			return botTokenSaveFailedMsg{err: err}
		}
		return botTokenSavedMsg{id: id}
	}
}
