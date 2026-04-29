![CI](https://github.com/Koha90/shopcore/actions/workflows/ci.yml/badge.svg)

# shopcore

[English](README.md) · [Русский](README.ru.md)

**shopcore** is an e-commerce platform with multiple interfaces and a shared sales, catalog and bot runtime core.

It is not just a bot manager. The project is moving toward a multi-interface commerce platform where Telegram, TUI and future Web admin surfaces use the same domain logic and application services.

---

## Current focus

- Telegram as the main customer sales channel
- TUI as the current operator/admin workspace
- Web as a future admin and storefront surface
- Shared catalog, order and flow logic across interfaces
- Per-bot runtime configuration and database wiring

---

## Interfaces

### Telegram

Telegram is currently the main runtime adapter.

Implemented:

- customer-facing catalog navigation
- reply and inline keyboards
- scenario-based `/start`
- schema-driven catalog drill-down
- order confirmation flow
- persisted order creation
- admin order notifications
- admin order status actions
- customer plain-text notifications to the admin chat
- admin replies to customers from message notifications
- admin replies to customers from order notifications
- admin text replies
- admin photo replies with optional captions
- per-bot flow service
- per-bot catalog provider through `database_id`

### TUI

TUI is the current operator/admin panel.

Implemented:

- bot list
- runtime status
- start/stop/restart actions
- token editing
- configuration editing
- `StartScenario` display and editing
- runtime spec synchronization through `manager.UpdateSpec(...)`

### Web

Planned.

Expected direction:

- admin catalog editing
- order management
- customer support tools
- storefront/admin views backed by shared services

---

## Architecture

shopcore follows a clean-ish ports/adapters style.

Core principles:

- transport-agnostic flow
- explicit dependencies
- small application services
- separated read and write paths where useful
- no SQL in Telegram runtime
- no Telegram-specific code in flow
- per-bot runtime wiring
- tests with every meaningful behavior change
- small steps instead of large rewrites

Main packages:

```text
internal/
  app/
    bootstrap/        application bootstrap and runtime wiring
    pgapp/            Postgres pool/config wiring
    runtime/
      telegram/       Telegram runtime adapter
      demo/           demo runtime pieces
    seed/             demo data seeding
    tuiapp/           TUI application wiring

  botconfig/          bot configuration domain/service
  catalog/
    postgres/         Postgres catalog adapter
    service/          catalog write-side application services
  flow/               transport-agnostic navigation and input flow
  manager/            bot lifecycle manager
  order/
    postgres/         Postgres order adapter
    service/          order application service
  tui/                terminal UI
  transport/          HTTP transport pieces
```

---

## Quick start

```bash
make up
make migrate
make run-tui
```

Run tests:

```bash
go test ./...
```

---

## Flow

`internal/flow` is transport-agnostic.

It owns:

- start scenarios
- screen/view models
- reply and inline keyboard models
- session state
- history-based back navigation
- pending text/photo input state
- catalog navigation
- transport-neutral effects

Start scenarios:

```text
reply_welcome
inline_catalog
```

Catalog drill-down uses schema-driven navigation:

```text
city -> category -> district -> product -> variant
```

Generic catalog actions and screens:

```text
catalog:select:<level>:<id>
catalog:screen:<path>
```

Back navigation uses `Session.History`. It does not jump to a hardcoded root.

---

## Catalog

Catalog data is stored in Postgres and loaded into `flow.Catalog` through a `CatalogProvider`.

Current tables:

```text
cities
catalog_categories
catalog_city_categories
catalog_districts
catalog_products
catalog_variants
```

Runtime wiring:

```text
bot config
↓
database_id
↓
database profile
↓
Postgres pool
↓
catalog provider
↓
per-bot flow service
↓
Telegram runtime
```

Important catalog rules:

- products without variants are skipped
- districts without products are skipped
- categories without a valid branch are skipped
- cities without children are skipped

---

## Orders

The order flow currently supports:

- selecting a catalog variant
- confirming an order
- persisting the order
- notifying the configured admin chat
- admin status actions:
  - take into work
  - close
- refreshing admin cards after callback actions
- replying to the customer from the order notification

---

## Customer messages and admin replies

Customer plain text is handled by Telegram runtime:

1. If the session has pending input, text goes to `flow.HandleText`.
2. If there is no pending input and the sender is a customer, text is sent to the configured admin chat.
3. Admin users are skipped by customer-message notifications.

Admin reply flow:

```text
customer message or order notification
↓
admin presses reply action
↓
bot opens pending input in admin private chat
↓
admin sends text or photo
↓
flow returns a transport-neutral effect
↓
Telegram runtime sends the response to the customer
```

Supported admin replies:

- plain text
- photo with optional caption

Photo replies use a transport-owned media token. Flow treats the token as opaque data and only forwards it through `EffectSendPhoto`.

---

## Runtime

Runtime behavior:

- each bot has its own runtime spec
- each bot has its own `flow.Service`
- multiple bots may share one database
- multiple bots may use different databases
- config updates are synchronized through `manager.UpdateSpec(...)`
- token/config/start scenario changes do not require full system restart

Telegram metadata support:

- Telegram bot id
- Telegram username
- Telegram bot display name

---

## Seed data

Seed creates:

- a real database profile
- a demo bot
- demo catalog data

Demo catalog includes:

- Moscow
- Saint Petersburg
- flower and gift categories
- districts
- products
- variants

If the token is empty, the demo bot is created as disabled.

---

## What works now

- Telegram bot runtime
- `/start`
- reply welcome scenario
- inline catalog scenario
- Postgres-backed catalog loading
- schema-driven catalog traversal
- history-based back navigation
- per-bot catalog wiring by `database_id`
- TUI bot management
- runtime spec sync
- Telegram metadata sync through `getMe`
- persisted order flow
- admin order workflow
- customer message notifications
- admin text replies
- admin photo replies

---

## Development status

Baseline:

```bash
go test ./...
```

Expected result: all tests pass.

Covered areas include:

- `internal/flow`
- catalog Postgres loader/build/repository
- catalog service write use cases
- Telegram runtime contracts
- bot configuration
- order service and storage
- manager behavior

---

## Near-term roadmap

### Catalog admin

- more write use cases
- edit category/city/district/product/variant
- TUI catalog management
- Telegram admin bot/catalog tools
- later Web admin catalog editor

### Customer communication

- persist customer inquiries
- persist admin replies
- delivery status
- retry support
- message history
- richer media support

### Orders

- richer order cards
- order comments
- order history
- better admin workflow
- customer-visible order status

### Web

- admin dashboard
- catalog editing
- order management
- customer support view

---

## Development principles

- keep working behavior intact
- make small, testable changes
- avoid premature overengineering
- keep dependencies explicit
- wrap meaningful errors
- document important packages, types and functions
- add tests with logic changes
- keep flow transport-agnostic
- keep runtime free from SQL
- treat shopcore as a commerce platform, not a temporary bot script
