package inmemory

import (
	"context"
	"sync"

	"github.com/koha90/shopcore/internal/botconfig"
)

// BotRepository stores bot configurations in process memory.
//
// It is intended for local development and tests.
type BotRepository struct {
	mu   sync.RWMutex
	bots map[string]*botconfig.BotConfig
}

// NewBotRepository creates a new in-memory bot configuration repository.
func NewBotRepository() *BotRepository {
	return &BotRepository{
		bots: make(map[string]*botconfig.BotConfig),
	}
}

// Save stores bot configuration in memory.
func (r *BotRepository) Save(ctx context.Context, bot *botconfig.BotConfig) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	cp := *bot
	r.bots[bot.ID] = &cp
	return nil
}

// ByID returns bot configuration by identifier.
func (r *BotRepository) ByID(ctx context.Context, id string) (*botconfig.BotConfig, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	bot, ok := r.bots[id]
	if !ok {
		return nil, botconfig.ErrBotNotFound
	}

	cp := *bot
	return &cp, nil
}

// List returns all stored bot configurations.
func (r *BotRepository) List(ctx context.Context) ([]botconfig.BotConfig, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]botconfig.BotConfig, 0, len(r.bots))
	for _, bot := range r.bots {
		result = append(result, *bot)
	}

	return result, nil
}

// Delete removes bot configuration by identifier.
func (r *BotRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.bots, id)
	return nil
}

// DatabaseProfileRepository stores database connection profile in process memory.
//
// It is intended for local development and tests.
type DatabaseProfileRepository struct {
	mu       sync.RWMutex
	profiles map[string]*botconfig.DatabaseProfile
}

// NewDatabaseProfileRepository creates a new in-memory database profile repository.
func NewDatabaseProfileRepository() *DatabaseProfileRepository {
	return &DatabaseProfileRepository{
		profiles: make(map[string]*botconfig.DatabaseProfile),
	}
}

// Save stores database profile in memory.
func (r *DatabaseProfileRepository) Save(ctx context.Context, profile *botconfig.DatabaseProfile) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	cp := *profile
	r.profiles[profile.ID] = &cp
	return nil
}

// ByID returns database profile by identifier.
func (r *DatabaseProfileRepository) ByID(ctx context.Context, id string) (*botconfig.DatabaseProfile, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	profile, ok := r.profiles[id]
	if !ok {
		return nil, botconfig.ErrDatabaseProfileNotFound
	}

	cp := *profile
	return &cp, nil
}

// List returns all stored database profiles.
func (r *DatabaseProfileRepository) List(ctx context.Context) ([]botconfig.DatabaseProfile, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]botconfig.DatabaseProfile, 0, len(r.profiles))
	for _, profile := range r.profiles {
		result = append(result, *profile)
	}

	return result, nil
}

// Delete removes database profile by identifier.
func (r *DatabaseProfileRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.profiles, id)
	return nil
}
