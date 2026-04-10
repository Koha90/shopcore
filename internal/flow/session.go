package flow

import (
	"sync"
)

// SessionKey uniquely identifies user interaction context inside one bot.
type SessionKey struct {
	BotID  string
	ChatID int64
	UserID int64
}

// PendingInputPayload stores small flow-local continuation data for pending text input.
type PendingInputPayload map[string]string

// PendingInput stores one active text-input state inside session.
type PendingInput struct {
	Kind    PendingInputKind
	Payload PendingInputPayload
}

// Active reports whether session currently expects text input.
func (p PendingInput) Active() bool {
	return p.Kind != PendingInputNone
}

// Value returns one payload value by key.
func (p PendingInput) Value(key string) string {
	if p.Payload == nil {
		return ""
	}
	return p.Payload[key]
}

// SetValue stores one payload value by key.
func (p *PendingInput) SetValue(key, value string) {
	if p.Payload == nil {
		p.Payload = make(PendingInputPayload)
	}
	p.Payload[key] = value
}

// Session stores current screen, backward navigation history,
// pending input state and explicit admin access flag.
type Session struct {
	Current  ScreenID
	History  []ScreenID
	Pending  PendingInput
	CanAdmin bool
}

// Store defines session storage required by flow.
type Store interface {
	Get(key SessionKey) (Session, bool)
	Put(key SessionKey, session Session)
	Delete(key SessionKey)
}

// MemoryStore is in-memory session storage.
type MemoryStore struct {
	mu    sync.RWMutex
	items map[SessionKey]Session
}

// NewMemoryStore creates in-memory session store
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		items: make(map[SessionKey]Session),
	}
}

func (s *MemoryStore) Get(key SessionKey) (Session, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	v, ok := s.items[key]
	return v, ok
}

func (s *MemoryStore) Put(key SessionKey, session Session) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.items[key] = session
}

func (s *MemoryStore) Delete(key SessionKey) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.items, key)
}
