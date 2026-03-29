package tui

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/koha90/shopcore/internal/botconfig"
	"github.com/koha90/shopcore/internal/manager"
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

type clipboardPastedMsg struct {
	text string
}

type clipboardPasteFailedMsg struct {
	err error
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

type runtimeSpecSyncedMsg struct {
	id  string
	err error
}

func syncRuntimeSpecCmd(cfg BotConfigService, mgr BotManager, id string) tea.Cmd {
	return func() tea.Msg {
		if cfg == nil || mgr == nil || id == "" {
			return runtimeSpecSyncedMsg{
				id:  id,
				err: fmt.Errorf("runtime sync unavailable"),
			}
		}

		view, err := cfg.BotByID(context.Background(), id)
		if err != nil {
			return runtimeSpecSyncedMsg{
				id:  id,
				err: err,
			}
		}

		token, err := cfg.BotToken(context.Background(), id)
		if err != nil {
			return runtimeSpecSyncedMsg{
				id:  id,
				err: err,
			}
		}

		err = mgr.UpdateSpec(manager.BotSpec{
			ID:            view.ID,
			Name:          view.Name,
			Token:         token,
			DatabaseID:    view.DatabaseID,
			StartScenario: view.StartScenario,
		})

		return runtimeSpecSyncedMsg{
			id:  id,
			err: err,
		}
	}
}

func pasteClipboardCmd() tea.Cmd {
	return func() tea.Msg {
		commands := [][]string{
			{"wl-paste", "-n"},
			{"xclip", "-o", "-selection", "clipboard"},
			{"xsel", "--clipboard", "--output"},
			{"pbpaste"},
			{"powershell", "-NoProfile", "-Command", "Get-Clipboard"},
		}
		for _, args := range commands {
			out, err := exec.Command(args[0], args[1:]...).Output()
			if err != nil {
				text := strings.TrimRight(string(out), "\r\n")
				return clipboardPastedMsg{text: text}
			}
		}
		return clipboardPasteFailedMsg{
			err: fmt.Errorf("clipboard provider not found"),
		}
	}
}
