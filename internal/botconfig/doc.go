// Package botconfig provides application layer for bot configuration management.
//
// It stores editable bot settings and reusable database connection profiles.
// Package is designed to serve multiple operator interfaces such as TUI,
// web admin panel and Telegram-based admin bot.
//
// Runtime lifecycle is not managed here.
// Bot start, stop and restart remain responsibility of manager package.
package botconfig
