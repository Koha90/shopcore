package botconfig

import (
	"context"
	"errors"
	"log/slog"
	"time"
)

var (
	ErrBotNotFound              = errors.New("bot config not found")
	ErrDatabaseProfileNotFound  = errors.New("database profile not found")
	ErrBotIDEmpty               = errors.New("bot id is required")
	ErrBotNameEmpty             = errors.New("bot name is required")
	ErrBotTokenEmpty            = errors.New("bot token is required")
	ErrBotStartScenarioEmpty    = errors.New("bot start scenario is required")
	ErrBotStartScenarioInvalid  = errors.New("bot start scenario is invalid")
	ErrDatabaseIDEmpty          = errors.New("database id is required")
	ErrDatabaseProfileIDEmpty   = errors.New("database profile id is required")
	ErrDatabaseProfileNameEmpty = errors.New("database profile name is required")
	ErrDatabaseDriverEmpty      = errors.New("database driver is required")
	ErrDatabaseDSNEmpty         = errors.New("database dsn is required")
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

// CreateDatabaseProfile creates a new reusable database profile.
func (s *Service) CreateDatabaseProfile(ctx context.Context, params CreateDatabaseProfileParams) error {
	if params.ID == "" {
		return ErrDatabaseProfileIDEmpty
	}
	if params.Name == "" {
		return ErrDatabaseProfileNameEmpty
	}
	if params.Driver == "" {
		return ErrDatabaseDriverEmpty
	}
	if params.DSN == "" {
		return ErrDatabaseDSNEmpty
	}

	profile := &DatabaseProfile{
		ID:        params.ID,
		Name:      params.Name,
		Driver:    params.Driver,
		DSN:       params.DSN,
		IsEnabled: params.IsEnabled,
		UpdatedAt: time.Now(),
	}

	return s.dbs.Save(ctx, profile)
}

// ListDatabaseProfiles returns safe views of all database profiles.
func (s *Service) ListDatabaseProfiles(ctx context.Context) ([]DatabaseProfileView, error) {
	profiles, err := s.dbs.List(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]DatabaseProfileView, 0, len(profiles))
	for _, profile := range profiles {
		result = append(result, toDatabaseProfileView(profile))
	}

	return result, nil
}

// DatabaseProfileByID returns safe view of databse profile by identifier.
func (s *Service) DatabaseProfileByID(ctx context.Context, id string) (DatabaseProfileView, error) {
	if id == "" {
		return DatabaseProfileView{}, ErrDatabaseProfileIDEmpty
	}

	profile, err := s.dbs.ByID(ctx, id)
	if err != nil {
		return DatabaseProfileView{}, err
	}

	return toDatabaseProfileView(*profile), nil
}

// CreateBot creates new bot configuration.
func (s *Service) CreateBot(ctx context.Context, params CreateBotParams) error {
	if params.ID == "" {
		return ErrBotIDEmpty
	}
	if params.Name == "" {
		return ErrBotNameEmpty
	}
	if params.Token == "" {
		return ErrBotTokenEmpty
	}
	if params.DatabaseID == "" {
		return ErrDatabaseIDEmpty
	}
	if params.StartScenario == "" {
		return ErrBotStartScenarioEmpty
	}
	if !IsValidStartScenario(params.StartScenario) {
		return ErrBotStartScenarioInvalid
	}

	dbProfile, err := s.dbs.ByID(ctx, params.DatabaseID)
	if err != nil || !dbProfile.IsEnabled {
		return ErrDatabaseProfileNotFound
	}

	bot := &BotConfig{
		ID:                   params.ID,
		Name:                 params.Name,
		Token:                params.Token,
		DatabaseID:           params.DatabaseID,
		StartScenario:        params.StartScenario,
		TelegramAdminUserIDs: params.TelegramAdminUserIDs,
		IsEnabled:            params.IsEnabled,
		UpdatedAt:            time.Now(),
	}

	return s.bots.Save(ctx, bot)
}

// UpdateBot updates editable bot configuration fields.
func (s *Service) UpdateBot(ctx context.Context, params UpdateBotParams) error {
	if params.ID == "" {
		return ErrBotIDEmpty
	}
	if params.Name == "" {
		return ErrBotNameEmpty
	}
	if params.DatabaseID == "" {
		return ErrDatabaseIDEmpty
	}
	if params.StartScenario == "" {
		return ErrBotStartScenarioEmpty
	}
	if !IsValidStartScenario(params.StartScenario) {
		return ErrBotStartScenarioInvalid
	}

	bot, err := s.bots.ByID(ctx, params.ID)
	if err != nil {
		return ErrBotNotFound
	}

	dbProfile, err := s.dbs.ByID(ctx, params.DatabaseID)
	if err != nil || !dbProfile.IsEnabled {
		return ErrDatabaseProfileNotFound
	}

	bot.Name = params.Name
	bot.DatabaseID = params.DatabaseID
	bot.StartScenario = params.StartScenario
	bot.IsEnabled = params.IsEnabled
	bot.TelegramAdminUserIDs = params.TelegramAdminUserIDs
	bot.UpdatedAt = time.Now()

	if params.Token != nil {
		if *params.Token == "" {
			return ErrBotTokenEmpty
		}
		bot.Token = *params.Token
	}

	return s.bots.Save(ctx, bot)
}

// ListBots returns safe views of all configured bots.
func (s *Service) ListBots(ctx context.Context) ([]BotView, error) {
	bots, err := s.bots.List(ctx)
	if err != nil {
		return nil, err
	}

	profiles, err := s.dbs.List(ctx)
	if err != nil {
		return nil, err
	}

	dbNames := make(map[string]string, len(profiles))
	for _, profile := range profiles {
		dbNames[profile.ID] = profile.Name
	}

	result := make([]BotView, 0, len(bots))
	for _, bot := range bots {
		result = append(result, toBotView(bot, dbNames[bot.DatabaseID]))
	}

	return result, nil
}

// BotByID returns safe view of bot configuration by identifier.
func (s *Service) BotByID(ctx context.Context, id string) (BotView, error) {
	if id == "" {
		return BotView{}, ErrBotIDEmpty
	}

	bot, err := s.bots.ByID(ctx, id)
	if err != nil {
		return BotView{}, ErrBotNotFound
	}

	profile, err := s.dbs.ByID(ctx, bot.DatabaseID)
	if err != nil {
		return BotView{}, ErrDatabaseProfileNotFound
	}

	return toBotView(*bot, profile.Name), nil
}

// BotToken returns raw bot token by identifier.
//
// This method is intended for internal operator workflows where token rotation
// or controlled reveal is required. Callers must treat returned value as secret.
func (s *Service) BotToken(ctx context.Context, id string) (string, error) {
	if id == "" {
		return "", ErrBotIDEmpty
	}

	bot, err := s.bots.ByID(ctx, id)
	if err != nil {
		return "", ErrBotNotFound
	}

	return bot.Token, nil
}

// UpdateBotToken replaces bot token without modifying other bot fields.
func (s *Service) UpdateBotToken(ctx context.Context, id string, token string) error {
	if id == "" {
		return ErrBotIDEmpty
	}
	if token == "" {
		return ErrBotTokenEmpty
	}

	bot, err := s.bots.ByID(ctx, id)
	if err != nil {
		return ErrBotNotFound
	}

	bot.Token = token
	bot.UpdatedAt = time.Now()

	return s.bots.Save(ctx, bot)
}
