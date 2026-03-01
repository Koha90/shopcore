package services

// Logger defines minimal logging capabilities
// required by application services.
type Logger interface {
	Info(msg string, fields map[string]any)
	Error(msg string, fields map[string]any)
}
