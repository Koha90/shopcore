package manager

import "errors"

var (

	// ErrDuplicateBotID is returned when a bot with the same ID
	// is already registreted in manager.
	ErrDuplicateBotID = errors.New("duplicate bot id")

	// ErrBotNotFound is returned when manager cannot find
	// a bot by the provided ID.
	ErrBotNotFound = errors.New("bot not found")

	// ErrBotAlreadyRunning is returned when start is requested
	// for a bot that is already starting or running.
	ErrBotAlreadyRunning = errors.New("bot already running")

	// ErrBotNotRunning is returned when stop is requested
	// for a bot that is not currently starting or running.
	ErrBotNotRunning = errors.New("bot not running")
)
