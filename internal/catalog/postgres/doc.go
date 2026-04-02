// Package postgres contains Postgres-backed catalog loading for flow.
//
// The package reads relational catalog tables and builds flow.Catalog:
//
//	city -> category -> district -> product -> variant
//
// This package is an adapter layer.
// It does not contain Telegram-specific logic and does not define flow behavior.
// Its responsibility is to load active catalog data from Postgres and map it
// into transport-agnostic flow models.
package postgres
