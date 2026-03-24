// Package telegram contains Telegram bot runtime implementation.
//
// The package provides a manager.Runner implementation based on
// github.com/go-telegram/bot.
//
// It is designed for multi-bot runtime:
// each Run call receives bot-specific settings through manager.BotSpec.
package telegram
