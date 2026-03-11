package manager

// Status represent current runtime state of a manager bot.
type Status string

const (
	// StatusRunning means bot start sequence has begun
	// but bot runtime is not yet considered stable.
	StatusStarting Status = "starting"

	// StatusRunning means bot runtime is active.
	StatusRunning Status = "running"

	// StatusStopping means stop sequence has begun.
	StatusStopping Status = "stopping"

	// StatusStopped means bot is currently not running.
	StatusStopped Status = "stopped"

	// StatusFailed means bot runtime exited with an error.
	StatusFailed Status = "failed"
)
