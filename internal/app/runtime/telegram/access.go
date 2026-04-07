package telegram

// AdminAccessResolver resolves explicit admin access for Telegram users.
//
// The resolver is transport-side. Flow receives only the final CanAdmin flag
// and does not depend on Telegram identity details.
type AdminAccessResolver interface {
	CanAdminTelegram(botID string, userID int64) bool
}

// DenyAllAdminAccessResolver denies admin access for every Telegram user.
//
// It is a safe default until explicit admin access wiring is configured.
type DenyAllAdminAccessResolver struct{}

// CanAdminTelegram reports whether Telegram user may access admin flow.
func (DenyAllAdminAccessResolver) CanAdminTelegram(botID string, userID int64) bool {
	return false
}
