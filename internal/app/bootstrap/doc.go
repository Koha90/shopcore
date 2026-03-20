// Package bootstrap wires persistent configuration with runtime manager.
//
// The package is responsible for application startup orchestration:
// loading enabled bots from botconfig storage, registring them in manager,
// and starting their runtimes.
package bootstrap
