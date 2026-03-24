package manager

import "context"

// Runner executes a bot runtime.
//
// Run is expected to block until bot stops or fails.
// Manager passes a cancelable context to control lifecycle.
type Runner interface {
	Run(ctx context.Context, spec BotSpec, ready func()) error
}
