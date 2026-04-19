package runtimelog

import "sync"

const defaultPerBotLimit = 300

// Store keeps recent runtime log entries grouped by bot ID.
//
// The store is safe for concurrent use.
// Each bot keeps only the most recent N entries to avoid unbounded growth.
type Store struct {
	mu          sync.RWMutex
	perBotLimit int
	byBotID     map[string][]Entry
}

// NewStore constructs an in-memory runtime log store.
//
// If perBotLimit is not positive, a default limit is used.
func NewStore(perBotLimit int) *Store {
	if perBotLimit <= 0 {
		perBotLimit = defaultPerBotLimit
	}

	return &Store{
		perBotLimit: perBotLimit,
		byBotID:     make(map[string][]Entry),
	}
}

// Append stores one entry.
//
// Entries without BotID are ignored becouse TUI cannot route them
// to a concrete bot panel.
func (s *Store) Append(entry Entry) {
	if s == nil || entry.BotID == "" {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	items := append(s.byBotID[entry.BotID], entry)
	if len(items) > s.perBotLimit {
		items = items[len(items)-s.perBotLimit:]
	}

	s.byBotID[entry.BotID] = items
}

// List returns a copy of recent entries for one bot.
func (s *Store) List(botID string) []Entry {
	if s == nil || botID == "" {
		return nil
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	src := s.byBotID[botID]
	if len(src) == 0 {
		return nil
	}

	out := make([]Entry, len(src))
	copy(out, src)
	return out
}

// Clear removes all retained entries for one bot.
func (s *Store) Clear(botID string) {
	if s == nil || botID == "" {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.byBotID, botID)
}
