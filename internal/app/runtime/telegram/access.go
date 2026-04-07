package telegram

import "github.com/koha90/shopcore/internal/manager"

// AdminAccessResolver resolves explicit admin access for Telegram users.
//
// The resolver is transport-side. Flow receives only the final CanAdmin flag
// and does not depend on Telegram identity details.
type AdminAccessResolver interface {
	CanAdminTelegram(spec manager.BotSpec, userID int64) bool
}

// DenyAllAdminAccessResolver denies admin access for every Telegram user.
//
// It is a safe default until explicit admin access wiring is configured.
type DenyAllAdminAccessResolver struct{}

// CanAdminTelegram reports whether Telegram user may access admin flow.
func (DenyAllAdminAccessResolver) CanAdminTelegram(spec manager.BotSpec, userID int64) bool {
	return false
}

func normalizeAdminAccessResolver(r AdminAccessResolver) AdminAccessResolver {
	if r == nil {
		return DenyAllAdminAccessResolver{}
	}

	return r
}

func (r *Runner) canAdminTelegram(spec manager.BotSpec, userID int64) bool {
	if r == nil || r.adminAccess == nil {
		return false
	}

	return r.adminAccess.CanAdminTelegram(spec, userID)
}
