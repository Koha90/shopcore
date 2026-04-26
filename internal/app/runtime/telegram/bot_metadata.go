package telegram

import "strings"

// BotMetadata contains Telegram bot identity returned by getMe.
type BotMetadata struct {
	ID       int64
	Username string
	Name     string
}

// normalizeTelegramUsername trims spaces and removes duplicate leading @.
func normalizeTelegramUsername(username string) string {
	return strings.TrimPrefix(strings.TrimSpace(username), "@")
}
