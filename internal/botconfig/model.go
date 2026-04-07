// Package botconfig
package botconfig

import "time"

// BotConfig describes editable bot configuration.
type BotConfig struct {
	ID                   string
	Name                 string
	Token                string
	DatabaseID           string
	StartScenario        string
	TelegramAdminUserIDs []int64
	IsEnabled            bool
	UpdatedAt            time.Time
}

// DatabaseProfile describes reusable database connection profile.
//
// DSN is infrastructure data and should not be exposed to operator UIs
// unless explicitly required by privileged admin workflows.
type DatabaseProfile struct {
	ID        string
	Name      string
	Driver    string
	DSN       string
	IsEnabled bool
	UpdatedAt time.Time
}
