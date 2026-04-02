// Package flow contains transport-agnostic bot view and action models.
//
// The package defines:
//   - start scenarios
//   - action identifiers
//   - view models for inline and reply keyboards
//   - session/history based navigation
//   - schema-driven catalog navigation
//   - catalog provider abstraction
//
// Catalog navigation is split into two concerns:
//   - StartScenario controls how user enters catalog
//   - CatalogSchema controls level order inside catalog
//
// Flow operates on flow.Catalog and does not depend on SQL or transport code.
package flow
