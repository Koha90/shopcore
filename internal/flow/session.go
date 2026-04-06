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
//
// Root and detail screen use stable named identiriers.
// Catalog drill-down screen are encoded dynamically from CatalogPath.
type ScreenID string

const (
	ScreenReplyWelcome ScreenID = "reply_welcome"
	ScreenRootCompact  ScreenID = "root_compact"
	ScreenRootExtended ScreenID = "root_extended"

	ScreenCabinet   ScreenID = "cabinet"
	ScreenSupport   ScreenID = "support"
	ScreenReviews   ScreenID = "reviews"
	ScreenBalance   ScreenID = "balance"
	ScreenBotsMine  ScreenID = "bots_mine"
	ScreenOrderLast ScreenID = "order_last"

	ScreenAdminRoot               ScreenID = "admin_root"
	ScreenAdminCatalog            ScreenID = "admin_catalog"
	ScreenAdminCategoryCreate     ScreenID = "admin_category_create"
	ScreenAdminCategoryCode       ScreenID = "admin_category_code"
	ScreenAdminCategoryCreateDone ScreenID = "admin_category_create_done"
)

// PendingInputKind identifies which text input flow currently expects.
type PendingInputKind string

const (
	PendingInputNone         PendingInputKind = ""
	PendingInputCategoryName PendingInputKind = "category_name"
	PendingInputCategoryCode PendingInputKind = "category_code"
)

const (
	// PendingValueName stores one entered name value inside pending input payload.
	PendingValueName = "name"

	// PendingValueCode stores one entered code value inside pending input payload.
	PendingValueCode = "code"
)

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

// Session stores current screen, backward navigation history and pending input state.
type Session struct {
	Current ScreenID
	History []ScreenID
	Pending PendingInput
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
