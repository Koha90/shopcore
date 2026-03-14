package botconfig

import (
	"context"
	"errors"
	"log/slog"
)

var (
	ErrBotNotFound              error = errors.New("bot config not found")
	ErrDatabaseProfileNotFound  error = errors.New("database profile not found")
	ErrBotIDEmpty               error = errors.New("bot id is required")
	ErrBotNameEmpty             error = errors.New("bot name is required")
	ErrBotTokenEmpty            error = errors.New("bot token is required")
	ErrDatabaseIDEmpty          error = errors.New("database id is required")
	ErrDatabaseProfileIDEmpty   error = errors.New("database profile id is required")
	ErrDatabaseProfileNameEmpty error = errors.New("database profile name is required")
	ErrDatabaseDriverEmpty      error = errors.New("database driver is required")
	ErrDatabaseDSNEmpty         error = errors.New("database dsn is required")
)

// Service orchestrate bot and database profile configuration use cases.
type Service struct {
	bots   BotRepository
	dbs    DatabaseProfileRepository
	logger *slog.Logger
}

// NewService creates configuration service instance.
func NewService(bots BotRepository, dbs DatabaseProfileRepository, logger *slog.Logger) *Service {
	if bots == nil {
		panic("botconfig: BotRepository is nil")
	}
	if dbs == nil {
		panic("botconfig: DatabaseProfileRepository is nil")
	}
	if logger == nil {
		logger = slog.Default()
	}

	return &Service{
		bots:   bots,
		dbs:    dbs,
		logger: logger,
	}
}

// CreateBot ...
func (s *Service) CreateBot(ctx context.Context, params CreateBotParams) error {
	return nil
}
