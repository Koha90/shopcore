// Package tuiapp contains wiring for the TUI runtime application.
//
// tuiapp is a composition root for the terminal operator interface of the
// e-commerce platform. It assembles infrastructure and application services
// required by cmd/tui.
//
// Responsibilities:
//   - open shared PostgreSQL resources
//   - build storage-backend configuration services
//   - build bot runtime manager
//
// It does not contain business logic, startup policy, or UI code.
package tuiapp
