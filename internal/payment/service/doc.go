// Package service contains payment application use cases.
//
// The package owns payment-facing service models and narrow ports used by
// high layers to read or change payment configuration. It does not depend on
// Postgres, Telegram, TUI or any other transport/storage adapter.
//
// Current responsiblities:
//   - describes payment method kinds
//   - expose bot-enabled payment method
//   - validate payment service input
//
// Storage adapter implement the reader/writer ports from this package.
// Runtime and flow layers should depend on this package through application
// services or their own narrow ports, not through database details.
package service
