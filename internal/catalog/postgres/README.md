# Catalog architecture

## Goal

Provide catalog data to `internal/flow` without coupling flow to SQL or transport-specific logic.

## Layers

### `internal/flow`
Responsible for:

- catalog navigation contract
- `CatalogSchema`
- `CatalogPath`
- `ViewModel`
- session/history navigation
- `ActionBack`

`flow` does not know where catalog data came from.

### `internal/catalog/postgres`
Responsible for:

- reading active catalog data from Postgres
- mapping relational rows into `flow.Catalog`
- adapting storage loading to `flow.CatalogProvider`

This package does not contain Telegram-specific logic.

### runtime
Responsible for:

- selecting bot runtime config
- choosing the right database by `database_id`
- creating a per-bot `flow.Service`
- wiring `CatalogProvider` into flow

## Current catalog path

The current drill-down path is:

`city -> category -> district -> product -> variant`

This path is represented in `flow` by `CatalogSchema`.

## Current relational model

Catalog storage currently uses these tables:

- `cities`
- `catalog_categories`
- `catalog_city_categories`
- `catalog_districts`
- `catalog_products`
- `catalog_variants`

## Important rules

### 1. `flow` receives `flow.Catalog`, not SQL rows
Storage shape and SQL stay outside flow.

### 2. Each bot runtime gets its own flow service
This allows different bots to use different databases.

### 3. Catalog source is selected by bot `database_id`
`bot_configs.database_id` determines which database profile should be used.

### 4. Incomplete catalog branches are filtered out
During Postgres catalog build:

- product without variants is skipped
- district without products is skipped
- category without district/product branch is skipped
- city without children is skipped

This keeps customer-facing catalog navigation free from empty dead-end screens.

## Why this shape

This design keeps responsibilities separated:

- Postgres adapter loads data
- flow controls navigation
- runtime performs per-bot wiring

Because of that, Telegram, TUI, and future Web surfaces can share the same catalog/navigation core.

## Future direction

Later steps are expected to include:

- wiring Postgres catalog provider by `database_id`
- bot-aware catalog loading in runtime
- admin editing of catalog entities from TUI
- Telegram admin bot using the same application services
- future Web admin on top of the same backend logic
