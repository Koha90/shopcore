package manager

import "context"

// BotSpec describes a single managed bot instance.
//
// ID must be unique inside manager.
// Token is runtime configuration and must not be treated
// as bot identity.
type BotSpec struct {
	ID         string
	Name       string
	Token      string
	DatabaseID string

	StartScenario string

	TelegramAdminUserIDs []int64
	AdminOrdersChatID    int64

	TelegramBotID    int64
	TelegramUsername string
	TelegramBotName  string
}

// Entry stores bot runtime state managed by Manager.
type Entry struct {
	spec      BotSpec
	cancel    context.CancelFunc
	done      chan struct{}
	status    Status
	lastError error
}
