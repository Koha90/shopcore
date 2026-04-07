package telegram

// StaticAdminAccessResolver resolves admin access from in-memory bot/user allowlist.
type StaticAdminAccessResolver struct {
	Allow map[string]map[int64]struct{}
}

// CanAdminTelegram reports whether user is explicitly allowed for bot admin flow.
func (r StaticAdminAccessResolver) CanAdminTelegram(botID string, userID int64) bool {
	if r.Allow == nil {
		return false
	}

	users, ok := r.Allow[botID]
	if !ok {
		return false
	}

	_, ok = users[userID]
	return ok
}
