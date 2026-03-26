package botconfig

import (
	"context"
	"sort"
)

// RuntimeBot contains data required to register bot runtime.
//
// This model is internal-facing. Unlike BotView, it includes raw token becouse
// runtime bootstrap mus be able to start real bot processes.
type RuntimeBot struct {
	ID            string
	Name          string
	Token         string
	DatabaseID    string
	IsEnabled     bool
	StartScenario string
}

// RuntimePort defines internal bot runtime bootstrap use cases.
//
// This port is intentionally separate from operator-facing ServicePort becouse
// bootstrap requires raw runtime data, while operator interfaces should use
// safe view model.
type RuntimePort interface {
	ListEnabledRuntimeBots(ctx context.Context) ([]RuntimeBot, error)
}

// ListEnabledRuntimeBots returns all enabled bots required for runtime startup.
//
// Bots are returned in stable order by ID so that bootstrap is deterministic.
func (s *Service) ListEnabledRuntimeBots(ctx context.Context) ([]RuntimeBot, error) {
	bots, err := s.bots.List(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]RuntimeBot, 0, len(bots))
	for _, bot := range bots {
		if !bot.IsEnabled {
			continue
		}

		result = append(result, RuntimeBot{
			ID:            bot.ID,
			Name:          bot.Name,
			Token:         bot.Token,
			DatabaseID:    bot.DatabaseID,
			IsEnabled:     bot.IsEnabled,
			StartScenario: startScenarioForBot(bot.ID),
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].ID < result[j].ID
	})

	return result, nil
}

func startScenarioForBot(id string) string {
	switch id {
	case "shop-main":
		return "reply_welcome"
	case "slow-bot":
		return "inline_catalog"
	default:
		return "reply_welcome"
	}
}
