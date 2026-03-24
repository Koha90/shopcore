package manager

import (
	"context"
	"sort"
	"sync"
)

// Manager coordinates lifecycle of registered bots.
//
// Manager is safe for concurrent use.
type Manager struct {
	runner  Runner
	mu      sync.RWMutex
	entries map[string]*Entry
}

// New creates a new Manager instance.
//
// runner must not be nil.
func New(runner Runner) *Manager {
	if runner == nil {
		panic("manager: Runner is nil")
	}

	return &Manager{
		runner:  runner,
		entries: make(map[string]*Entry),
	}
}

// Register adds a new bot to manager in stopped state.
//
// It does not start bot runtime.
func (m *Manager) Register(spec BotSpec) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.entries[spec.ID]; exists {
		return ErrDuplicateBotID
	}

	m.entries[spec.ID] = &Entry{
		spec:   spec,
		status: StatusStopped,
	}

	return nil
}

// Start launches bot runtime in a separate goroutine.
//
// Bot must be registered and not already be starting or running.
func (m *Manager) Start(ctx context.Context, id string) error {
	m.mu.Lock()

	entry, ok := m.entries[id]
	if !ok {
		m.mu.Unlock()
		return ErrBotNotFound
	}

	if entry.status == StatusStarting || entry.status == StatusRunning {
		m.mu.Unlock()
		return ErrBotAlreadyRunning
	}

	runCtx, cancel := context.WithCancel(ctx)
	done := make(chan struct{})

	entry.cancel = cancel
	entry.done = done
	entry.status = StatusStarting
	entry.lastError = nil

	spec := entry.spec
	m.mu.Unlock()

	go func() {
		defer close(done)

		readyOnce := sync.OnceFunc(func() {
			m.mu.Lock()
			defer m.mu.Unlock()

			if current, ok := m.entries[id]; ok && current.status == StatusStarting {
				current.status = StatusRunning
			}
		})

		err := m.runner.Run(runCtx, spec, readyOnce)

		m.mu.Lock()
		defer m.mu.Unlock()

		current, ok := m.entries[id]
		if !ok {
			return
		}

		current.cancel = nil
		current.done = nil

		if err != nil {
			current.status = StatusFailed
			current.lastError = err
			return
		}

		current.status = StatusStopped
		current.lastError = nil
	}()

	return nil
}

// Stop requests bot runtime shutdown.
//
// Bot must be registered and currently starting or running.
func (m *Manager) Stop(id string) error {
	m.mu.Lock()

	entry, ok := m.entries[id]
	if !ok {
		m.mu.Unlock()
		return ErrBotNotFound
	}

	if entry.status != StatusStarting && entry.status != StatusRunning {
		m.mu.Unlock()
		return ErrBotNotRunning
	}

	cancel := entry.cancel
	done := entry.done
	entry.status = StatusStopping
	m.mu.Unlock()

	if cancel != nil {
		cancel()
	}

	if done != nil {
		<-done
	}

	return nil
}

// Restart stops bot if it is running and then starts it again.
func (m *Manager) Restart(ctx context.Context, id string) error {
	status, err := m.Status(id)
	if err != nil {
		return err
	}

	if status == StatusStarting || status == StatusRunning {
		if err := m.Stop(id); err != nil {
			return err
		}
	}

	return m.Start(ctx, id)
}

// Status returns current runtime status for bot by ID.
func (m *Manager) Status(id string) (Status, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	entry, ok := m.entries[id]
	if !ok {
		return "", ErrBotNotFound
	}

	return entry.status, nil
}

// Info returns current read model for a single bot.
func (m *Manager) Info(id string) (Info, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	entry, ok := m.entries[id]
	if !ok {
		return Info{}, ErrBotNotFound
	}

	return infoFromEntry(entry), nil
}

// List returns all registered bots as sorted read models.
func (m *Manager) List() []Info {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]Info, 0, len(m.entries))
	for _, entry := range m.entries {
		result = append(result, infoFromEntry(entry))
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].ID < result[j].ID
	})

	return result
}

// Rename updates display name of registered bot runtime.
func (m *Manager) Rename(id string, name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	entry, ok := m.entries[id]
	if !ok {
		return ErrBotNotFound
	}

	entry.spec.Name = name
	return nil
}

func infoFromEntry(entry *Entry) Info {
	lastError := ""
	if entry.lastError != nil {
		lastError = entry.lastError.Error()
	}

	return Info{
		ID:        entry.spec.ID,
		Name:      entry.spec.Name,
		Status:    entry.status,
		LastError: lastError,
	}
}
