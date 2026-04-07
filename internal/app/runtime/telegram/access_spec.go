package telegram

import "github.com/koha90/shopcore/internal/manager"

// SpecAdminAccessResolver resolves admin access directly from bot runtime spec.
type SpecAdminAccessResolver struct{}

// CanAdminTelegram reports whether user is listed in bot Telegram admin user IDs.
func (SpecAdminAccessResolver) CanAdminTelegram(spec manager.BotSpec, userID int64) bool {
	for _, id := range spec.TelegramAdminUserIDs {
		if id == userID {
			return true
		}
	}

	return false
}
