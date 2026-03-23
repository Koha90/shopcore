// Package app provides PostgreSQL bootstrap helpers for application entrypoints.
//
// The package centralizes DSN construction and pgx pool initialization so that
// commands such as migrate and tui do not duplicate connection setup logic.
package pgapp
