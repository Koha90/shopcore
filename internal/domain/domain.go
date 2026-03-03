// Package domain contains the core business model of the application.
//
// Domain layer rules:
//   - No dependencies on infrastructure (DB, HTTP, Telegram, etc.).
//   - Aggregates enforce business invariants and control state transitions.
//   - All identifiers are assigned by the repository layer (database).
//   - Money is stored as int64 in smallest currency unit (e.g. cents), unless the project decides otherwise.
//   - Aggregates may emit DomainEvent values that are later published by the application layer.
//
// Building blocks used in this package:
//
// Aggregates
//
//   - Product
//     Owns ProductVariant entities as part of the Product aggregate boundary.
//     Product is responsible for:
//
//   - validating product state changes (rename, category change, etc.)
//
//   - adding variants and preventing duplicates among active variants
//
//   - archiving variants (with "cannot archive last active variant" invariant)
//
//   - optimistic concurrency via Version()
//
//   - buffering domain events (PullEvents)
//
//   - Order
//     Represents a purchase intent with immutable monetary snapshot at creation time.
//     Order is responsible for:
//
//   - status transitions (pending -> paid / cancelled)
//
//   - guarding rules like "paid order cannot be cancelled"
//
//   - optimistic concurrency via Version()
//
//   - buffering domain events (PullEvents)
//
// Entities / Value Objects
//
//   - ProductVariant
//     Represents packaging/offer for a product.
//     Typically includes: pack size, district reference, price, archivedAt, version.
//
//   - User
//     Represents an application user.
//     A user may authenticate via Telegram (tgID) or email/password.
//     Admin access may be time-limited (adminAccessExpiresAt).
//
//   - Category, City, District (when present)
//     Reference entities used by products and variants.
//
// # Optimistic locking
//
// Many aggregates in this project contain a version counter.
// Repository implementations should use it for optimistic concurrency.
// Typical SQL pattern:
//
//	UPDATE products
//	SET name=$1, version=version+1
//	WHERE id=$2 AND version=$3
//
// If no rows are affected, it means the aggregate was updated concurrently.
//
// # Domain events
//
// DomainEvent represents "a fact that happened" in the domain.
// Aggregates buffer events internally and expose them via PullEvents().
// Publishing is responsibility of the application layer (services) through an EventBus.
//
// # Persistence notes
//
// Domain types intentionally keep fields unexported to enforce invariants.
// Repositories should load aggregates using dedicated constructors like NewXFromDB(...)
// (or builder functions) rather than scanning directly into unexported fields.
//
// This package is designed to be tested without database access.
// Infrastructure-specific tests belong to repository packages.
package domain
