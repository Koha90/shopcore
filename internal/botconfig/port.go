package botconfig

import "context"

// ServicePort defines bot configuration use cases for operator interfaces.
type ServicePort interface {
	ListBots(ctx context.Context) ([]BotView, error)
	BotByID(ctx context.Context, id string) (BotView, error)
	CreateBot(ctx context.Context, params CreateBotParams) error
	UpdateBot(ctx context.Context, params UpdateBotParams) error

	ListDatabaseProfiles(ctx context.Context) ([]DatabaseProfileView, error)
	DatabaseProfileByID(ctx context.Context, id string) (DatabaseProfileView, error)
	CreateDatabaseProfile(ctx context.Context, params CreateDatabaseProfileParams) error
}
