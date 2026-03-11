package manager

// Info is a read model for external consumers such as TUI,
// admin panel, monitoring or HTTP API.
type Info struct {
	ID        string
	Name      string
	Status    Status
	LastError string
}
