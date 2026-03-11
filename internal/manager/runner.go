package manager

import "context"

// Runner execute a bot runtime.
//
// Run is expected to block until bot stops or fails.
// Manager passes a cancelable context to control lifecycle.
type Runner interface {
	Run(ctx context.Context, spec BotSpec) error
}
