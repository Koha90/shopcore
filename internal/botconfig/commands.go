package botconfig

// CreateBotParams contains data required to create bot configuration.
type CreateBotParams struct {
	ID            string
	Name          string
	Token         string
	DatabaseID    string
	StartScenario string
	IsEnabled     bool
}

// UpdateBotParams contains editable bot fields.
type UpdateBotParams struct {
	ID            string
	Name          string
	Token         *string
	DatabaseID    string
	StartScenario string
	IsEnabled     bool
}

// CreateDatabaseProfileParams contains data required to create database profile.
type CreateDatabaseProfileParams struct {
	ID        string
	Name      string
	Driver    string
	DSN       string
	IsEnabled bool
}

// UpdateDatabaseProfileParams contains editable profile fields.
type UpdateDatabaseProfileParams struct {
	ID        string
	Name      string
	Driver    string
	DSN       *string
	IsEnabled bool
}
