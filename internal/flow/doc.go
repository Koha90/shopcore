// Package flow contains transport-agnostic bot view and action models.
//
// The package defines:
//   - start scenarios
//   - action identifiers
//   - view models for inline and reply keyboards
//   - session/history based navigation
//   - schema-driven catalog navigation
//
// Catalog navigation is split into two concerns:
//   - StartScenario controls how user enters catalog
//   - CatalogSchema controls level order inside catalog
//
// Telegram runtime should render these models, not invent business behavior.
package flow
