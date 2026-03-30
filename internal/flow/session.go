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

// ScreenID identifies current logical screen in flow.
type ScreenID string

const (
	ScreenReplyWelcome ScreenID = "reply_welcome"
	ScreenRootCompact  ScreenID = "root_compact"
	ScreenRootExtended ScreenID = "root_extended"
	ScreenEntity1      ScreenID = "entity_1"
	ScreenEntity2      ScreenID = "entity_2"
	ScreenEntity3      ScreenID = "entity_3"
	ScreenEntity4      ScreenID = "entity_4"
	ScreenCabinet      ScreenID = "cabinet"
	ScreenSupport      ScreenID = "support"
	ScreenReviews      ScreenID = "reviews"
	ScreenBalance      ScreenID = "balance"
	ScreenBotsMine     ScreenID = "bots_mine"
	ScreenOrderLast    ScreenID = "order_last"
)

type Session struct {
	Current ScreenID
	History []ScreenID
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
