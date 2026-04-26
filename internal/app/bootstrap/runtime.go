package bootstrap

import (
	"context"
	"errors"

	"github.com/koha90/shopcore/internal/botconfig"
	"github.com/koha90/shopcore/internal/manager"
)

// RuntimeSource defines runtime bot configuration required for startup.
type RuntimeSource interface {
	ListEnabledRuntimeBots(ctx context.Context) ([]botconfig.RuntimeBot, error)
}

// RuntimeManager defines lifecycle operations required by bootstrap.
type RuntimeManager interface {
	Register(spec manager.BotSpec) error
	Start(ctx context.Context, id string) error
	Rename(id string, name string) error
}

// Result contains bootstrap result for a single bot.
type Result struct {
	ID         string
	Registered bool
	Started    bool
	Err        error
}

// Starter bootstraps bot runtimes from persistent configuration.
type Starter struct {
	source  RuntimeSource
	manager RuntimeManager
}

// NewStarter creates runtime bootstrap starter.
func NewStarter(source RuntimeSource, manager RuntimeManager) *Starter {
	if source == nil {
		panic("bootstrap: runtime source is nil")
	}
	if manager == nil {
		panic("bootstrap: runtime manager is nil")
	}

	return &Starter{
		source:  source,
		manager: manager,
	}
}

// StartEnabled loads enabled bots from configuration storage, registers them in
// manager, and starts their runtimes.
//
// Bootstrap is best-effort:
//   - duplicate registration is treated as non-fatal
//   - already running bot is treated as non-fatal
//   - all bots are attempted even if some fail
//
// Returned error if an aggregate of individual failures, if any.
func (s *Starter) StartEnabled(ctx context.Context) ([]Result, error) {
	bots, err := s.source.ListEnabledRuntimeBots(ctx)
	if err != nil {
		return nil, err
	}

	results := make([]Result, 0, len(bots))
	var errs []error

	for _, bot := range bots {
		result := Result{ID: bot.ID}

		spec := manager.BotSpec{
			ID:                   bot.ID,
			Name:                 bot.Name,
			Token:                bot.Token,
			DatabaseID:           bot.DatabaseID,
			StartScenario:        bot.StartScenario,
			TelegramAdminUserIDs: bot.TelegramAdminUserIDs,
			AdminOrdersChatID:    bot.AdminOrdersChatID,
			TelegramBotID:        bot.TelegramBotID,
			TelegramUsername:     bot.TelegramUsername,
			TelegramBotName:      bot.TelegramBotName,
		}

		err := s.manager.Register(spec)
		switch {
		case err == nil:
			result.Registered = true

		case errors.Is(err, manager.ErrDuplicateBotID):
			// Runtime may already be registered. Keep runtime name in sync with
			// canonical botconfig name and continue startup path.
			if renameErr := s.manager.Rename(bot.ID, bot.Name); renameErr != nil {
				result.Err = renameErr
				errs = append(errs, renameErr)
				results = append(results, result)
				continue
			}

		default:
			result.Err = err
			errs = append(errs, err)
			results = append(results, result)
			continue
		}

		err = s.manager.Start(ctx, bot.ID)
		switch {
		case err == nil:
			result.Started = true

		case errors.Is(err, manager.ErrBotAlreadyRunning):
			// Bot is already running. Treat as successful bootstrap outcome.
			result.Started = true

		default:
			result.Err = err
			errs = append(errs, err)
		}

		results = append(results, result)
	}

	return results, errors.Join(errs...)
}
