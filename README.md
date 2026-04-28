![CI](https://github.com/Koha90/shopcore/actions/workflows/ci.yml/badge.svg)

# shopcore

**shopcore** is an e-commerce platform with multiple interfaces and a shared sales/navigation core.

Primary direction right now:

- **Telegram** as the main sales channel
- **TUI** as operator/admin workspace
- **Web** as the next admin surface
- later possible: online catalog views built on the same flow core

---

## What it is

shopcore is **not just a bot manager**.

It is a sales platform where different interfaces use the same domain and flow ideas:

- catalog navigation
- сценарии входа
- заявки и действия пользователя
- bot runtime lifecycle
- configuration editing
- future admin and storefront surfaces

---

## Current interfaces

### Telegram
- customer-facing catalog flow
- reply and inline navigation
- scenario-based start behavior
- customer plain text notifications to the configured admin chat

### TUI
- bot operator panel
- runtime status
- config editing
- token editing
- start/stop/restart flows
- StartScenario switching

### Web
- planned admin and store-facing surface

---

## Architecture

shopcore follows a **clean-ish ports/adapters approach**.

Core ideas:

- transport-agnostic flow
- separated bot runtime lifecycle
- isolated bot configuration service
- step-by-step evolution instead of large rewrites
- platform-first thinking, not one-off bot scripting

Main building blocks:

- `internal/flow` — transport-agnostic navigation and view models
- `internal/manager` — bot lifecycle and runtime state
- `internal/botconfig` — bot configuration domain/service
- `internal/app/runtime/telegram` — Telegram runtime adapter
- `internal/tui` — operator interface

---

## Quick start

```bash
make up
make migrate
make run-tui
```

---

## Project structure

```text
cmd/
  migrate/     database migrations entry point
  tui/         TUI application entry point

internal/
  app/
    bootstrap/ bootstrap enabled bots
    pgapp/     postgres pool/config wiring
    runtime/
      telegram/ telegram runtime adapter
      demo/     demo runtime pieces
    seed/      demo data
    tuiapp/    app wiring for TUI mode

  botconfig/   bot configuration business logic
  flow/        transport-agnostic flow/navigation core
  manager/     bot lifecycle manager
  tui/         terminal UI
  transport/   HTTP transport pieces
```

---

## Flow navigation

`internal/flow` is transport-agnostic and drives user navigation for Telegram and future interfaces.

Navigation is split into two separate concerns:

- `StartScenario` controls how the user enters catalog:
  - `reply_welcome`
  - `inline_catalog`
- `CatalogSchema` controls catalog drill-down order:
  - current demo schema: `city -> category -> district -> product -> variant`

State is tracked through:

- `SessionKey`
- `Session`
- `Session.History`

`ActionBack` always uses session history instead of jumping to a hardcoded root.

Catalog drill-down uses generic encoded actions and screens:

- action: `catalog:select:<level>:<id>`
- screen: `catalog:screen:<path>`

This allows catalog order to evolve without rewriting transport logic.

## Catalog storage

Catalog data is stored in Postgres and loaded into `internal/flow` through `CatalogProvider`.

Current relational model:

- `cities`
- `catalog_categories`
- `catalog_city_categories`
- `catalog_districts`
- `catalog_products`
- `catalog_variants`

Runtime does not query catalog tables directly.

Instead:

- bot config provides `database_id`
- runtime builds a per-bot flow service
- flow service uses a catalog provider
- Postgres catalog provider loads rows and builds `flow.Catalog`

This keeps flow transport-agnostic and allows different bots to use different databases.

Current catalog drill-down path:

`city -> category -> district -> product -> variant`

---

## Telegram customer messages

Telegram runtime treats plain text in two steps:

1. If the current flow session has pending input, text is passed to `flow.HandleText`.
2. If there is no pending input and the sender is a customer, text is sent to `AdminOrdersChatID` as an operator notification.

Admin users are skipped by customer-text notifications. Their text remains reserved for admin pending input and future admin commands.

The runtime still does not query SQL or mutate catalog data directly. Customer message forwarding is a transport concern and only builds an admin-facing Telegram notification card.

---

## What already works

### Flow
- transport-agnostic service in `internal/flow`
- scenario-aware `/start`
- compact and extended catalog roots
- session/history-based back navigation
- schema-driven catalog navigation
- generic catalog actions and catalog screens

### Bot configuration
- `StartScenario` stored in config
- validation for start scenario
- postgres and in-memory repositories updated
- runtime list uses stored start scenario

### TUI
- shows StartScenario
- edits StartScenario
- edits tokens
- displays enabled/disabled bots
- runtime summary and actions
- syncs runtime spec after config/token changes

### Runtime
- `manager.UpdateSpec(...)` already updates runtime spec
- token/config changes can be picked up without full system restart
- bot restart is enough to apply updated runtime spec
- customer plain text messages can be forwarded to the configured admin chat

### Infra
- Postgres
- migrations
- seed data
- bootstrap startup flow

---

## Development status

Current baseline:

- `go test ./...` passes
- `internal/flow` has strong coverage on stable navigation behavior
- `internal/app/runtime/telegram` is covered on stable transport contracts
- TUI bot management is in a solid state
- flow is now moving from demo entities to a real sales-tree model

---

## Near-term roadmap

### 1. Catalog evolution
Move from demo navigation to richer commerce structure:

- cities
- categories
- districts
- products
- variants
- later: cart and payment

### 2. Catalog source abstraction
Current flow uses demo in-memory catalog data.
Next step is to inject a catalog provider so flow can later consume real config/storage-backed data.

### 3. UX improvements
- better Telegram catalog rendering
- richer product/variant cards
- stronger TUI ergonomics
- future Web admin/store integration

---

## Principles

- do not rewrite the whole house at night
- keep working behavior intact
- prefer testable seams
- keep names short and explicit
- reduce magic
- treat shopcore as a sales platform, not a temporary bot utility

