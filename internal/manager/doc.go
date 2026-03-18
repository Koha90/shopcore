// Package manager coordinates runtime lifecycle of bot-based channels.
//
// Manager is responsible for registration, start, stop, restart, status,
// and read-model access for bot runtimes. It does not manager bot
// configuration storage. Configuration concerns belong to the botconfig
// package
//
// Typical flow:
//   - register bot runtime specifications
//   - start or stop runtimes by ID
//   - query current runtime state for TUI, web admin, or telegram admin bot.
package manager
