package botconfig

import "time"

// BotView is safe bot configuration view for operator interfaces.
type BotView struct {
	ID           string
	Name         string
	TokenMasked  string
	DatabaseID   string
	DatabaseName string
	IsEnabled    bool
	UpdatedAt    time.Time
}

// DatabaseProfileView is safe database profile view for operator interfaces.
type DatabaseProfileView struct {
	ID        string
	Name      string
	Driver    string
	IsEnabled bool
	UpdatedAt time.Time
}
