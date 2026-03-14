package botconfig

import "context"

// BotRepository stores bot configuration.
type BotRepository interface {
	Save(ctx context.Context, bot *BotConfig) error
	ByID(ctx context.Context, id string) (*BotConfig, error)
	List(ctx context.Context) ([]*BotConfig, error)
	Delete(ctx context.Context, id string) error
}

// DatabaseProfileRepository stores reusable database profiles.
type DatabaseProfileRepository interface {
	Save(ctx context.Context, profile *DatabaseProfile) error
	ByID(ctx context.Context, id string) (*DatabaseProfile, error)
	List(ctx context.Context) ([]*DatabaseProfile, error)
	Delete(ctx context.Context, id string) error
}
